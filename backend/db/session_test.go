package db

import (
	"imaginaer/models"
	_ "imaginaer/testing_init"
	"testing"
	"time"
)

func TestSessionCreateGetDelete(t *testing.T) {
	store := NewStore()

	userID, _ := store.CreateUser(models.User{Username: "TestSessionCreateDelete", Password: ""})
	sessionID, err := store.CreateSession(userID, time.Now().Add(1000000))
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	session, err := store.SessionByID(sessionID)
	if err != nil {
		t.Fatalf("SessionByID: %v", err)
	}
	if session.UserID != userID {
		t.Fatalf("SessionByID: expected usedID = %v, got %v", userID, session.UserID)
	}

	t.Cleanup(func() {
		err = store.DeleteSession(sessionID)
		if err != nil {
			t.Fatalf("DeleteSession: %v", err)
		}
	})
}
