package handlers

import (
	"bytes"
	"codex/src/memory"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func chdirTemp(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd error: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir error: %v", err)
	}
	return func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("restore cwd error: %v", err)
		}
	}
}

func TestChatHandlerStoresConversation(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	origSend := llamaSendPrompt
	llamaSendPrompt = func(prompt string) (string, error) {
		return "echo: " + prompt, nil
	}
	defer func() { llamaSendPrompt = origSend }()

	body, _ := json.Marshal(ChatRequest{Prompt: "hello", Project: "proj", Role: "user"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ChatHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var resp ChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp: %v", err)
	}
	if resp.Response != "echo: hello" {
		t.Fatalf("unexpected response: %s", resp.Response)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	entries, err := memory.LastNEntries(db, "proj", 10)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	seen := map[string]string{}
	for _, e := range entries {
		seen[e.Role] = e.Content
	}
	if seen["assistant"] != "echo: hello" {
		t.Fatalf("assistant entry missing or unexpected: %+v", entries)
	}
	if seen["user"] != "hello" {
		t.Fatalf("user entry missing or unexpected: %+v", entries)
	}
}

func TestChatHandlerUsesActiveProject(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	if err := memory.AddProject(db, "alpha"); err != nil {
		t.Fatalf("AddProject error: %v", err)
	}
	if err := memory.SetActiveProject(db, "alpha"); err != nil {
		t.Fatalf("SetActiveProject error: %v", err)
	}
	db.Close()

	origSend := llamaSendPrompt
	llamaSendPrompt = func(prompt string) (string, error) { return "ok", nil }
	defer func() { llamaSendPrompt = origSend }()

	body, _ := json.Marshal(ChatRequest{Prompt: "ping"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ChatHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	db, err = memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	entries, err := memory.LastNEntries(db, "alpha", 5)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Project != "alpha" {
			t.Fatalf("entries not stored under active project: %+v", entries)
		}
	}
}

func TestChatHandlerRejectsEmptyPrompt(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader([]byte(`{"prompt":"   "}`)))
	rec := httptest.NewRecorder()

	ChatHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rec.Code)
	}
}

func TestChatHandlerDefaultsToProject(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	origSend := llamaSendPrompt
	llamaSendPrompt = func(prompt string) (string, error) { return "result", nil }
	defer func() { llamaSendPrompt = origSend }()

	body, _ := json.Marshal(ChatRequest{Prompt: "hi"})
	req := httptest.NewRequest(http.MethodPost, "/api/chat", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ChatHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	db, err := memory.InitDB()
	if err != nil {
		t.Fatalf("InitDB error: %v", err)
	}
	defer db.Close()

	entries, err := memory.LastNEntries(db, "default", 5)
	if err != nil {
		t.Fatalf("LastNEntries error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	projects, err := memory.ListProjects(db)
	if err != nil {
		t.Fatalf("ListProjects error: %v", err)
	}
	found := false
	for _, p := range projects {
		if p == "default" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("default project not created: %+v", projects)
	}
}

func TestChatHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
	rec := httptest.NewRecorder()

	ChatHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed, got %d", rec.Code)
	}
}
