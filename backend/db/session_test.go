package db

import (
	"imaginaer/auth"
	_ "imaginaer/testing_init"
	"testing"
	"time"
)

func TestSessionCreateGetDelete(t *testing.T) {
	store := NewStore()

	userID, _ := store.CreateUser(auth.User{Username: "TestSessionCreateDelete", Password: ""})
	token := auth.SessionToken("TestSessionCreateGetDelete")
	err := store.CreateSession(token, userID, time.Now().Add(1000000))
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	session, err := store.SessionByToken(token)
	if err != nil {
		t.Fatalf("SessionByToken: %v", err)
	}
	if session.UserID != userID {
		t.Fatalf("SessionByToken: expected usedID = %v, got %v", userID, session.UserID)
	}

	t.Cleanup(func() {
		err = store.DeleteSession(token)
		if err != nil {
			t.Fatalf("DeleteSession: %v", err)
		}
	})
}
