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
	GetSongs(params dto.GetSongsDto) ([]domain.Song, int, error)
	GetSong(songID int32, params dto.PaginationParamsDto) ([]string, int, error)
	Delete(songID int32) error
	Update(songID int32, updateSongInput dto.SongParamsDto) (domain.Song, error)
	Create(createSongInput dto.CreateSongDto) (domain.Song, error)
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

func (h SongsHandler) getSongs(w http.ResponseWriter, r *http.Request, params dto.GetSongsDto) {
	songs, totalPages, err := h.songsService.GetSongs(params)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSongs)
		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingSongs})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, h.toSongsDto(songs, totalPages))
}

func (h SongsHandler) getSongText(w http.ResponseWriter, r *http.Request, songID int, params dto.PaginationParamsDto) {
	verses, totalVerses, err := h.songsService.GetSong(int32(songID), params)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrGettingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JsonError{Error: delivery.ErrGettingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, dto.VersesDto{
		Verses:      verses,
		TotalVerses: totalVerses,
	})
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

func (h SongsHandler) updateSong(w http.ResponseWriter, r *http.Request, songID int, updateSongInput dto.SongParamsDto) {
	song, err := h.songsService.Update(int32(songID), updateSongInput)
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

	delivery.RespondWithJSON(w, http.StatusOK, h.toSongDto(song))
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

	delivery.RespondWithJSON(w, http.StatusCreated, h.toSongDto(song))
}

func (h SongsHandler) toSongsDto(songs []domain.Song, totalPages int) dto.SongsDto {
	var songsDto []dto.SongDto

	for _, song := range songs {
		songDto := h.toSongDto(song)
		songsDto = append(songsDto, songDto)
	}

	return dto.SongsDto{
		Songs:      songsDto,
		TotalPages: totalPages,
	}
}

func (h SongsHandler) toSongDto(song domain.Song) dto.SongDto {
	var releaseDate string
	if !song.ReleaseDate.IsZero() {
		releaseDate = song.ReleaseDate.Format(domain.DateFormat)
	}

	songDto := dto.SongDto{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: releaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}

	return songDto
}
