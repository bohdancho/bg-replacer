package auth_test

import (
	"bytes"
	"encoding/json"
	"imaginaer/auth"
	"imaginaer/db"
	_ "imaginaer/testing_init"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthFlow(t *testing.T) {
	store := db.NewStore()
	mux := auth.NewMux(store)

	regPayload := auth.RegistrationDTO{Username: "TestLoginLogout", Password: "invalid_password"}
	var regReqBuf bytes.Buffer
	err := json.NewEncoder(&regReqBuf).Encode(regPayload)
	if err != nil {
		t.Fatal(err)
	}
	regReq, err := http.NewRequest("POST", "/registration", &regReqBuf)
	if err != nil {
		t.Fatal(err)
	}
	regRR := httptest.NewRecorder()
	mux.ServeHTTP(regRR, regReq)
	if status := regRR.Code; status != http.StatusOK {
		t.Fatalf("RegistrationHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	t.Cleanup(func() {
		user, err := store.UserByUsername(regPayload.Username)
		if err != nil {
			t.Fatal(err)
		}
		store.DeleteUser(user.ID)
	})

	var loginReqBuf bytes.Buffer
	loginPayload := auth.LoginDTO(regPayload)
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

	getUserReq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	for _, c := range loginRR.Result().Cookies() {
		getUserReq.AddCookie(c)
	}

	rr := httptest.NewRecorder()
	_, err = auth.GetCurrentUser(rr, getUserReq, store)
	if err != nil {
		t.Errorf("GetCurrentUser: %v", err)
	}

	logoutReq, err := http.NewRequest("POST", "/logout", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range loginRR.Result().Cookies() {
		logoutReq.AddCookie(c)
	}

	logoutRR := httptest.NewRecorder()
	mux.ServeHTTP(logoutRR, logoutReq)
	if status := logoutRR.Code; status != http.StatusOK {
		t.Errorf("LogoutHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	store := db.NewStore()
	mux := auth.NewMux(store)

	payload := auth.RegistrationDTO{Username: "TestRegisterInvalidPassword", Password: ""}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/registration", &buf)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	status := rr.Code
	expectedStatus := http.StatusBadRequest
	if status != expectedStatus {
		t.Errorf("wrong status code: got %v expected %v",
			status, expectedStatus)
	}
}

func TestGetCurrentUserNoCookie(t *testing.T) {
	store := db.NewStore()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	_, err = auth.GetCurrentUser(w, req, store)
	if err == nil {
		t.Errorf("expected an error")
	}
}

func TestGetCurrentUserInvalidCookie(t *testing.T) {
	store := db.NewStore()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := auth.NewSessionCookie("invalid", 1000000)
	req.AddCookie(&cookie)

	w := httptest.NewRecorder()
	_, err = auth.GetCurrentUser(w, req, store)
	if err == nil {
		t.Errorf("expected an error")
	}
}

func TestLoginNoUser(t *testing.T) {
	store := db.NewStore()
	mux := auth.NewMux(store)

	payload := auth.LoginDTO{
		Username: "TestLoginUnauthorized",
		Password: "some_password",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/login", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	store := db.NewStore()
	mux := auth.NewMux(store)

	username := "TestLoginWrongPassword"
	id, err := store.CreateUser(auth.User{Username: username, Password: "valid_password"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := store.DeleteUser(id)
		if err != nil {
			t.Fatal(err)
		}
	})

	payload := auth.LoginDTO{
		Username: username,
		Password: "invalid_password",
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", &buf)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusUnauthorized)
	}
}
