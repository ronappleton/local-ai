package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"codex/src/auth"
	"codex/src/email"
	"codex/src/memory"
)

func TestPasswordResetFlow(t *testing.T) {
	cleanup := setupTemp(t)
	defer cleanup()

	db, _ := memory.InitDB()
	if err := auth.CreateUser(db, "alice", "a@b.com", "oldpwd"); err != nil {
		t.Fatalf("create user: %v", err)
	}
	auth.MarkVerified(db, 1)
	db.Close()

	var body string
	oldSend := email.SendFunc
	email.SendFunc = func(to, subject, b string) error {
		body = b
		return nil
	}
	defer func() { email.SendFunc = oldSend }()

	reqBody := bytes.NewBufferString(`{"Email":"a@b.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/reset/request", reqBody)
	w := httptest.NewRecorder()
	ResetRequestHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("reset request failed: %d", w.Result().StatusCode)
	}

	r := regexp.MustCompile(`token=([a-f0-9]+)`)
	m := r.FindStringSubmatch(body)
	if len(m) != 2 {
		t.Fatalf("token not found in email: %q", body)
	}
	token := m[1]

	reqBody = bytes.NewBufferString(`{"Token":"` + token + `","Password":"newpwd"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/reset", reqBody)
	w = httptest.NewRecorder()
	ResetPasswordHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("reset password failed: %d", w.Result().StatusCode)
	}

	reqBody = bytes.NewBufferString(`{"Username":"alice","Password":"newpwd"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/login", reqBody)
	w = httptest.NewRecorder()
	LoginHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("login with new password failed: %d", w.Result().StatusCode)
	}
}
