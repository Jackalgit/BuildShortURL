package handlers

import (
	"flag"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	"io"
	"log"
	"net/http"
)

type ShortURL struct {
	url dicturl.DictURL
}

func NewShortURL() *ShortURL {
	return &ShortURL{
		url: make(dicturl.DictURL),
	}

}

func (s *ShortURL) MakeShortURL(w http.ResponseWriter, r *http.Request) {
	// Оставляю проверку метода т.к. во 2 инкременте мы тестируем работу функции, а не работу запущенного сервера.
	// Или как вариан из тестов убрать проверку методов запроса.

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

	shortURLKey := s.AddOriginalURL(originalURL)

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	flag.Parse()
	fmt.Println(fmt.Sprint(config.Config.BaseAddress, config.Config.ServerPort, "/", shortURLKey))
	w.Write([]byte(fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)))

}

func (s *ShortURL) GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	shortURLKey := r.URL.Path[1:]
	if shortURLKey == "" {
		http.Error(w, "Don't shortUrlKey", http.StatusBadRequest)
		return
	}

	originalURL, found := s.url[shortURLKey]
	if !found {
		http.Error(w, "originalURL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (s *ShortURL) AddOriginalURL(originalURL []byte) string {
	shortURLKey := util.GenerateKey()
	s.url[shortURLKey] = originalURL

	return shortURLKey

}
