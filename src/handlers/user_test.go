package handlers

import (
	"bytes"
	"encoding/json"
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
	body := bytes.NewBufferString(`{"Username":"alice","Email":"a@b.com","Password":"secret"}`)
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

	body = bytes.NewBufferString(`{"Email":"a@b.com","Password":"secret"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/login", body)
	w = httptest.NewRecorder()
	LoginHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login failed: %d", res.StatusCode)
	}
	found := false
	for _, c := range res.Cookies() {
		if c.Name == "session" && c.Value != "" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("session cookie not set")
	}
	var u struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Admin    bool   `json:"admin"`
	}
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if u.Username != "alice" || u.Email != "a@b.com" || !u.Verified {
		t.Fatalf("unexpected user response: %+v", u)
	}
}
