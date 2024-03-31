package auth

import (
	"encoding/json"
	"net/http"
)

func (s Server) currentUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(w, r, s.store)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// TODO: put this in context?
func GetCurrentUser(w http.ResponseWriter, r *http.Request, store Store) (UserDTO, error) {
	sessionToken, err := sessionTokenFromCookie(r)
	if err != nil {
		return UserDTO{}, err
	}

	user, err := store.UserBySessionToken(sessionToken)
	if err != nil {
		if err == ErrUserNotFound {
			removeSessionCookie(w)
		}
		return UserDTO{}, err
	}

	return UserDTO{ID: user.ID, Username: user.Username}, nil
}
