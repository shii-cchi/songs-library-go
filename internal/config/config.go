package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var errEnvVarNotDefined = errors.New("environment variable is not defined")

type Config struct {
	Port            string
	DbUser          string
	DbPassword      string
	DbHost          string
	DbPort          string
	DbName          string
	MusicInfoApiUrl string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT: %w", errEnvVarNotDefined)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER: %w", errEnvVarNotDefined)
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD: %w", errEnvVarNotDefined)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST: %w", errEnvVarNotDefined)
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("DB_PORT: %w", errEnvVarNotDefined)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME: %w", errEnvVarNotDefined)
	}

	musicInfoApiUrl := os.Getenv("MUSIC_INFO_API_URL")
	if musicInfoApiUrl == "" {
		return nil, fmt.Errorf("MUSIC_INFO_API_URL: %w", errEnvVarNotDefined)
	}

	return &Config{
		Port:            port,
		DbUser:          dbUser,
		DbPassword:      dbPassword,
		DbHost:          dbHost,
		DbPort:          dbPort,
		DbName:          dbName,
		MusicInfoApiUrl: musicInfoApiUrl,
	}, nil
}
