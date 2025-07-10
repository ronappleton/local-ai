package llama

// Package llama contains a minimal HTTP client used to communicate with the
// local LLM server. The AI uses this package to send prompts and receive
// completions.
//
// AI Awareness: Any change to this client affects how the assistant interacts
// with its underlying language model. Additional parameters or authentication
// mechanisms can be added here to support new model backends.
import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type completionRequest struct {
	// Prompt is the text sent to the LLM.
	Prompt string `json:"prompt"`
	// NPredict controls how many tokens to generate.
	NPredict int `json:"n_predict"`
	// Temperature controls sampling randomness.
	Temperature float64 `json:"temperature"`
}

type completionResponse struct {
	// Content is the raw text returned by the LLM server.
	Content string `json:"content"`
}

// SendPrompt sends the provided prompt to the local LLM server and returns the
// generated text. It is used by higher level handlers such as the HTTP API.
func SendPrompt(prompt string) (string, error) {
	reqBody := completionRequest{
		Prompt:      prompt,
		NPredict:    300,
		Temperature: 0.7,
	}

	body, _ := json.Marshal(reqBody)

	// TODO: make the endpoint configurable so alternative LLM backends can
	// be targeted in the future.
	resp, err := http.Post("http://localhost:8080/completion", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	// Parse the JSON response returned by the LLM server. If the format
	// changes in future this section will need to adapt.
	var respData completionResponse
	json.Unmarshal(resBody, &respData)

	return respData.Content, nil
}
