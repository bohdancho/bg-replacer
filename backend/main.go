package main

import (
	"fmt"
	"log"
	"net/http"

	"imaginaer/auth"
	"imaginaer/db"
	"imaginaer/processing"
)

func main() {
	store := db.NewStore()

	authMux := auth.NewMux(store)
	http.Handle("/api/", http.StripPrefix("/api", authMux))

	processingMux := processing.NewMux()
	http.Handle("/api/processing/", http.StripPrefix("/api/processing", processingMux))

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":"true"}`))
	})

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
