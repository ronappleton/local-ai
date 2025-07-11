package handlers

import (
	"codex/src/models"
	"encoding/json"
	"net/http"
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
