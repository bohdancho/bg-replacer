package main

import (
	"fmt"
	"log"
	"net/http"

	"imaginaer/auth"
	_ "imaginaer/db"
	"imaginaer/processing"
)

func main() {
	http.HandleFunc("/api/grayscale", processing.GrayscaleHandler)
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":"true"}`))
	})
	http.HandleFunc("/api/login", auth.LoginHandler)
	http.HandleFunc("/api/logout", auth.LogoutHandler)
	http.HandleFunc("/api/registration", auth.RegistrationHandler)

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
