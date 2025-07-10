package handlers

// Project related HTTP handlers providing CRUD operations for conversation
// contexts.  Each project maps to a separate memory namespace in the database.

import (
	"codex/src/memory"
	"encoding/json"
	"net/http"
	"strings"
)

// This file implements HTTP handlers for managing assistant projects. Projects
// allow the AI to maintain multiple independent memory banks.

// ProjectsResponse is returned by GET /projects
type ProjectsResponse struct {
	// Projects contains the list of all known project names.
	Projects []string `json:"projects"`
	// Active indicates which project is currently selected.
	Active string `json:"active"`
}

// ProjectsHandler handles listing and creating projects. It responds to both
// GET and POST on the /projects endpoint and interacts with the memory package
// for persistence. Extension Point: additional methods such as PUT could be
// added here to extend project metadata management.
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Additional HTTP methods (PUT for rename etc.) could be supported here
	// in the future.

	switch r.Method {
	case http.MethodGet:
		list, err := memory.ListProjects(db)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		active, _ := memory.GetActiveProject(db)
		json.NewEncoder(w).Encode(ProjectsResponse{Projects: list, Active: active})
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		if err := memory.AddProject(db, req.Name); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// SwitchProjectHandler sets the active project. Clients post a project name to
// /projects/switch to change context for subsequent chat operations. AI
// Awareness: by switching project, the assistant focuses memory queries on a
// different conversation context.
func SwitchProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Use the shared memory database to update the active project setting.
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := memory.SetActiveProject(db, req.Name); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProjectHandler removes a project identified by /projects/{name}. It
// will also unset the active project if that project is being deleted. This is
// primarily used by the CLI and HTTP API when the user wants to discard a
// conversation history.
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Extract the project name from the URL path.
	name := strings.TrimPrefix(r.URL.Path, "/projects/")
	if name == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Remove the project using the memory package which will also tidy up
	// any associated settings.
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	if err := memory.DeleteProject(db, name); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RenameProjectHandler renames a project. This updates all stored memories and
// the active project setting so the assistant remains consistent across the
// database.
func RenameProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Old == "" || req.New == "" {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := memory.RenameProject(db, req.Old, req.New); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
