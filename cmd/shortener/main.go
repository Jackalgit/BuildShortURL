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
	config.ConfigServerPort()
	config.ConfigBaseAddress()

}

func main() {

	flag.Parse()

	if err := runServer(); err != nil {
		log.Println("runServer ERROR: ", err)
	}

}

func runServer() error {

	dictURL := handlers.NewShortURL()

	router := mux.NewRouter()

	router.HandleFunc("/", dictURL.MakeShortURL).Methods("POST")
	router.HandleFunc("/{id}", dictURL.GetURL).Methods("GET")

	return http.ListenAndServe(config.Config.ServerPort, router)

}
