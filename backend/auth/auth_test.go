package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthFlow(t *testing.T) {
	regPayload := registrationDTO{Username: "TestLoginLogout", Password: "invalid_password"}
	var regReqBuf bytes.Buffer
	err := json.NewEncoder(&regReqBuf).Encode(regPayload)
	if err != nil {
		t.Fatal(err)
	}
	regReq, err := http.NewRequest("POST", "/api/registration", &regReqBuf)
	if err != nil {
		t.Fatal(err)
	}
	regRR := httptest.NewRecorder()
	regHandler := http.HandlerFunc(RegistrationHandler)
	regHandler.ServeHTTP(regRR, regReq)
	if status := regRR.Code; status != http.StatusOK {
		t.Errorf("RegistrationHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	t.Cleanup(func() {
		user, err := userByUsername(regPayload.Username)
		if err != nil {
			t.Fatal(err)
		}
		deleteUser(user.ID)
	})

	var loginReqBuf bytes.Buffer
	loginPayload := loginDTO(regPayload)
	err = json.NewEncoder(&loginReqBuf).Encode(loginPayload)
	if err != nil {
		t.Fatal(err)
	}

	loginReq, err := http.NewRequest("POST", "/api/login", &loginReqBuf)
	if err != nil {
		t.Fatal(err)
	}
	loginRR := httptest.NewRecorder()
	loginHandler := http.HandlerFunc(LoginHandler)
	loginHandler.ServeHTTP(loginRR, loginReq)
	if status := loginRR.Code; status != http.StatusOK {
		t.Fatalf("LoginHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	cookies := loginRR.Result().Cookies()
	var sessionCookie *http.Cookie
	fmt.Println(cookies)
	for _, c := range cookies {
		if c.Name == sessionCookieName {
			sessionCookie = c
		}
	}
	if sessionCookie == nil {
		t.Fatal("LoginHandler did not add session cookie")
	}

	getUserReq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	getUserReq.AddCookie(sessionCookie)
	_, err = GetCurrentUser(getUserReq)
	if err != nil {
		t.Errorf("GetCurrentUser: %v", err)
	}

	logoutReq, err := http.NewRequest("POST", "/api/logout", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	logoutReq.AddCookie(sessionCookie)

	logoutRR := httptest.NewRecorder()
	logoutHandler := http.HandlerFunc(LogoutHandler)
	logoutHandler.ServeHTTP(logoutRR, logoutReq)
	if status := logoutRR.Code; status != http.StatusOK {
		t.Errorf("LogoutHandler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}
}

func TestGetCurrentUserNoCookie(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = GetCurrentUser(req)
	if err == nil {
		t.Errorf("expected an error")
	}
}

func TestGetCurrentUserInvalidCookie(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	cookie := newSessionCookie(0, sessionCookieMaxAge)
	req.AddCookie(&cookie)

	_, err = GetCurrentUser(req)
	if err == nil {
		t.Errorf("expected an error")
	}
}

func TestLoginNoUser(t *testing.T) {
	payload := loginDTO{
		Username: "TestLoginUnauthorized",
		Password: "some_password",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/login", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	username := "TestLoginWrongPassword"
	id, err := createUser(User{username: username, password: "valid_password"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		deleteUser(id)
	})

	payload := loginDTO{
		Username: username,
		Password: "invalid_password",
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/login", &buf)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusUnauthorized)
	}
}
