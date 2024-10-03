package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"songs-library-go/internal/config"
	"songs-library-go/internal/delivery/handlers"
	"songs-library-go/internal/repository"
	"songs-library-go/internal/service"
	"songs-library-go/internal/validator"
	"time"
)

const (
	errLoadingConfig  = "error loading config"
	errConnectingToDb = "error connecting to db"

	successfulConfigLoad     = "config has been loaded successfully"
	successfulConnectionToDb = "successfully connected to db"
	serverStart              = "server starting on port"
)

const pingInterval = 10 * time.Second

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal(errLoadingConfig)
	}
	log.Info(successfulConfigLoad)

	conn, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		log.WithError(err).Fatal(errConnectingToDb)
	}
	defer conn.Close()
	log.Info(successfulConnectionToDb)

	go pingDatabase(conn)

	songsRepo := repository.NewSongsRepo(conn)

	v := validator.Init()
	songsService := service.NewSongsService(songsRepo, cfg.MusicInfoApiUrl)

	r := chi.NewRouter()
	songsHandler := handlers.NewSongsHandler(v, songsService)
	songsHandler.RegisterRoutes(r)

	log.Infof(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func pingDatabase(conn *sql.DB) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.Ping()
			if err != nil {
				log.WithError(err).Fatal(errConnectingToDb)
			}
		}
	}
}
