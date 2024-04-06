package processing

import (
	"image"
	"image/draw"
	"imaginaer/codecs"
	"net/http"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /grayscale", grayscaleHandler)
	return mux
}

func grayscaleHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	format, err := codecs.AssertSupportedImageType(contentType)
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
