package main

import (
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	if err := runServer(); err != nil {
		log.Println("runServer ERROR: ", err)
	}

}

func runServer() error {

	dictURL := handlers.NewShortURL()

	router := mux.NewRouter()

	router.HandleFunc("/", dictURL.MakeShortURL)
	router.HandleFunc("/{id}", dictURL.GetURL)

	return http.ListenAndServe(":8080", router)

}
