package app

import (
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"songs-library-go/internal/config"
	"songs-library-go/internal/delivery/handlers"
	"songs-library-go/internal/repository"
	"songs-library-go/internal/service"
	"songs-library-go/internal/validator"
)

const serverStart = "server starting on port"

// Run initializes whole application.
func Run() {
	cfg := config.Init()

	conn := repository.Init(cfg)
	defer conn.Close()

	songsRepo := repository.NewSongsRepo(conn)

	v := validator.Init()
	songsService := service.NewSongsService(songsRepo, cfg.MusicInfoAPIURL)

	r := chi.NewRouter()
	songsHandler := handlers.NewSongsHandler(v, songsService)
	songsHandler.RegisterRoutes(r)

	log.Infof(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
