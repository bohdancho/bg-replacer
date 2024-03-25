package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ImgRequest struct {
	Img string `json:"img"`
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var request ImgRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	img, err := decodeImage(request.Img)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	grayImg := grayscaleImage(img)

	grayImg64, err := encodeImage(grayImg)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json, err := json.Marshal(grayImg64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func decodeImage(s string) (image.Image, error) {
	header, content, found := strings.Cut(s, ";base64,")
	if !found {
		return nil, errors.New("no ';base64,' separator found")
	}
	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	switch strings.TrimPrefix(header, "data:") {
	case "image/png":
		return png.Decode(bytes.NewReader(b))
	case "image/jpeg":
		return jpeg.Decode(bytes.NewReader(b))
	}
	return nil, errors.New("invalid mime type")
}

func encodeImage(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	return string(append([]byte("data:image/png;base64,"), s...)), nil

}

func grayscaleImage(img image.Image) image.Image {
	result := image.NewGray(img.Bounds())
	draw.Draw(result, result.Bounds(), img, img.Bounds().Min, draw.Src)
	return result
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/grayscale", grayscaleHandler)

	http.Handle("/", r)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
