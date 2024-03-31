package db

import (
	"imaginaer/auth"
	_ "imaginaer/testing_init"
	"testing"
)

func TestUserCreateGetDelete(t *testing.T) {
	store := NewStore()

	user := auth.User{Username: "TestAddGetDeleteUser", Password: "pwd"}
	id, err := store.CreateUser(user)
	user.ID = id

	if err != nil {
		t.Fatalf("createUser: %v", err)
	}

	t.Cleanup(func() {
		err = store.DeleteUser(id)
		if err != nil {
			t.Fatalf("deleteUser: %v", err)
		}
	})

	byUsername, err := store.UserByUsername(user.Username)
	if err != nil {
		t.Fatalf("userByUsername: %v", err)
	}
	if byUsername.ID != user.ID {
		t.Fatalf("userByUsername expected id: %v, received id: %v", user.ID, byUsername.ID)
	}
	if byUsername.Username != user.Username {
		t.Fatalf("userByUsername expected Username: %v, received Username: %v", user.Username, byUsername.Username)
	}
	if byUsername.Password != user.Password {
		t.Fatalf("userByUsername expected Password: %v, received Password: %v", user.Password, byUsername.Password)
	}

	byId, err := store.UserByID(id)
	if err != nil {
		t.Fatalf("userById: %v", err)
	}
	if byId.ID != user.ID {
		t.Fatalf("userByID expected id: %v, received id: %v", user.ID, byId.ID)
	}
	if byId.Username != user.Username {
		t.Fatalf("userByID expected Username: %v, received Username: %v", user.Username, byId.Username)
	}
	if byId.Password != user.Password {
		t.Fatalf("userByID expected Password: %v, received Password: %v", user.Password, byId.Password)
	}
}

func TestAddUserConflict(t *testing.T) {
	store := NewStore()

	user := auth.User{Username: "TestAddUserConflict", Password: "pwd"}
	id, err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("error while creating the first user: %v", err)
	}
	t.Cleanup(func() {
		store.DeleteUser(id)
	})

	_, err = store.CreateUser(user)
	if err != auth.ErrUsernameTaken {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestDeleteNonexistentUser(t *testing.T) {
	store := NewStore()

	err := store.DeleteUser(1000000000000000)
	if err != auth.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
