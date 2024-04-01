package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"imaginaer/auth"
	"imaginaer/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	store := db.NewStore()
	mux := auth.NewMux(store)

	user := auth.User{Username: "TestCurrentUser", Password: "password123123"}
	id, err := store.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := store.DeleteUser(id)
		if err != nil {
			t.Fatal(err)
		}
	})

	var loginReqBuf bytes.Buffer
	loginPayload := auth.LoginDTO{Username: user.Username, Password: user.Password}
	err = json.NewEncoder(&loginReqBuf).Encode(loginPayload)
	if err != nil {
		t.Fatal(err)
	}

	loginReq, err := http.NewRequest("POST", "/login", &loginReqBuf)
	if err != nil {
		t.Fatal(err)
	}
	loginRR := httptest.NewRecorder()
	mux.ServeHTTP(loginRR, loginReq)
	if status := loginRR.Code; status != http.StatusOK {
		t.Fatalf("LoginHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	getUserRR := httptest.NewRecorder()
	getUserReq, err := http.NewRequest("GET", "/current-user", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range loginRR.Result().Cookies() {
		getUserReq.AddCookie(c)
	}
	mux.ServeHTTP(getUserRR, getUserReq)

	status := getUserRR.Code
	expectedStatus := http.StatusOK
	if status != expectedStatus {
		t.Errorf("wrong status code: expected %v, got %v", expectedStatus, status)
	}

	body := getUserRR.Body.String()
	expectedBody := fmt.Sprintf(`{"id":%d,"username":"%s"}\n`, id, user.Username)
	if body == expectedBody {
		t.Errorf("wrong body: expected %v, got %v", expectedBody, body)
	}
}
