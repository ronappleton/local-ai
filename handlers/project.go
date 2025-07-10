package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"local-ai/memory"
)

// ProjectsResponse is returned by GET /projects
type ProjectsResponse struct {
	Projects []string `json:"projects"`
	Active   string   `json:"active"`
}

// ProjectsHandler handles listing and creating projects.
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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

// SwitchProjectHandler sets the active project.
func SwitchProjectHandler(w http.ResponseWriter, r *http.Request) {
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

// DeleteProjectHandler deletes a project with /projects/{name}
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/projects/")
	if name == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
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
