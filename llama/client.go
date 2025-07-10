package llama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type completionRequest struct {
	Prompt      string  `json:"prompt"`
	NPredict    int     `json:"n_predict"`
	Temperature float64 `json:"temperature"`
}

type completionResponse struct {
	Content string `json:"content"`
}

func SendPrompt(prompt string) (string, error) {
	reqBody := completionRequest{
		Prompt:      prompt,
		NPredict:    300,
		Temperature: 0.7,
	}

	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8080/completion", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	var respData completionResponse
	json.Unmarshal(resBody, &respData)

	return respData.Content, nil
}
