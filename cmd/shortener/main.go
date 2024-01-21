package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/initialization"
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
	config.ConfigDatabaseDSN()
}

func main() {
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	ctx := context.Background()
	//ctx, _ := signal.NotifyContext(
	//	context.Background(),
	//	syscall.SIGINT,
	//	syscall.SIGTERM,
	//	syscall.SIGQUIT,
	//)

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

	storage := initialization.InitStorage(ctx)

	router := mux.NewRouter()

	router.HandleFunc("/ping", storage.PingDB).Methods("GET")
	router.HandleFunc("/", storage.MakeShortURL).Methods("POST")
	router.HandleFunc("/{id}", storage.GetURL).Methods("GET")
	router.HandleFunc("/api/shorten", storage.APIShortURL).Methods("POST")
	router.HandleFunc("/api/shorten/batch", storage.Batch).Methods("POST")

	router.Use(logger.LoggingMiddleware, zip.GzipMiddleware)

	if err := http.ListenAndServe(config.Config.ServerPort, router); err != nil {
		return fmt.Errorf("[ListenAndServe] запустить сервер: %q", err)

	}

	return http.ListenAndServe(config.Config.ServerPort, router)

}
