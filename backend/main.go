package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"net/http"

	"imaginaer/codecs"
)

func main() {
	http.HandleFunc("/api/grayscale", grayscaleHandler)
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"ok\":\"true\"}")
	})
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := 8080
	fmt.Printf("Server started at http://localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	format, err := codecs.AssertSupportedFormat(contentType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	img, err := codecs.DecodeImage(r.Body, format)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	grayImg := grayscaleImage(img)
	grayImgBytes, err := codecs.EncodeImage(grayImg, format)
	if err != nil {
		fmt.Println(err)
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
