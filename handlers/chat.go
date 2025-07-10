package handlers

import (
	"encoding/json"
	_ "io"
	"local-ai/llama"
	"net/http"
)

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := llama.SendPrompt(req.Prompt)
	if err != nil {
		http.Error(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := ChatResponse{Response: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
