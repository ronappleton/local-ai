package handlers

import (
	"codex/src/auth"
	"codex/src/memory"
	"encoding/json"
	"net/http"
	"strconv"
)

// RegisterHandler creates a new user account.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct{ Username, Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	if err := auth.CreateUser(db, req.Username, req.Email, req.Password); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler authenticates a user and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct{ Username, Password, TOTP string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	u, err := auth.Authenticate(db, req.Username, req.Password, req.TOTP)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: strconv.Itoa(u.ID), Path: "/"})
	w.WriteHeader(http.StatusOK)
}

// LogoutHandler clears the session cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	w.WriteHeader(http.StatusOK)
}

// UsersHandler lists users for the admin UI.
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	list, err := auth.List(db)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}
