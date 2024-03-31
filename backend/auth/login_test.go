package auth

import (
	"testing"
)

func TestNewSessionToken(t *testing.T) {
	s, err := newSessionToken()
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 255 {
		t.Fatalf("wrong token length, want: %d got: %d", 255, len(s))
	}
}
