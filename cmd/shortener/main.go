package main

import (
	"flag"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/Jackalgit/BuildShortURL/internal/logger"
	"github.com/Jackalgit/BuildShortURL/internal/zip"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func init() {
	config.ConfigServerPort()
	config.ConfigBaseAddress()
	config.ConfigLogger()
	config.ConfigFileStorage()

}

func main() {

	flag.Parse()

	if err := runServer(); err != nil {
		log.Println("runServer ERROR: ", err)
	}

}

func runServer() error {

	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", config.Config.ServerPort))

	dictURL := handlers.NewShortURL()

	router := mux.NewRouter()

	router.HandleFunc("/", dictURL.MakeShortURL).Methods("POST")
	router.HandleFunc("/{id}", dictURL.GetURL).Methods("GET")
	router.HandleFunc("/api/shorten", dictURL.APIShortURL).Methods("POST")
	router.Use(logger.LoggingMiddleware, zip.GzipMiddleware)

	return http.ListenAndServe(config.Config.ServerPort, router)

}
