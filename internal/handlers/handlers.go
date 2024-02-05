package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/jobertask"
	"github.com/Jackalgit/BuildShortURL/internal/jsondecoder"
	"github.com/Jackalgit/BuildShortURL/internal/jwt"
	"github.com/Jackalgit/BuildShortURL/internal/logger"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/userid"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"time"
)

type Repository interface {
	AddURL(ctx context.Context, userID uuid.UUID, shortURLKey string, originalURL []byte) error
	GetURL(ctx context.Context, userID uuid.UUID, shortURLKey string) ([]byte, models.StatusURL)
	AddBatchURL(ctx context.Context, userID uuid.UUID, batchList []models.BatchURL) error
	UserURLList(ctx context.Context, userID uuid.UUID) ([]models.ResponseUserURL, bool, error)
}

type ShortURL struct {
	Storage         Repository
	DictUserIDToken userid.DictUserIDToken
}

func (s *ShortURL) MakeShortURL(w http.ResponseWriter, r *http.Request) {
	// Оставляю проверку метода т.к. во 2 инкременте мы тестируем работу функции, а не работу запущенного сервера.
	// Или как вариан из тестов убрать проверку методов запроса.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
		return
	}
	cookieStr := cookie.Value
	userID, err := jwt.GetUserID(cookieStr)
	if err != nil {
		http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
		return
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
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
	shortURLKeyFull := fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)

	ctx := r.Context()

	if err := s.Storage.AddURL(ctx, userID, shortURLKeyFull, originalURL); err != nil {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURLKeyFull))

}

func (s *ShortURL) GetURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
		return
	}
	cookieStr := cookie.Value
	userID, err := jwt.GetUserID(cookieStr)
	if err != nil {
		http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
		return
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	logger.Log.Info("Передаваемый ключ в пути запроса", zap.String("url", r.URL.Path[1:]))

	shortURLKey := r.URL.Path[1:]
	if shortURLKey == "" {
		http.Error(w, "Don't shortUrlKey", http.StatusBadRequest)
		return
	}

	shortURLKeyFull := fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)

	ctx := r.Context()

	originalURL, status := s.Storage.GetURL(ctx, userID, shortURLKeyFull)

	logger.Log.Info("originalURL при GET запросе", zap.String("url", string(originalURL)))

	if status.Delete {
		http.Error(w, "originalURL delete", http.StatusGone)
		return
	}
	if !status.Found {
		http.Error(w, "originalURL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (s *ShortURL) JSONShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
		return
	}
	cookieStr := cookie.Value
	userID, err := jwt.GetUserID(cookieStr)
	if err != nil {
		http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
		return
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	request, err := jsondecoder.RequestJSONToStruct(r.Body)
	if err != nil {
		http.Error(w, "Not read body", http.StatusBadRequest)
		return
	}
	originalURL := request.URL

	logger.Log.Info("originalURL при запросе эндпоинта /api/shorten", zap.String("url", originalURL))

	shortURLKey := util.GenerateKey()

	shortURLKeyFull := fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey)

	ctx := r.Context()

	if err := s.Storage.AddURL(ctx, userID, shortURLKeyFull, []byte(originalURL)); err != nil {
		respons := models.Response{
			Result: err.Error(),
		}
		responsJSON, err := json.Marshal(respons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write(responsJSON)

		return
	}

	respons := models.Response{
		Result: shortURLKeyFull,
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

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (s *ShortURL) Batch(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")

	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
		return
	}
	cookieStr := cookie.Value
	userID, err := jwt.GetUserID(cookieStr)
	if err != nil {
		http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
		return
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	requestList, err := jsondecoder.RequestListJSONToStruct(r.Body)
	if err != nil {
		http.Error(w, "Not read body", http.StatusBadRequest)
		return
	}
	// создаем структуру для ответа хендлера
	var responseList []models.ResponseBatch
	// создаем структуру которую передадим для хранения в память или в базу данных
	var batchList []models.BatchURL

	for _, v := range requestList {
		shortURLKey := util.GenerateKey()

		responseBatch := models.ResponseBatch{
			Correlation: v.Correlation,
			ShortURL:    fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey),
		}
		responseList = append(responseList, responseBatch)

		batchURL := models.BatchURL{
			Correlation: v.Correlation,
			ShortURL:    fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey),
			OriginalURL: v.OriginalURL,
		}
		batchList = append(batchList, batchURL)
	}

	ctx := r.Context()

	if err := s.Storage.AddBatchURL(ctx, userID, batchList); err != nil {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))

		return
	}

	responsJSON, err := json.Marshal(responseList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responsJSON)

}

func (s *ShortURL) TokenMiddleware(next http.Handler) http.Handler {
	tokenFn := func(w http.ResponseWriter, r *http.Request) {
		// добываем токен из запроса
		cookie, err := r.Cookie("token")
		// если токена нет, то генерируем его и возвращаем клиенту
		if errors.Is(err, http.ErrNoCookie) {
			s.SetCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}
		// добываем значение токена, а из него userId.
		// если токен не валидный, то генерируем новый токен и возвращаем его клиенту
		cookieStr := cookie.Value
		userID, err := jwt.GetUserID(cookieStr)
		if err != nil {
			s.SetCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}
		// если в токене нет userId, то возвращаем ошибку авторизации
		if userID.String() == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// если токен валидный, но userId нет в DictUserId, то генерируем токен и возвращаем его клиенту
		if _, ok := s.DictUserIDToken[userID]; !ok {
			s.SetCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(tokenFn)
}

func (s *ShortURL) SetCookie(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()

	tokenString := jwt.BuildJWTString(id)
	s.DictUserIDToken.AddUserID(id, tokenString)

	cookie := http.Cookie{Name: "token", Value: tokenString}
	r.AddCookie(&cookie)
	http.SetCookie(w, &cookie)
}

func (s *ShortURL) UserDictURL(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		cookie, err := r.Cookie("token")

		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
			return
		}
		cookieStr := cookie.Value
		userID, err := jwt.GetUserID(cookieStr)
		if err != nil {
			http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
			return
		}
		if userID.String() == "" {
			http.Error(w, "No User ID in token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		userURLList, foundDictUser, err := s.Storage.UserURLList(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !foundDictUser {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		responsJSON, err := json.Marshal(userURLList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responsJSON)
	}

	if r.Method == http.MethodDelete {

		cookie, err := r.Cookie("token")

		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "[MakeShortURL] No Cookie", http.StatusUnauthorized)
			return
		}
		cookieStr := cookie.Value
		userID, err := jwt.GetUserID(cookieStr)
		if err != nil {
			http.Error(w, "[MakeShortURL] Token is not valid", http.StatusUnauthorized)
			return
		}
		if userID.String() == "" {
			http.Error(w, "No User ID in token", http.StatusUnauthorized)
			return
		}

		requestList, err := jsondecoder.RequestListURLDelete(r.Body)
		if err != nil {
			http.Error(w, "Not read body", http.StatusBadRequest)
			return
		}

		jobID := uuid.New()

		job := jobertask.NewJober(jobID, userID, requestList).DeleteURL()
		jobertask.JobDict[jobID] = job

		w.WriteHeader(http.StatusAccepted)
	}

}
