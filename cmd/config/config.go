package config

import (
	"flag"
	"fmt"
	"os"
)

var Config struct {
	ServerPort      string
	BaseAddress     string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
}

func ConfigServerPort() {
	flag.StringVar(&Config.ServerPort, "a", "localhost:8080", "Addres local port")

	if envRunServerAddr := os.Getenv("SERVER_ADDRESS"); envRunServerAddr != "" {
		Config.ServerPort = envRunServerAddr

	}

}

func ConfigBaseAddress() {
	flag.StringVar(&Config.BaseAddress, "b", "http://localhost:8080", "Base local addres")

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		Config.BaseAddress = envBaseURL

	}

}

func ConfigLogger() {
	flag.StringVar(&Config.LogLevel, "l", "info", "log level")

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		Config.LogLevel = envLogLevel
	}

}

func ConfigFileStorage() {
	flag.StringVar(&Config.FileStoragePath, "f", "tmp/short-url-db.json", "Path to FileStorage")

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		Config.FileStoragePath = envFileStoragePath
	}

}

func ConfigDatabaseDSN() {
	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `5432`, `ivan`, `XXXXXXXX`, `shorturl`)

	flag.StringVar(&Config.DatabaseDSN, "d", ps, "Database source name")

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		Config.DatabaseDSN = envDatabaseDSN
	}

}
