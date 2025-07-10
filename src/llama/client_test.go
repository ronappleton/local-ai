package llama

// Unit tests for the LLM client wrapper. They spin up a temporary HTTP server
// that mimics the real model endpoint.

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSendPrompt runs the client against a mocked LLM HTTP server to ensure
// prompts are encoded correctly and responses are decoded as expected.
func TestSendPrompt(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Skip("port 8080 unavailable")
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req completionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode err: %v", err)
		}
		resp := completionResponse{Content: "echo: " + req.Prompt}
		json.NewEncoder(w).Encode(resp)
	}))
	srv.Listener = ln
	srv.Start()
	defer srv.Close()

	out, err := SendPrompt("hello")
	if err != nil {
		t.Fatalf("SendPrompt error: %v", err)
	}
	if out != "echo: hello" {
		t.Fatalf("unexpected response: %s", out)
	}
}
