package config

import (
	"flag"
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
	path := "tmp/short-url-db.json"

	flag.StringVar(&Config.FileStoragePath, "f", path, "Path to FileStorage")

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		Config.FileStoragePath = envFileStoragePath
	}

}

func ConfigDatabaseDSN() {
	//ps := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable",
	//	`localhost`, `5432`, `shorturl`)

	flag.StringVar(&Config.DatabaseDSN, "d", "", "Database source name")

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		Config.DatabaseDSN = envDatabaseDSN
	}

}
