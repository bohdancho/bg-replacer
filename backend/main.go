package main

import (
	"fmt"
	"log"
	"net/http"

	"imaginaer/auth"
	"imaginaer/processing"
)

func main() {
	http.HandleFunc("/api/grayscale", processing.GrayscaleHandler)
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"ok\":\"true\"}"))
	})
	http.HandleFunc("/api/login", auth.LoginHandler)
	http.HandleFunc("/api/logout", auth.LogoutHandler)
	http.HandleFunc("/api/registration", auth.RegistrationHandler)
	http.HandleFunc("/api/protected", protectedHandler)

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "all good, you are %v", user)
}
