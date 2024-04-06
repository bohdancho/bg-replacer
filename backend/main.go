package main

import (
	"fmt"
	"log"
	"net/http"

	"imaginaer/auth"
	"imaginaer/db"
	"imaginaer/gallery"
	"imaginaer/processing"
)

func main() {
	store := db.NewStore()

	authMux := auth.NewMux(store)
	http.Handle("/api/", http.StripPrefix("/api", authMux))

	processingMux := processing.NewMux()
	http.Handle("/api/processing/", http.StripPrefix("/api/processing", processingMux))

	galleryMux := gallery.NewMux(store)
	http.Handle("/api/gallery/", http.StripPrefix("/api/gallery", galleryMux))

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":"true"}`))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
