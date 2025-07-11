package handlers

// Tests covering the project management HTTP API. They validate that the
// handler layer properly persists data using the memory package.

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

// setupTempDB creates a temporary working directory and initialises the SQLite
// database there. Tests use this to ensure isolation. It returns a cleanup
// function that restores the original working directory.
func setupTempDB(t *testing.T) func() {
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

// TestProjectAPI exercises the full project management API: creating,
// switching, listing, renaming and deleting projects. It validates the
// integration between the HTTP handlers and the memory layer.
func TestProjectAPI(t *testing.T) {
	cleanup := setupTempDB(t)
	defer cleanup()

	val, _ := sc.Encode("session", 1)
	cookie := &http.Cookie{Name: "session", Value: val}

	// create project p1
	body := bytes.NewBufferString(`{"name":"p1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/projects", body)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()
	WithAuth(ProjectsHandler)(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201")
	}

	// switch to p1
	body = bytes.NewBufferString(`{"name":"p1"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/projects/switch", body)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	WithAuth(SwitchProjectHandler)(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("switch failed")
	}

	// list
	req = httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	WithAuth(ProjectsHandler)(w, req)
	var resp ProjectsResponse
	if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Projects) != 1 || resp.Active != "p1" {
		t.Fatalf("unexpected list %+v", resp)
	}

	// rename project p1 to p2
	body = bytes.NewBufferString(`{"old":"p1","new":"p2"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/projects/rename", body)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	WithAuth(RenameProjectHandler)(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("rename failed")
	}

	// verify rename
	req = httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	WithAuth(ProjectsHandler)(w, req)
	if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Projects) != 1 || resp.Projects[0] != "p2" || resp.Active != "p2" {
		t.Fatalf("unexpected list after rename %+v", resp)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/api/projects/p1", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	WithAuth(DeleteProjectHandler)(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("delete failed")
	}
}
