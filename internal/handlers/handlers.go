package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type ShortUrl struct {
	url map[string][]byte
}

func NewShortUrl() *ShortUrl {
	return &ShortUrl{
		url: make(map[string][]byte),
	}

}

func (s *ShortUrl) MakeShortUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read originalURL ERROR: ", err)
	}
	if string(originalURL) == "" {
		http.Error(w, "Body don't url", http.StatusBadRequest)
		return
	}
	shortUrlKey := uuid.New().String()
	s.url[shortUrlKey] = originalURL

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://localhost:8080/%v\n", shortUrlKey)))

}

func (s *ShortUrl) GetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	shortUrlKey := r.URL.Path[1:]
	//fmt.Println(shortUrlKey)
	if shortUrlKey == "" {
		http.Error(w, "Don't shortUrlKey", http.StatusBadRequest)
		return
	}

	originalURL, found := s.url[shortUrlKey]
	fmt.Println(string(originalURL))
	if !found {
		http.Error(w, "originalURL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
	//http.Redirect(w, r, string(originalURL), http.StatusTemporaryRedirect)
}
