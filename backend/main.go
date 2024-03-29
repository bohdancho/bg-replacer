package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"net/http"

	"imaginaer/auth"
	"imaginaer/codecs"
)

func main() {
	http.HandleFunc("/api/grayscale", grayscaleHandler)
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"ok\":\"true\"}"))
	})
	http.HandleFunc("/api/login", auth.LoginHandler)
	http.HandleFunc("/api/logout", auth.LogoutHandler)
	http.HandleFunc("/api/registration", auth.RegistrationHandler)
	http.HandleFunc("/api/protected", func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUser(r)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "all good, you are %v", user)
	})

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	format, err := codecs.AssertSupportedFormat(contentType)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	img, err := codecs.DecodeImage(r.Body, format)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	grayImg := grayscaleImage(img)
	grayImgBytes, err := codecs.EncodeImage(grayImg, format)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", string(format))
	w.WriteHeader(http.StatusOK)
	w.Write(grayImgBytes)
}

func grayscaleImage(img image.Image) image.Image {
	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)
	return result
}
