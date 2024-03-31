package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	var payload LoginDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.store.UserByUsername(payload.Username)
	if err != nil {
		if err == ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.Password != payload.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expires := time.Now().Add(sessionMaxAge)
	sessionID, err := s.store.CreateSession(user.ID, expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionID)
}
