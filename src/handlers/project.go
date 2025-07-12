package handlers

// Project related HTTP handlers providing CRUD operations for conversation
// contexts.  Each project maps to a separate memory namespace in the database.

import (
	"codex/src/memory"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// This file implements HTTP handlers for managing assistant projects. Projects
// allow the AI to maintain multiple independent memory banks.

// ProjectsResponse is returned by GET /api/projects
type ProjectsResponse struct {
	// Projects contains the list of all known project names.
	Projects []string `json:"projects"`
	// Active indicates which project is currently selected.
	Active string `json:"active"`
}

// ProjectsHandler handles listing and creating projects. It responds to both
// GET and POST on the /api/projects endpoint and interacts with the memory package
// for persistence. Extension Point: additional methods such as PUT could be
// added here to extend project metadata management.
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("ProjectsHandler InitDB error: %v", err)
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
			log.Printf("ProjectsHandler ListProjects error: %v", err)
			return
		}
		active, _ := memory.GetActiveProject(db)
		resp := ProjectsResponse{Projects: list, Active: active}
		log.Printf("ProjectsHandler response %+v", resp)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("ProjectsHandler encode error: %v", err)
		}
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			log.Printf("ProjectsHandler decode error: %v", err)
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		if err := memory.AddProject(db, req.Name); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			log.Printf("ProjectsHandler AddProject error: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("ProjectsHandler method not allowed: %s", r.Method)
	}
}

// SwitchProjectHandler sets the active project. Clients post a project name to
// /api/projects/switch to change context for subsequent chat operations. AI
// Awareness: by switching project, the assistant focuses memory queries on a
// different conversation context.
func SwitchProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("SwitchProjectHandler method not allowed")
		return
	}
	// Use the shared memory database to update the active project setting.
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("SwitchProjectHandler InitDB error: %v", err)
		return
	}
	defer db.Close()
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		log.Printf("SwitchProjectHandler decode error: %v", err)
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := memory.SetActiveProject(db, req.Name); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("SwitchProjectHandler SetActiveProject error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProjectHandler removes a project identified by /api/projects/{name}. It
// will also unset the active project if that project is being deleted. This is
// primarily used by the CLI and HTTP API when the user wants to discard a
// conversation history.
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("DeleteProjectHandler method not allowed")
		return
	}
	// Extract the project name from the URL path. Requests are routed with
	// the /api prefix so we strip that as well.
	name := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if name == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Remove the project using the memory package which will also tidy up
	// any associated settings.
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("DeleteProjectHandler InitDB error: %v", err)
		return
	}
	defer db.Close()
	if err := memory.DeleteProject(db, name); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("DeleteProjectHandler DeleteProject error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RenameProjectHandler renames a project. This updates all stored memories and
// the active project setting so the assistant remains consistent across the
// database.
func RenameProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("RenameProjectHandler method not allowed")
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("RenameProjectHandler InitDB error: %v", err)
		return
	}
	defer db.Close()
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Old == "" || req.New == "" {
		log.Printf("RenameProjectHandler decode error: %v", err)
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := memory.RenameProject(db, req.Old, req.New); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("RenameProjectHandler RenameProject error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
