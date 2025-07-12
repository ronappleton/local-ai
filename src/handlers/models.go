package handlers

import (
	"codex/src/llama"
	"codex/src/memory"
	"codex/src/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// ModelsHandler lists Hugging Face models for a pipeline type.
// It requires the `pipeline` query parameter.
func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("ModelsHandler method not allowed: %s", r.Method)
		return
	}
	pipeline := r.URL.Query().Get("pipeline")
	if pipeline == "" {
		http.Error(w, "pipeline required", http.StatusBadRequest)
		log.Printf("ModelsHandler missing pipeline")
		return
	}
	refresh := r.URL.Query().Get("refresh") == "1"
	db, err := memory.InitDB()
	if err == nil {
		defer db.Close()
	} else if err != nil {
		log.Printf("ModelsHandler InitDB error: %v", err)
	}

	var list []models.ModelInfo
	if db != nil && !refresh {
		list, err = memory.GetModelList(db, pipeline)
		if err != nil {
			log.Printf("ModelsHandler GetModelList error: %v", err)
		}
	}
	if len(list) == 0 || refresh {
		list, err = models.ListModelsByType(pipeline)
		if err != nil {
			log.Printf("ModelsHandler ListModelsByType error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if db != nil {
			if err := memory.SaveModelList(db, pipeline, list); err != nil {
				log.Printf("ModelsHandler SaveModelList error: %v", err)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("ModelsHandler response count=%d", len(list))
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Printf("ModelsHandler encode error: %v", err)
	}
}

// ModelActionHandler exposes actions such as download and enable for a specific
// model. The action is taken from the URL path after the model ID.
func ModelActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/models/"), "/")
	if len(parts) > 0 {
		if decoded, err := url.PathUnescape(parts[0]); err == nil {
			parts[0] = decoded
		}
	}
	if len(parts) == 1 {
		// GET /api/models/{id}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Printf("ModelActionHandler method not allowed: %s", r.Method)
			return
		}
		id := parts[0]
		refresh := r.URL.Query().Get("refresh") == "1"
		db, err := memory.InitDB()
		if err == nil {
			defer db.Close()
		} else if err != nil {
			log.Printf("ModelActionHandler InitDB error: %v", err)
		}

		var md *models.ModelMetadata
		if db != nil && !refresh {
			md, err = memory.GetModelMetadata(db, id)
			if err != nil {
				log.Printf("ModelActionHandler GetModelMetadata error: %v", err)
			}
		}
		if md == nil || refresh {
			md, err = models.GetModelMetadata(id)
			if err != nil {
				log.Printf("ModelActionHandler GetModelMetadata remote error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if db != nil {
				if err := memory.SaveModelMetadata(db, "", md); err != nil {
					log.Printf("ModelActionHandler SaveModelMetadata error: %v", err)
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		log.Printf("ModelActionHandler detail %+v", md)
		if err := json.NewEncoder(w).Encode(md); err != nil {
			log.Printf("ModelActionHandler encode error: %v", err)
		}
		return
	}

	if len(parts) < 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Printf("ModelActionHandler bad request path=%s", r.URL.Path)
		return
	}
	id, action := parts[0], parts[1]
	if decoded, err := url.PathUnescape(id); err == nil {
		id = decoded
	}

	// Handle special case /api/models/stats/global
	if id == "stats" && action == "global" {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		stats, err := models.GetGlobalStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
		return
	}

	switch action {
	case "enable":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Printf("ModelActionHandler enable method not allowed")
			return
		}
		log.Printf("ModelActionHandler enable id=%s", id)
		path, err := models.ActivateModel(id)
		if err != nil {
			log.Printf("ActivateModel error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := llama.LoadModel(path); err != nil {
			log.Printf("LoadModel error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "download":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Printf("ModelActionHandler download method not allowed")
			return
		}
		log.Printf("ModelActionHandler download id=%s", id)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			log.Printf("Flusher unsupported")
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		progress := func(done, total int) {
			pct := int(float64(done) / float64(total) * 100)
			fmt.Fprintf(w, "data: %d\n\n", pct)
			flusher.Flush()
		}
		if _, err := models.DownloadModelWithProgress(id, progress); err != nil {
			log.Printf("DownloadModelWithProgress error: %v", err)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
		fmt.Fprintf(w, "event: done\ndata: ok\n\n")
		flusher.Flush()
	case "stats":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Printf("ModelActionHandler stats method not allowed")
			return
		}
		log.Printf("ModelActionHandler stats id=%s", id)
		stats, err := models.GetModelDetail(id)
		if err != nil {
			log.Printf("GetModelDetail stats error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Printf("ModelActionHandler stats encode error: %v", err)
		}
	default:
		http.NotFound(w, r)
	}
}

// RefreshModelsHandler re-fetches model listings from Hugging Face and updates
// the local database cache. The pipeline parameter must be supplied via query
// string and the handler responds with the latest list.
func RefreshModelsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("RefreshModelsHandler method not allowed")
		return
	}
	pipeline := r.URL.Query().Get("pipeline")
	if pipeline == "" {
		http.Error(w, "pipeline required", http.StatusBadRequest)
		log.Printf("RefreshModelsHandler missing pipeline")
		return
	}
	list, err := models.ListModelsByType(pipeline)
	if err != nil {
		log.Printf("RefreshModelsHandler ListModelsByType error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if db, err := memory.InitDB(); err == nil {
		if err := memory.SaveModelList(db, pipeline, list); err != nil {
			log.Printf("RefreshModelsHandler SaveModelList error: %v", err)
		}
		db.Close()
	} else if err != nil {
		log.Printf("RefreshModelsHandler InitDB error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("RefreshModelsHandler response count=%d", len(list))
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Printf("RefreshModelsHandler encode error: %v", err)
	}
}
