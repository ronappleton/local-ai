package handlers

import (
	"codex/src/auth"
	"codex/src/memory"
	"net/http"
)

// WithAuth ensures the request has a valid session cookie. If the cookie is
// missing or does not map to a user account a 401 response is returned.
func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getSessionID(r)
		if err != nil || id == 0 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		db, err := memory.InitDB()
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer db.Close()
		if _, err := auth.GetByID(db, id); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// WithAdmin ensures the requester is both authenticated and marked as an admin
// in the database. Non-admin users receive a 403 response.
func WithAdmin(next http.HandlerFunc) http.HandlerFunc {
	return WithAuth(func(w http.ResponseWriter, r *http.Request) {
		id, _ := getSessionID(r)
		db, err := memory.InitDB()
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer db.Close()
		u, err := auth.GetByID(db, id)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.Admin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}
