package handlers

// This file provides HTTP endpoints for user account management including
// registration, login, logout and listing users for the admin UI. It relies on
// the auth package for credential handling and the memory package for database
// access.

import (
	"codex/src/auth"
	"codex/src/memory"
	"encoding/json"
	"log"
	"net/http"
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
	u, err := auth.GetByUsername(db, req.Username)
	if err == nil {
		auth.SendVerification(db, u)
	}
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler authenticates a user and sets a session cookie.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("LoginHandler: method not allowed")
		return
	}
	var req struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		log.Printf("LoginHandler: invalid request body: %v", err)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("LoginHandler: db error: %v", err)
		return
	}
	defer db.Close()
	u, err := auth.Authenticate(db, req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		log.Printf("LoginHandler: authentication failed for %s: %v", req.Email, err)
		return
	}
	if err := setSessionCookie(w, u.ID); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Printf("LoginHandler: failed to set session cookie: %v", err)
		return
	}

	// 👇 Fix: ensure Content-Type is set
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Printf("LoginHandler: failed to marshal user: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// LogoutHandler clears the session cookie.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
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

// VerifyHandler marks an account as verified using a token sent via email.
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.NotFound(w, r)
		return
	}
	db, err := memory.InitDB()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	id, err := auth.ConsumeToken(db, token, "verify")
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	auth.MarkVerified(db, id)
	w.WriteHeader(http.StatusOK)
}

// ResetRequestHandler sends a password reset email when given an address.
func ResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct{ Email string }
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
	u, err := auth.GetByEmail(db, req.Email)
	if err == nil {
		auth.SendReset(db, u)
	}
	w.WriteHeader(http.StatusOK)
}

// ResetPasswordHandler sets a new password using a reset token.
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct{ Token, Password string }
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
	id, err := auth.ConsumeToken(db, req.Token, "reset")
	if err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	auth.SetPassword(db, id, req.Password)
	// Successfully resetting the password proves ownership of the account
	// so mark the user as verified.
	auth.MarkVerified(db, id)
	w.WriteHeader(http.StatusOK)
}
