package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/logger"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"time"
)

type Repository interface {
	AddURL(ctx context.Context, shortURLKey string, originalURL []byte)
	GetURL(ctx context.Context, shortURLKey string) ([]byte, bool)
}

type ShortURL struct {
	Ctx     context.Context
	Storage Repository
}

func (s *ShortURL) MakeShortURL(w http.ResponseWriter, r *http.Request) {
	// Оставляю проверку метода т.к. во 2 инкременте мы тестируем работу функции, а не работу запущенного сервера.
	// Или как вариан из тестов убрать проверку методов запроса.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	originalURL, err := io.ReadAll(r.Body)
	logger.Log.Info("originalURL при запросе на эндпоинта /", zap.String("url", string(originalURL)))

	if err != nil {
		log.Println("Read originalURL ERROR: ", err)
	}
	if string(originalURL) == "" {
		http.Error(w, "Body don't url", http.StatusBadRequest)
		return
	}

	shortURLKey := util.GenerateKey()

	s.Storage.AddURL(s.Ctx, shortURLKey, originalURL)

	if config.Config.DatabaseDSN == "" {
		util.SaveURLToJSONFile(config.Config.FileStoragePath, string(originalURL), shortURLKey)
	}

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)))

}

func (s *ShortURL) GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	logger.Log.Info("Передаваемый ключ в пути запроса", zap.String("url", r.URL.Path[1:]))

	shortURLKey := r.URL.Path[1:]
	if shortURLKey == "" {
		http.Error(w, "Don't shortUrlKey", http.StatusBadRequest)
		return
	}

	originalURL, found := s.Storage.GetURL(s.Ctx, shortURLKey)

	logger.Log.Info("originalURL при GET запросе", zap.String("url", string(originalURL)))
	if !found {
		http.Error(w, "originalURL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (s *ShortURL) APIShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	request, err := util.RequestJSONToStruct(r.Body)
	if err != nil {
		http.Error(w, "Not read body", http.StatusBadRequest)
		return
	}

	originalURL := []byte(request.URL)

	logger.Log.Info("originalURL при запросе эндпоинта /api/shorten", zap.String("url", string(originalURL)))

	shortURLKey := util.GenerateKey()
	s.Storage.AddURL(s.Ctx, shortURLKey, originalURL)

	if config.Config.DatabaseDSN == "" {
		util.SaveURLToJSONFile(config.Config.FileStoragePath, string(originalURL), shortURLKey)
	}

	shortURL := fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)
	respons := models.Response{
		Result: shortURL,
	}

	responsJSON, err := json.Marshal(respons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responsJSON)

}

func (s *ShortURL) PingDB(w http.ResponseWriter, r *http.Request) {

	ps := config.Config.DatabaseDSN

	db, err := sql.Open("pgx", ps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	ctx, cancel := context.WithTimeout(s.Ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Close()
	w.WriteHeader(http.StatusOK)

}
