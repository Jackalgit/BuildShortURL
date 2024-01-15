package main

import (
	"context"
	"flag"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/handlers"
	"github.com/Jackalgit/BuildShortURL/internal/logger"
	"github.com/Jackalgit/BuildShortURL/internal/zip"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.ConfigServerPort()
	config.ConfigBaseAddress()
	config.ConfigLogger()
	config.ConfigFileStorage()

}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	flag.Parse()

	if err := runServer(ctx); err != nil {
		log.Println("runServer ERROR: ", err)
	}

}

func runServer(ctx context.Context) error {

	if err := logger.Initialize(config.Config.LogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", config.Config.ServerPort))

	dictURL := handlers.NewShortURL(ctx)

	router := mux.NewRouter()

	router.HandleFunc("/ping", dictURL.PingDB).Methods("GET")
	router.HandleFunc("/", dictURL.MakeShortURL).Methods("POST")
	router.HandleFunc("/{id}", dictURL.GetURL).Methods("GET")
	router.HandleFunc("/api/shorten", dictURL.APIShortURL).Methods("POST")

	router.Use(logger.LoggingMiddleware, zip.GzipMiddleware)

	return http.ListenAndServe(config.Config.ServerPort, router)

}
