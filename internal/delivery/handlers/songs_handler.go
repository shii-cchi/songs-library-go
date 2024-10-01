package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"songs-library-go/internal/delivery"
	"songs-library-go/internal/delivery/dto"
	"songs-library-go/internal/delivery/middleware"
	"songs-library-go/internal/domain"
)

type SongsService interface {
	Create(createSongInput dto.CreateSongDto) (dto.SongDto, error)
	Update(updateSongInput dto.UpdateSongDto, songID int32) (dto.SongDto, error)
	Delete(songID int32) error
}

type SongsHandler struct {
	validator    *validator.Validate
	songsService SongsService
}

func NewSongsHandler(validator *validator.Validate, songsService SongsService) *SongsHandler {
	return &SongsHandler{
		validator:    validator,
		songsService: songsService,
	}
}

func (h SongsHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/songs", func(r chi.Router) {
		r.With(middleware.ValidateCreateSongInput(h.validator)).Post("/", h.createSong)
		r.With(middleware.ValidateIDInput, middleware.ValidateUpdateSongInput(h.validator)).Put("/{id}", h.updateSong)
		r.With(middleware.ValidateIDInput).Delete("/{id}", h.deleteSong)
	})
}

func (h SongsHandler) createSong(w http.ResponseWriter, r *http.Request) {
	createSongInput := r.Context().Value(delivery.CreateSongInputKey).(dto.CreateSongDto)

	song, err := h.songsService.Create(createSongInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrCreatingSong)

		if errors.Is(err, domain.ErrSongAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrCreatingSong, Message: domain.ErrSongAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrCreatingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusCreated, song)
}

func (h SongsHandler) updateSong(w http.ResponseWriter, r *http.Request) {
	updateSongInput := r.Context().Value(delivery.UpdateSongInputKey).(dto.UpdateSongDto)
	songID := r.Context().Value(delivery.IDInputKey).(int32)

	song, err := h.songsService.Update(updateSongInput, songID)
	if err != nil {
		log.WithError(err).Error(delivery.ErrUpdatingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JsonError{Error: delivery.ErrUpdatingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrSongAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrUpdatingSong, Message: domain.ErrSongAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrUpdatingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, song)
}

func (h SongsHandler) deleteSong(w http.ResponseWriter, r *http.Request) {
	songID := r.Context().Value(delivery.IDInputKey).(int32)

	if err := h.songsService.Delete(songID); err != nil {
		log.WithError(err).Error(delivery.ErrDeletingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JsonError{Error: delivery.ErrDeletingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrDeletingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, nil)
}
