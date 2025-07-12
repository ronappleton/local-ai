package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"codex/src/auth"
	"codex/src/memory"
)

// setupTemp creates a clean working directory with a fresh database.
func setupTemp(t *testing.T) func() {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	return func() { os.Chdir(cwd) }
}

func TestRegisterAndLogin(t *testing.T) {
	cleanup := setupTemp(t)
	defer cleanup()

	// register
	body := bytes.NewBufferString(`{"Email":"alice","Email":"a@b.com","Password":"secret"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	w := httptest.NewRecorder()
	RegisterHandler(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("register failed")
	}

	// manually verify user for login
	db, _ := memory.InitDB()
	auth.MarkVerified(db, 1)
	db.Close()

	body = bytes.NewBufferString(`{"Email":"alice","Password":"secret"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/login", body)
	w = httptest.NewRecorder()
	LoginHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("login failed: %d", w.Result().StatusCode)
	}
}
