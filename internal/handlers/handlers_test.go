package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	dicturl "github.com/Jackalgit/BuildShortURL/internal/dictURL"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/Jackalgit/BuildShortURL/internal/util"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortURL_GetURL(t *testing.T) {
	ctx := context.Background()
	dictURL := dicturl.NewDictURL()
	userID := uuid.New()

	dictURL.AddURL(ctx, userID, "qweQWErtyQ", []byte("long long long url"))

	s := ShortURL{Ctx: ctx, Storage: dictURL}

	tokenString := util.BuildJWTString(userID)
	cookie := http.Cookie{Name: "token", Value: tokenString}

	tests := []struct {
		name       string
		method     string
		statusCode int
		Location   string
		request    string
	}{
		{name: "Test MethodPut", method: http.MethodPut, request: "/", statusCode: http.StatusMethodNotAllowed},
		{name: "Test MethodDelete", method: http.MethodDelete, request: "/", statusCode: http.StatusMethodNotAllowed},
		{name: "Test MethodPost", method: http.MethodPost, request: "/", statusCode: http.StatusMethodNotAllowed},
		{name: "Test MethodGet with out request", method: http.MethodGet, request: "/", statusCode: http.StatusBadRequest},
		{name: "Test MethodGet with not found", method: http.MethodGet, request: "/hKuwkBVYTG", statusCode: http.StatusNotFound},
		{name: "Test MethodGet", method: http.MethodGet, request: "/qweQWErtyQ", statusCode: http.StatusTemporaryRedirect, Location: "long long long url"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			r := httptest.NewRequest(tc.method, tc.request, nil)
			r.AddCookie(&cookie)
			w := httptest.NewRecorder()

			s.GetURL(w, r)

			result := w.Result()
			fmt.Println(tc.statusCode)
			fmt.Println(w.Code)
			fmt.Println(result)
			require.Equal(t, tc.statusCode, w.Code, "The response code does not match what is expected")
			assert.Equal(t, tc.Location, result.Header.Get("Location"))
			err := result.Body.Close()
			require.NoError(t, err)

		})
	}
}

func TestShortURL_MakeShortURL(t *testing.T) {
	ctx := context.Background()
	dictURL := dicturl.NewDictURL()
	s := ShortURL{Ctx: ctx, Storage: dictURL}

	userID := uuid.New()
	tokenString := util.BuildJWTString(userID)
	cookie := http.Cookie{Name: "token", Value: tokenString}

	tests := []struct {
		name        string
		method      string
		statusCode  int
		Body        string
		contentType string
		shortURL    string
	}{
		{name: "Test MethodGet", method: http.MethodGet, statusCode: http.StatusMethodNotAllowed, contentType: "text/plain; charset=utf-8"},
		{name: "Test MethodPut", method: http.MethodPut, statusCode: http.StatusMethodNotAllowed, contentType: "text/plain; charset=utf-8"},
		{name: "Test MethodDelete", method: http.MethodDelete, statusCode: http.StatusMethodNotAllowed, contentType: "text/plain; charset=utf-8"},
		{name: "Test MethodPost and empty body", method: http.MethodPost, Body: "", statusCode: http.StatusBadRequest, contentType: "text/plain; charset=utf-8"},
		{
			name:        "Test MethodPost and body with URL",
			method:      http.MethodPost,
			Body:        "long long long url",
			statusCode:  http.StatusCreated,
			contentType: "text/plain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			bodyReader := strings.NewReader(tc.Body)

			r := httptest.NewRequest(tc.method, "/", bodyReader)
			r.AddCookie(&cookie)
			w := httptest.NewRecorder()

			s.MakeShortURL(w, r)

			require.Equal(t, tc.statusCode, w.Code, "The response code does not match what is expected")
			if tc.Body == "" {
				assert.Equal(t, tc.statusCode, w.Code, "The response code does not match what is expected")
			}

			result := w.Result()

			assert.Equal(t, tc.statusCode, result.StatusCode)
			assert.Equal(t, tc.contentType, result.Header.Get("Content-Type"))

			bodyResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			//shortURLKeyFull := fmt.Sprint(config.Config.BaseAddress, string(bodyResult))
			originalURL, _ := s.Storage.GetURL(ctx, userID, string(bodyResult)[1:])
			assert.Equal(t, tc.Body, string(originalURL))

		})
	}

}

func TestShortURL_JSONShortURL(t *testing.T) {
	ctx := context.Background()
	dictURL := dicturl.NewDictURL()
	s := ShortURL{Ctx: ctx, Storage: dictURL}

	userID := uuid.New()
	tokenString := util.BuildJWTString(userID)
	cookie := http.Cookie{Name: "token", Value: tokenString}

	handler := http.HandlerFunc(s.JSONShortURL)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	RequestBody := `{"url": "https://practicum.yandex.ru"}`

	tests := []struct {
		name         string
		method       string
		body         string
		statusCode   int
		expectedBody string
	}{
		{name: "Test MethodGet", method: http.MethodGet, statusCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{name: "Test MethodPut", method: http.MethodPut, statusCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{name: "Test MethodDelete", method: http.MethodDelete, statusCode: http.StatusMethodNotAllowed, expectedBody: ""},

		{name: "Test MethodPost with out body", method: http.MethodPost, statusCode: http.StatusBadRequest, expectedBody: ""},
		{
			name:         "Test MethodPost with not correct json",
			method:       http.MethodPost,
			body:         "wrong",
			statusCode:   http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "Test MethodPost with correct json",
			method:       http.MethodPost,
			body:         RequestBody,
			statusCode:   http.StatusCreated,
			expectedBody: "https://practicum.yandex.ru",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL
			req.SetCookie(&cookie)

			if len(tc.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.body)
			}

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.statusCode, resp.StatusCode(), "Response code didn't match expected")
			//fmt.Println(s.Storage)
			//fmt.Println(string(resp.Body()))
			// проверяем, что сохранилось в dictURL
			if tc.expectedBody != "" {
				var respons models.Response
				// десериализуем resp.Body json в go model Response
				json.Unmarshal(resp.Body(), &respons)

				originalURL, _ := s.Storage.GetURL(ctx, userID, respons.Result[1:])
				assert.Equal(t, tc.expectedBody, string(originalURL))
			}
		})
	}

}
