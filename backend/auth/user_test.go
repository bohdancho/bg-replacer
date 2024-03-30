package auth

import (
	_ "imaginaer/testing_init"
	"testing"
)

func TestAddGetDeleteUser(t *testing.T) {
	user := User{username: "TestAddGetDeleteUser", password: "pwd"}
	id, err := addUser(user)
	if err != nil {
		t.Fatalf("addUser: %v", err)
	}

	t.Cleanup(func() {
		err = deleteUser(id)
		if err != nil {
			t.Fatalf("deleteUser: %v", err)
		}
	})

	u, err := userByUsername(user.username)
	if u.ID != id {
		t.Fatalf("userByUsername expected id: %v, received id: %v", id, u.ID)
	}
	if u.username != user.username {
		t.Fatalf("userByUsername expected username: %v, received username: %v", user.username, u.username)
	}
	if u.password != user.password {
		t.Fatalf("userByUsername expected password: %v, received password: %v", user.password, u.password)
	}
}

func TestAddUserConflict(t *testing.T) {
	user := User{username: "TestAddUserConflict", password: "pwd"}
	id, err := addUser(user)
	if err != nil {
		t.Fatalf("error while creating the first user: %v", err)
	}
	t.Cleanup(func() {
		deleteUser(id)
	})

	_, err = addUser(user)
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
