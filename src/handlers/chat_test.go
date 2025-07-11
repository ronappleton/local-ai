package handlers

// Tests for the HTTP chat endpoint. The handler communicates with a mocked LLM
// server so that responses are deterministic.

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"codex/src/auth"
	"codex/src/memory"
)

// startMockLLM spins up a temporary HTTP server that mimics the behaviour of
// the local LLM endpoint. It allows the ChatHandler tests to run without an
// actual model. The returned function shuts the server down.
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

// setupAuthed prepares a temporary database with a verified user and returns a
// cleanup function. It is used to test endpoints that require authentication.
func setupAuthed(t *testing.T) func() {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB: %v", err)
	}
	if err := auth.CreateUser(db, "bob", "b@c.com", "pwd"); err != nil {
		t.Fatalf("create: %v", err)
	}
	auth.MarkVerified(db, 1)
	db.Close()
	return func() { os.Chdir(cwd) }
}

// TestChatHandlerSuccess posts a prompt to ChatHandler and verifies that the
// mocked LLM response is returned. This exercises the HTTP entry point used by
// the web server.
func TestChatHandlerSuccess(t *testing.T) {
	closeSrv := startMockLLM(t)
	if closeSrv != nil {
		defer closeSrv()
	}
	cleanup := setupAuthed(t)
	defer cleanup()

	body := bytes.NewBufferString(`{"prompt":"hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", body)
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})
	w := httptest.NewRecorder()
	WithAuth(ChatHandler)(w, req)
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

// TestAnonCookie ensures that an anonymous cookie is issued when the user is
// not logged in. This helps associate chat history with unauthenticated users.
func TestAnonCookie(t *testing.T) {
	closeSrv := startMockLLM(t)
	if closeSrv != nil {
		defer closeSrv()
	}

	body := bytes.NewBufferString(`{"prompt":"hi"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/chat", body)
	w := httptest.NewRecorder()
	ChatHandler(w, req)
	res := w.Result()
	found := false
	for _, c := range res.Cookies() {
		if c.Name == "anon" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("anon cookie not set")
	}
}

// TestChatHandlerMethod confirms that ChatHandler rejects non-POST requests
// with a method-not-allowed status.
func TestChatHandlerMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
	w := httptest.NewRecorder()
	ChatHandler(w, req)
	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405")
	}
}

// TestAuthMiddlewareUnauthorized verifies that WithAuth blocks requests lacking
// a valid session cookie.
func TestAuthMiddlewareUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewBufferString(`{"prompt":"hi"}`))
	w := httptest.NewRecorder()
	WithAuth(ChatHandler)(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401")
	}
}
