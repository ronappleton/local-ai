package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTempDB(t *testing.T) func() {
	dir := t.TempDir()
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	return func() { os.Chdir(cwd) }
}

func TestProjectAPI(t *testing.T) {
	cleanup := setupTempDB(t)
	defer cleanup()

	// create project p1
	body := bytes.NewBufferString(`{"name":"p1"}`)
	req := httptest.NewRequest(http.MethodPost, "/projects", body)
	w := httptest.NewRecorder()
	ProjectsHandler(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201")
	}

	// switch to p1
	body = bytes.NewBufferString(`{"name":"p1"}`)
	req = httptest.NewRequest(http.MethodPost, "/projects/switch", body)
	w = httptest.NewRecorder()
	SwitchProjectHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("switch failed")
	}

	// list
	req = httptest.NewRequest(http.MethodGet, "/projects", nil)
	w = httptest.NewRecorder()
	ProjectsHandler(w, req)
	var resp ProjectsResponse
	if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Projects) != 1 || resp.Active != "p1" {
		t.Fatalf("unexpected list %+v", resp)
	}

	// rename project p1 to p2
	body = bytes.NewBufferString(`{"old":"p1","new":"p2"}`)
	req = httptest.NewRequest(http.MethodPost, "/projects/rename", body)
	w = httptest.NewRecorder()
	RenameProjectHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("rename failed")
	}

	// verify rename
	req = httptest.NewRequest(http.MethodGet, "/projects", nil)
	w = httptest.NewRecorder()
	ProjectsHandler(w, req)
	if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Projects) != 1 || resp.Projects[0] != "p2" || resp.Active != "p2" {
		t.Fatalf("unexpected list after rename %+v", resp)
	}

	// delete
	req = httptest.NewRequest(http.MethodDelete, "/projects/p1", nil)
	w = httptest.NewRecorder()
	DeleteProjectHandler(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("delete failed")
	}
}
