package auth

import "testing"

func TestSessionCreateGetDelete(t *testing.T) {
	userID, _ := createUser(User{username: "TestSessionCreateDelete", password: ""})
	sessionID, err := createSession(userID)
	if err != nil {
		t.Fatalf("createSession: %v", err)
	}

	session, err := sessionByID(sessionID)
	if err != nil {
		t.Fatalf("userByUsername: %v", err)
	}
	if session.userID != userID {
		t.Fatalf("userByUsername: expected usedID = %v, got %v", userID, session.userID)
	}

	t.Cleanup(func() {
		err = deleteSession(sessionID)
		if err != nil {
			t.Fatalf("deleteSession: %v", err)
		}
	})
}
