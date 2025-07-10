package handlers

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startMockLLM(t *testing.T) func() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Skip("port 8080 unavailable")
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Prompt string `json:"prompt"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		json.NewEncoder(w).Encode(map[string]string{"content": "echo: " + req.Prompt})
	}))
	srv.Listener = ln
	srv.Start()
	return srv.Close
}

func TestChatHandlerSuccess(t *testing.T) {
	closeSrv := startMockLLM(t)
	if closeSrv != nil {
		defer closeSrv()
	}

	body := bytes.NewBufferString(`{"prompt":"hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	w := httptest.NewRecorder()
	ChatHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	var resp ChatResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Response != "echo: hi" {
		t.Fatalf("unexpected response: %s", resp.Response)
	}
}

func TestChatHandlerMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	w := httptest.NewRecorder()
	ChatHandler(w, req)
	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405")
	}
}
