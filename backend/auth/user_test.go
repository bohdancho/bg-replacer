package auth

import (
	_ "imaginaer/testing_init"
	"testing"
)

func TestUserAddGetDelete(t *testing.T) {
	user := User{username: "TestAddGetDeleteUser", password: "pwd"}
	id, err := createUser(user)
	user.ID = id

	if err != nil {
		t.Fatalf("createUser: %v", err)
	}

	t.Cleanup(func() {
		err = deleteUser(id)
		if err != nil {
			t.Fatalf("deleteUser: %v", err)
		}
	})

	byUsername, err := userByUsername(user.username)
	if err != nil {
		t.Fatalf("userByUsername: %v", err)
	}
	if byUsername.ID != user.ID {
		t.Fatalf("userByUsername expected id: %v, received id: %v", user.ID, byUsername.ID)
	}
	if byUsername.username != user.username {
		t.Fatalf("userByUsername expected username: %v, received username: %v", user.username, byUsername.username)
	}
	if byUsername.password != user.password {
		t.Fatalf("userByUsername expected password: %v, received password: %v", user.password, byUsername.password)
	}

	byId, err := userByID(id)
	if err != nil {
		t.Fatalf("userById: %v", err)
	}
	if byId.ID != user.ID {
		t.Fatalf("userByID expected id: %v, received id: %v", user.ID, byId.ID)
	}
	if byId.username != user.username {
		t.Fatalf("userByID expected username: %v, received username: %v", user.username, byId.username)
	}
	if byId.password != user.password {
		t.Fatalf("userByID expected password: %v, received password: %v", user.password, byId.password)
	}
}

func TestAddUserConflict(t *testing.T) {
	user := User{username: "TestAddUserConflict", password: "pwd"}
	id, err := createUser(user)
	if err != nil {
		t.Fatalf("error while creating the first user: %v", err)
	}
	t.Cleanup(func() {
		deleteUser(id)
	})

	_, err = createUser(user)
	if err != ErrUsernameTaken {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestDeleteNonexistentUser(t *testing.T) {
	err := deleteUser(1000000000000000)
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
