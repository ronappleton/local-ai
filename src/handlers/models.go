package handlers

import (
	"codex/src/llama"
	"codex/src/models"
	"encoding/json"
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

// EnableModelHandler downloads and activates the requested model ID then
// instructs the llama server to reload the model from disk.
func EnableModelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/models/"), "/")
	if len(parts) < 2 || parts[1] != "enable" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id := parts[0]
	path, err := models.EnableModel(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := llama.LoadModel(path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
