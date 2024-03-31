package auth

import "net/http"

func (s Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := sessionIDFromCookie(r)
	if err != nil {
		removeSessionCookie(w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = s.store.DeleteSession(sessionID)
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
