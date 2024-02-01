package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
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
	GetURL(ctx context.Context, userID uuid.UUID, shortURLKey string) ([]byte, bool)
	AddBatchURL(ctx context.Context, userID uuid.UUID, batchList []models.BatchURL) error
	UserURLList(ctx context.Context, userID uuid.UUID) ([]models.ResponseUserURL, bool)
}

type ShortURL struct {
	Ctx             context.Context
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
	if err == http.ErrNoCookie {
		log.Println("[MakeShortURL] No Cookie:", err)
	}
	cookieStr := cookie.Value
	userID, err := util.GetUserID(cookieStr)
	if err != nil {
		log.Println("[MakeShortURL] Token is not valid", err)
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

	if err := s.Storage.AddURL(s.Ctx, userID, shortURLKey, originalURL); err != nil {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprint(config.Config.BaseAddress, "/", err.Error())))

		return
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

	cookie, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		log.Println("[GetURL] No Cookie:", err)
	}
	cookieStr := cookie.Value
	userID, err := util.GetUserID(cookieStr)
	if err != nil {
		log.Println("[GetURL] Token is not valid", err)
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

	originalURL, found := s.Storage.GetURL(s.Ctx, userID, shortURLKey)

	logger.Log.Info("originalURL при GET запросе", zap.String("url", string(originalURL)))
	if !found {
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
	if err == http.ErrNoCookie {
		log.Println("[JSONShortURL] No Cookie:", err)
	}
	cookieStr := cookie.Value
	userID, err := util.GetUserID(cookieStr)
	if err != nil {
		log.Println("[JSONShortURL] Token is not valid", err)
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	request, err := util.RequestJSONToStruct(r.Body)
	if err != nil {
		http.Error(w, "Not read body", http.StatusBadRequest)
		return
	}
	originalURL := request.URL

	logger.Log.Info("originalURL при запросе эндпоинта /api/shorten", zap.String("url", string(originalURL)))

	shortURLKey := util.GenerateKey()

	if err := s.Storage.AddURL(s.Ctx, userID, shortURLKey, []byte(originalURL)); err != nil {
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

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	ctx, cancel := context.WithTimeout(s.Ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (s *ShortURL) Batch(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		log.Println("[Batch] No Cookie:", err)
	}
	cookieStr := cookie.Value
	userID, err := util.GetUserID(cookieStr)
	if err != nil {
		log.Println("[Batch] Token is not valid", err)
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	requestList, err := util.RequestListJSONToStruct(r.Body)
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

	if err := s.Storage.AddBatchURL(s.Ctx, userID, batchList); err != nil {
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
		if err == http.ErrNoCookie {
			s.SetCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}
		// добываем значение токена, а из него userId.
		// если токен не валидный, то генерируем новый токен и возвращаем его клиенту
		cookieStr := cookie.Value
		userID, err := util.GetUserID(cookieStr)
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

	tokenString := util.BuildJWTString(id)
	s.DictUserIDToken.AddUserID(id, tokenString)

	cookie := http.Cookie{Name: "token", Value: tokenString}
	r.AddCookie(&cookie)
	http.SetCookie(w, &cookie)
}

func (s *ShortURL) UserDictURL(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		log.Println("[UserDictURL] No Cookie:", err)
	}
	cookieStr := cookie.Value
	userID, err := util.GetUserID(cookieStr)
	if err != nil {
		log.Println("[UserDictURL] Token is not valid", err)
	}
	if userID.String() == "" {
		http.Error(w, "No User ID in token", http.StatusUnauthorized)
		return
	}

	userURLList, foundDictUser := s.Storage.UserURLList(s.Ctx, userID)
	if !foundDictUser {
		w.WriteHeader(http.StatusNoContent)
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
