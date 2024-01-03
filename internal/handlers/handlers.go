package handlers

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/models"
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
		url: dicturl.NewDictURL(),
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

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	flag.Parse()
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
	s.url.AddURL(shortURLKey, originalURL)

	return shortURLKey

}

func (s *ShortURL) APIShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	var request models.Request

	var buf bytes.Buffer
	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Not read body", http.StatusBadRequest)
		return
	}
	// десериализуем JSON в Request
	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(w, "Not parsing request json", http.StatusBadRequest)
		return
	}

	originalURL := []byte(request.URL)
	shortURLKey := s.AddOriginalURL(originalURL)
	flag.Parse()
	result := fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)

	respons := models.Response{
		Result: result,
	}

	resp, err := json.Marshal(respons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
