package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"imagine/codecs"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/grayscale", grayscaleHandler)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the image processing service!")
	})

	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type ImgDTO struct {
	Src string `json:"img"`
}

func (imgDto *ImgDTO) fromJson(bytes []byte) error {
	return json.Unmarshal(bytes, imgDto)
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	payloadBuffer, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var payload ImgDTO
	err = payload.fromJson(payloadBuffer)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	img, format, err := codecs.DecodeImage(payload.Src)
	if err == codecs.ErrUnsupportedImageFormat {
		http.Error(w, err.Error(), 400)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	grayImg := grayscaleImage(img)

	grayImgSrc, err := codecs.EncodeImage(grayImg, format)
	if err == codecs.ErrUnsupportedImageFormat {
		http.Error(w, err.Error(), 400)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsonResponse, err := json.Marshal(grayImgSrc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func grayscaleImage(img image.Image) image.Image {
	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)
	return result
}
