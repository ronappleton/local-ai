package handlers

import (
	"codex/src/llama"
	"codex/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ModelsHandler lists Hugging Face models for a pipeline type.
// It requires the `pipeline` query parameter.
func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pipeline := r.URL.Query().Get("pipeline")
	if pipeline == "" {
		http.Error(w, "pipeline required", http.StatusBadRequest)
		return
	}
	list, err := models.ListModelsByType(pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// ModelActionHandler exposes actions such as download and enable for a specific
// model. The action is taken from the URL path after the model ID.
func ModelActionHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/models/"), "/")
	if len(parts) < 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id, action := parts[0], parts[1]

	switch action {
	case "enable":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		path, err := models.ActivateModel(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := llama.LoadModel(path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "download":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
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
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
		fmt.Fprintf(w, "event: done\ndata: ok\n\n")
		flusher.Flush()
	default:
		http.NotFound(w, r)
	}
}
