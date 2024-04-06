package gallery

import (
	"bytes"
	"encoding/json"
	"errors"
	"imaginaer/auth"
	"imaginaer/codecs"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func NewMux(store Store) http.Handler {
	mux := http.NewServeMux()
	server := Server{store: store}

	mux.HandleFunc("GET /", server.getImageHandler)
	mux.HandleFunc("POST /", server.uploadImageHandler)
	mux.HandleFunc("DELETE /", server.deleteImageHandler)
	return mux
}

type Image struct {
	Url     string
	OwnerId int
}

type Server struct {
	store Store
}

type Store interface {
	ImageStore
	auth.UserStore
	auth.SessionStore
}

var ErrImageNotFound = errors.New("image not found")
var ErrImageAlreadyUploaded = errors.New("delete the existing image first before uploading a new one")

type ImageStore interface {
	CreateImageUrl(url string, ownerID auth.UserID) error
	DeleteImageUrlByOwner(ownerID auth.UserID) error
	ImageUrlByOwner(ownerID auth.UserID) (string, error)
}

func (s Server) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(w, r, s.store)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	contentType := r.Header.Get("Content-Type")
	imageType, err := codecs.AssertSupportedImageType(contentType)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var reqBytes bytes.Buffer
	reqBytes.ReadFrom(r.Body)

	fileName := uuid.NewString() + imageType.Extension()

	url, err := writeImageToFile(fileName, reqBytes.Bytes())
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = s.store.CreateImageUrl(url, user.ID)
	if err != nil {
		if err == ErrImageAlreadyUploaded {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		http.Error(w, err.Error(), 500)
	}

	writeJSON(w, map[string]string{"url": url})
}

func writeImageToFile(fileName string, bytes []byte) (url string, err error) {
	url = "static/uploads/" + fileName
	err = os.WriteFile(url, bytes, 0644)
	return url, err
}

func (s Server) getImageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(w, r, s.store)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	url, err := s.store.ImageUrlByOwner(user.ID)
	if err != nil {
		if err == ErrImageNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"url": url})
}

func (s Server) deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(w, r, s.store)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = s.store.DeleteImageUrlByOwner(user.ID)
	if err != nil {
		if err == ErrImageNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"url": url})
}

func writeJSON(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
