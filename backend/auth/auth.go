package auth

import (
	"encoding/json"
	"net/http"
)

type registrationDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// TODO: validation

	var payload registrationDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{username: payload.Username, password: payload.Password}
	_, err = createUser(user)
	if err == ErrUsernameTaken {
		http.Error(w, ErrUsernameTaken.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type loginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload loginDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userByUsername(payload.Username)
	if err != nil {
		if err == ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user.password != payload.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := createSession(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionID)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := sessionIDFromCookie(r)
	if err != nil {
		deleteSessionCookie(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = deleteSession(sessionID)
	if err != nil {
		if err == ErrSessionNotFound {
			deleteSessionCookie(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deleteSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

func GetCurrentUser(r *http.Request) (User, error) {
	sessionID, err := sessionIDFromCookie(r)
	if err != nil {
		return User{}, err
	}

	return userBySessionId(sessionID)
}
