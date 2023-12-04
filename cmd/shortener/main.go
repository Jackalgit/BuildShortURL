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

	dictUrl := handlers.NewShortUrl()

	router := mux.NewRouter()

	router.HandleFunc("/", dictUrl.MakeShortUrl)
	router.HandleFunc("/{id}", dictUrl.GetUrl)

	return http.ListenAndServe(":8080", router)

}
