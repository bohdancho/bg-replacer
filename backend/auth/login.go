package auth

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
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

	token, err := newSessionToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expires := newSessionExpiresTime()

	err = s.store.CreateSession(token, user.ID, expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, token)
}

func newSessionToken() (SessionToken, error) {
	token, error := newRandomString(255)
	return SessionToken(token), error
}

func newRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range length {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}
