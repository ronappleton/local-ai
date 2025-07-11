package handlers

import (
	"github.com/gorilla/securecookie"
	"net/http"
	"os"
)

var sc *securecookie.SecureCookie

func init() {
	hashKey := []byte(os.Getenv("CODEX_COOKIE_HASH"))
	if len(hashKey) == 0 {
		// default 32 bytes key if not provided
		hashKey = []byte("default-secret-change-me-----")
	}
	sc = securecookie.New(hashKey, nil)
}

func setSessionCookie(w http.ResponseWriter, id int) error {
	encoded, err := sc.Encode("session", id)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func getSessionID(r *http.Request) (int, error) {
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		return 0, err
	}
	var id int
	if err := sc.Decode("session", c.Value, &id); err != nil {
		return 0, err
	}
	return id, nil
}
