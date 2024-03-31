package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegistrationDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Server) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	var payload RegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := payload.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := User{Username: payload.Username, Password: payload.Password}
	_, err = s.store.CreateUser(user)
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
func (r RegistrationDTO) validate() error {
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
