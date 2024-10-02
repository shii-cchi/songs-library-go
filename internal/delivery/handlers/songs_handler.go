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
	GetSongs(page int, limit int, filters map[string]string) ([]dto.SongDto, error)
	GetSong(songID int32, page int, limit int) (dto.VerseDto, error)
	Delete(songID int32) error
	Update(updateSongInput dto.UpdateSongDto, songID int32) (dto.SongDto, error)
	Create(createSongInput dto.CreateSongDto) (dto.SongDto, error)
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
		r.Get("/", middleware.ValidateGetSongsParam(h.validator, h.getSongs))
		r.Get("/{id}", middleware.ValidateGetSongParam(h.getSongText))
		r.Delete("/{id}", middleware.ValidateIDInput(h.deleteSong))
		r.Put("/{id}", middleware.ValidateUpdateSongInput(h.validator, h.updateSong))
		r.Post("/", middleware.ValidateCreateSongInput(h.validator, h.createSong))
	})
}

func (h SongsHandler) getSongs(w http.ResponseWriter, r *http.Request, params delivery.PaginationParams, filters map[string]string) {
	songs, err := h.songsService.GetSongs(params.Page, params.Limit, filters)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSongs)

		if errors.Is(err, domain.ErrPageDoesntExist) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JsonError{Error: delivery.ErrGettingSongs, Message: domain.ErrPageDoesntExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingSongs})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, songs)
}

func (h SongsHandler) getSongText(w http.ResponseWriter, r *http.Request, songID int, params delivery.PaginationParams) {
	songText, err := h.songsService.GetSong(int32(songID), params.Page, params.Limit)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrPageDoesntExist) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JsonError{Error: delivery.ErrGettingSong, Message: domain.ErrPageDoesntExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, songText)
}

func (h SongsHandler) deleteSong(w http.ResponseWriter, r *http.Request, songID int) {
	if err := h.songsService.Delete(int32(songID)); err != nil {
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

func (h SongsHandler) updateSong(w http.ResponseWriter, r *http.Request, songID int, updateSongInput dto.UpdateSongDto) {
	song, err := h.songsService.Update(updateSongInput, int32(songID))
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

func (h SongsHandler) createSong(w http.ResponseWriter, r *http.Request, createSongInput dto.CreateSongDto) {
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
