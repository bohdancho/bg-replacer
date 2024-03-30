package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type registrationDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	minUsernameLength = 3
	maxUsernameLength = 24
	minPasswordLength = 8
	maxPasswordLength = 40
)

var (
	ErrInvalidUsernameLength = fmt.Errorf(
		"username length must be between %v and %v symbols", minUsernameLength, maxUsernameLength)
	ErrInvalidPasswordLength = fmt.Errorf(
		"password length must be between %v and %v symbols", minPasswordLength, maxPasswordLength)
)

// TODO: tests
func (r registrationDTO) validate() error {
	nameLen := len(r.Username)
	if nameLen < minUsernameLength || nameLen > maxUsernameLength {
		return ErrInvalidUsernameLength
	}

	pwdLen := len(r.Password)
	if pwdLen < minPasswordLength || pwdLen > maxPasswordLength {
		return ErrInvalidPasswordLength
	}

	return nil
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload registrationDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := payload.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		removeSessionCookie(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = deleteSession(sessionID)
	if err != nil {
		if err == ErrSessionNotFound {
			removeSessionCookie(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	removeSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) (User, error) {
	sessionID, err := sessionIDFromCookie(r)
	if err != nil {
		if err == ErrInvalidSessionCookie {
			removeSessionCookie(w)
		}
		return User{}, err
	}

	return userBySessionId(sessionID)
}
