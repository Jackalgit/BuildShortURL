package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//func TestNewShortURL(t *testing.T) {
//	tests := []struct {
//		name string
//		want *ShortURL
//	}{
//		{
//			name: "NewShortURL test #1",
//			want: &ShortURL{
//				url: nil,
//			},
//		},
//	}
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			fmt.Println(tc.want)
//			fmt.Println(NewShortURL())
//			if got := NewShortURL(); !reflect.DeepEqual(got, tc.want) {
//				t.Errorf("NewShortURL() = %v, want %v", got, tc.want)
//			}
//		})
//	}
//}

func TestShortURL_GetURL(t *testing.T) {
	s := ShortURL{url: map[string][]byte{"qweQWErtyQ": []byte("long long long url")}}

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
			w := httptest.NewRecorder()

			s.GetURL(w, r)

			result := w.Result()

			require.Equal(t, tc.statusCode, w.Code, "The response code does not match what is expected")
			assert.Equal(t, tc.Location, result.Header.Get("Location"))
			err := result.Body.Close()
			require.NoError(t, err)

		})
	}
}

func TestShortURL_MakeShortURL(t *testing.T) {

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
			shortURL:    "http://localhost",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			bodyReader := strings.NewReader(tc.Body)

			r := httptest.NewRequest(tc.method, "/", bodyReader)
			w := httptest.NewRecorder()

			s := NewShortURL()
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

			assert.Equal(t, tc.Body, string(s.url[string(bodyResult)[1:]]))

		})
	}

}
