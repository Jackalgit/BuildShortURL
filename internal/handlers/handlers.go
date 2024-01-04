package handlers

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/logger"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	"go.uber.org/zap"
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
	logger.Log.Info("тело запроса при урл /", zap.String("url", fmt.Sprint(originalURL)))

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

	logger.Log.Info("Словарь урлов", zap.String("url", fmt.Sprint(s.url)))
	logger.Log.Info("кусок пути как ключ", zap.String("url", r.URL.Path[1:]))

	shortURLKey := r.URL.Path[1:]
	if shortURLKey == "" {
		http.Error(w, "Don't shortUrlKey", http.StatusBadRequest)
		return
	}

	originalURL, found := s.url[shortURLKey]
	logger.Log.Info("Оригинальный урл байт", zap.String("url", fmt.Sprint(originalURL)))
	logger.Log.Info("Оригинальный урл", zap.String("url", string(originalURL)))
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
	logger.Log.Info("тело запроса при урл /shorten", zap.String("url", fmt.Sprint(originalURL)))
	logger.Log.Info("тело запроса при урл /shorten", zap.String("url", string(originalURL)))
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
