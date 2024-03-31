package processing

import (
	"image"
	"image/draw"
	"net/http"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/grayscale", GrayscaleHandler)
	return mux
}

func GrayscaleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	format, err := AssertSupportedFormat(contentType)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	img, err := DecodeImage(r.Body, format)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	grayImg := grayscaleImage(img)
	grayImgBytes, err := EncodeImage(grayImg, format)
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
