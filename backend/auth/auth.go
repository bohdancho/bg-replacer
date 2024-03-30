package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SessionID int64
type Session struct {
	ID      SessionID
	expires time.Time
	userID  UserID
}

type registrationDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	user := User{username: payload.Username, password: payload.Password}
	_, err = addUser(user)
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
	}
	if user.password != payload.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprint(w, "cool: ", user)

	// sessionId := uuid.New().String()
	// expires := time.Now().Add(sessionMaxAge)

	// sessions[sessionId] = Session{
	// 	user:    User{username: payload.Username},
	// 	expires: expires,
	// }
	// setSessionCookie(w, sessionId)
	// w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	//
	// sessionCookie, err := getSessionCookie(r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// sessionId := sessionCookie.Value
	// _, ok := sessions[sessionId]
	// if !ok {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	//
	// sessions[sessionId] = Session{}
	// deleteSessionCookie(w)
	// w.WriteHeader(http.StatusOK)
}

// func GetUser(r *http.Request) (User, error) {
// 	sessionCookie, err := getSessionCookie(r)
// 	if err != nil {
// 		return User{}, errors.New("no session cookie found")
// 	}
// 	sessionId := sessionCookie.Value
//
// 	session, ok := sessions[sessionId]
//
// 	if !ok {
// 		return User{}, errors.New("no session matches the session cookie")
// 	}
//
// 	return session.user, nil
// }
