package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	errLoadingConfig     = "error loading config"
	errEnvVarNotDefined  = "environment variable is not defined"
	successfulConfigLoad = "config has been loaded successfully"
)

// Config is a struct that holds the configuration settings for the application.
type Config struct {
	Port            string
	DbUser          string
	DbPassword      string
	DbHost          string
	DbPort          string
	DbName          string
	MusicInfoAPIURL string
}

// Init loads environment variables from the .env file and returns a Config struct.
func Init() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.WithError(err).Fatal(errLoadingConfig)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT: %s", errEnvVarNotDefined)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatalf("DB_USER: %s", errEnvVarNotDefined)
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatalf("DB_PASSWORD: %s", errEnvVarNotDefined)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatalf("DB_HOST: %s", errEnvVarNotDefined)
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatalf("DB_PORT: %s", errEnvVarNotDefined)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("DB_NAME: %s", errEnvVarNotDefined)
	}

	musicInfoAPIURL := os.Getenv("MUSIC_INFO_API_URL")
	if musicInfoAPIURL == "" {
		log.Fatalf("MUSIC_INFO_API_URL: %s", errEnvVarNotDefined)
	}

	log.Info(successfulConfigLoad)

	return &Config{
		Port:            port,
		DbUser:          dbUser,
		DbPassword:      dbPassword,
		DbHost:          dbHost,
		DbPort:          dbPort,
		DbName:          dbName,
		MusicInfoAPIURL: musicInfoAPIURL,
	}
}
