package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type registrationDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO: test full auth flow
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

// TODO: test bad case
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

// TODO: test bad case
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sessionCookie, err := getSessionCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionIDInt, err := strconv.Atoi(sessionCookie.Value)
	sessionID := SessionID(sessionIDInt)
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

// TODO: test

// func GetUser(r *http.Request) (User, error) {
// 	sessionCookie, err := getSessionCookie(r)
// 	if err != nil {
// 		return User{}, errors.New("no session cookie found")
// 	}
// 	sessionId := sessionCookie.Value
//
// 	session, ok := sessions[sessionId]

// TODO: validate session.expired

// 	if !ok {
// 		return User{}, errors.New("no session matches the session cookie")
// 	}
//
// 	return session.user, nil
// }
