package main

import (
	"flag"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func init() {
	config.ConfigPort()
	config.ConfigBaseAddress()

}

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

	flag.Parse()

	return http.ListenAndServe(config.Config.Port, router)

}
