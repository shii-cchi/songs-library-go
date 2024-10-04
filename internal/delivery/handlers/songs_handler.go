package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	_ "songs-library-go/docs"
	"songs-library-go/internal/delivery"
	"songs-library-go/internal/delivery/dto"
	"songs-library-go/internal/delivery/middleware"
	"songs-library-go/internal/domain"
)

// SongsService defines the methods for managing songs, including retrieval, creation, updating, and deletion.
type SongsService interface {
	GetSongs(params dto.GetSongsDto) ([]domain.Song, int, error)
	GetSongText(songID int32, params dto.PaginationParamsDto) ([]string, int, error)
	Delete(songID int32) error
	Update(songID int32, updateSongInput dto.SongParamsDto) (domain.Song, error)
	Create(createSongInput dto.CreateSongDto) (domain.Song, error)
}

// SongsHandler manages HTTP requests related to songs and validates input using the provided validator.
type SongsHandler struct {
	validator    *validator.Validate
	songsService SongsService
}

// NewSongsHandler initializes and returns a new instance of SongsHandler with the provided validator and songs service.
func NewSongsHandler(validator *validator.Validate, songsService SongsService) *SongsHandler {
	return &SongsHandler{
		validator:    validator,
		songsService: songsService,
	}
}

// RegisterRoutes sets up the HTTP routes for song-related operations using the Chi router.
func (h SongsHandler) RegisterRoutes(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/songs", func(r chi.Router) {
		r.Get("/", middleware.ValidateGetSongsParam(h.validator, h.getSongs))
		r.Get("/{id}", middleware.ValidateGetSongParam(h.validator, h.getSongText))
		r.Delete("/{id}", middleware.ValidateIDInput(h.deleteSong))
		r.Put("/{id}", middleware.ValidateUpdateSongInput(h.validator, h.updateSong))
		r.Post("/", middleware.ValidateCreateSongInput(h.validator, h.createSong))
	})
}

// @Summary Get list of songs
// @Description Retrieve a paginated list of songs based on various filters like group, song, release date, text, and link.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param body body dto.SongParamsDto true "Filters"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of songs per page"
// @Success 200 {object} dto.SongsDto "List of songs"
// @Failure 500 {object} delivery.JSONError "Internal Server Error"
// @Router /songs [get]
func (h SongsHandler) getSongs(w http.ResponseWriter, r *http.Request, params dto.GetSongsDto) {
	songs, totalPages, err := h.songsService.GetSongs(params)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSongs)
		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JSONError{Error: delivery.ErrGettingSongs})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, h.toSongsDto(songs, totalPages))
}

// @Summary Get song text by song ID
// @Description Retrieve the verses of a song based on its ID with pagination.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songID path int true "Song ID"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of verses per page"
// @Success 200 {object} dto.VersesDto "List of song verses"
// @Failure 404 {object} delivery.JSONError "Not Found"
// @Failure 500 {object} delivery.JSONError "Internal Server Error"
// @Router /songs/{songID} [get]
func (h SongsHandler) getSongText(w http.ResponseWriter, r *http.Request, songID int, params dto.PaginationParamsDto) {
	verses, totalPages, err := h.songsService.GetSongText(int32(songID), params)
	if err != nil {
		log.WithError(err).Error(delivery.ErrGettingSongText)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JSONError{Error: delivery.ErrGettingSongText, Message: domain.ErrSongNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JSONError{Error: delivery.ErrGettingSongText})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, dto.VersesDto{
		Verses:     verses,
		TotalPages: totalPages,
	})
}

// @Summary Delete a song by song ID
// @Description Delete a song from the database based on its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songID path int true "Song ID"
// @Success 200 "Song successfully deleted"
// @Failure 404 {object} delivery.JSONError "Not Found"
// @Failure 500 {object} delivery.JSONError "Internal Server Error"
// @Router /songs/{songID} [delete]
func (h SongsHandler) deleteSong(w http.ResponseWriter, r *http.Request, songID int) {
	if err := h.songsService.Delete(int32(songID)); err != nil {
		log.WithError(err).Error(delivery.ErrDeletingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JSONError{Error: delivery.ErrDeletingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JSONError{Error: delivery.ErrDeletingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, nil)
}

// @Summary Update a song by song ID
// @Description Update the details of an existing song based on its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songID path int true "Song ID"
// @Param body body dto.SongParamsDto true "Song details to update"
// @Success 200 {object} dto.SongDto "Updated song"
// @Failure 400 {object} delivery.JSONError "Bad Request"
// @Failure 404 {object} delivery.JSONError "Not Found"
// @Failure 500 {object} delivery.JSONError "Internal Server Error"
// @Router /songs/{songID} [put]
func (h SongsHandler) updateSong(w http.ResponseWriter, r *http.Request, songID int, updateSongInput dto.SongParamsDto) {
	song, err := h.songsService.Update(int32(songID), updateSongInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrUpdatingSong)

		if errors.Is(err, domain.ErrSongNotFound) {
			delivery.RespondWithJSON(w, http.StatusNotFound, delivery.JSONError{Error: delivery.ErrUpdatingSong, Message: domain.ErrSongNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrSongAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrUpdatingSong, Message: domain.ErrSongAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JSONError{Error: delivery.ErrUpdatingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, h.toSongDto(song))
}

// @Summary Create a new song
// @Description Add a new song to the database.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param body body dto.CreateSongDto true "Song details to create"
// @Success 201 {object} dto.SongDto "Created song"
// @Failure 400 {object} delivery.JSONError "Bad Request"
// @Failure 500 {object} delivery.JSONError "Internal Server Error"
// @Router /songs [post]
func (h SongsHandler) createSong(w http.ResponseWriter, r *http.Request, createSongInput dto.CreateSongDto) {
	song, err := h.songsService.Create(createSongInput)
	if err != nil {
		log.WithError(err).Error(delivery.ErrCreatingSong)

		if errors.Is(err, domain.ErrSongAlreadyExist) {
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrCreatingSong, Message: domain.ErrSongAlreadyExist.Error()})
			return
		}

		delivery.RespondWithJSON(w, http.StatusInternalServerError, delivery.JSONError{Error: delivery.ErrCreatingSong})
		return
	}

	delivery.RespondWithJSON(w, http.StatusCreated, h.toSongDto(song))
}

func (h SongsHandler) toSongsDto(songs []domain.Song, totalPages int) dto.SongsDto {
	songsDto := make([]dto.SongDto, 0)

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
