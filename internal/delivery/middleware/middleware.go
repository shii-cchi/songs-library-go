package middleware

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"songs-library-go/internal/delivery"
	"songs-library-go/internal/delivery/dto"
	"strconv"
	"strings"
)

func ValidateGetSongsParam(v *validator.Validate, next func(http.ResponseWriter, *http.Request, delivery.PaginationParams, map[string]string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := validatePage(w, r)
		if err != nil {
			return
		}

		limit, err := validateLimit(w, r, delivery.DefaultSongsLimit)
		if err != nil {
			return
		}

		filters, err := validateFilters(w, r, v)
		if err != nil {
			return
		}

		next(w, r, delivery.PaginationParams{
			Page:  page,
			Limit: limit,
		}, filters)
	}
}

func ValidateGetSongParam(next func(http.ResponseWriter, *http.Request, int, delivery.PaginationParams)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		page, err := validatePage(w, r)
		if err != nil {
			return
		}

		limit, err := validateLimit(w, r, delivery.DefaultVerseLimit)
		if err != nil {
			return
		}

		next(w, r, songID, delivery.PaginationParams{
			Page:  page,
			Limit: limit,
		})
	}
}

func ValidateIDInput(next func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		next(w, r, songID)
	}
}

func ValidateUpdateSongInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, int, dto.UpdateSongDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		var updateSongInput dto.UpdateSongDto

		if err := json.NewDecoder(r.Body).Decode(&updateSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.ErrInvalidJSON})
			return
		}

		if !isAnyFieldProvided(updateSongInput) {
			log.Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.MesEmptyUpdateSongInput})
			return
		}

		if err := v.Struct(updateSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.MesInvalidUpdateSongInput})
			return
		}

		next(w, r, songID, updateSongInput)
	}
}

func ValidateCreateSongInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, dto.CreateSongDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createSongInput dto.CreateSongDto

		if err := json.NewDecoder(r.Body).Decode(&createSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidCreateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidCreateSongInput, Message: delivery.ErrInvalidJSON})
			return
		}

		if err := v.Struct(createSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidCreateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidCreateSongInput, Message: delivery.MesInvalidCreateSongInput})
			return
		}

		next(w, r, createSongInput)
	}
}

func validatePage(w http.ResponseWriter, r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: delivery.MesInvalidPage})
			return 0, errors.New(delivery.MesInvalidPage)
		}
		return page, nil
	}

	return delivery.DefaultPage, nil
}

func validateLimit(w http.ResponseWriter, r *http.Request, defaultLimit int) (int, error) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > delivery.LimitSongsPerPage {
			log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: delivery.MesInvalidLimit})
			return 0, errors.New(delivery.MesInvalidLimit)
		}
		return limit, nil
	}

	return defaultLimit, nil
}

func validateFilters(w http.ResponseWriter, r *http.Request, v *validator.Validate) (map[string]string, error) {
	validFilters := map[string]bool{
		"group":        true,
		"song":         true,
		"release_date": true,
		"text":         true,
		"link":         true,
	}

	filters := make(map[string]string)
	var dtoFilters dto.SongFiltersDto
	anyFieldExists := false

	for filter, isValid := range validFilters {
		if value := r.URL.Query().Get(filter); value != "" {
			if !isValid {
				log.Error(delivery.ErrInvalidGetSongsParam)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: delivery.MesInvalidFilterName})
				return nil, errors.New(delivery.MesInvalidFilterName)
			}

			filters[filter] = value

			if filter == "release_date" {
				filter = "releaseDate"
			}

			reflect.ValueOf(&dtoFilters).Elem().FieldByName(strings.Title(filter)).Set(reflect.ValueOf(&value))
			anyFieldExists = true
		}
	}

	if anyFieldExists {
		if err := v.Struct(dtoFilters); err != nil {
			log.Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: delivery.MesInvalidFilters})
			return nil, errors.New(delivery.MesInvalidFilters)
		}
	}

	return filters, nil
}

func extractAndValidateID(w http.ResponseWriter, r *http.Request) (int, error) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil || songID <= 0 {
		log.WithError(err).Error(delivery.ErrInvalidIDInput)
		delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidIDInput, Message: delivery.MesInvalidIDInput})
		return 0, errors.New(delivery.MesInvalidIDInput)
	}

	return songID, nil
}

func isAnyFieldProvided(input dto.UpdateSongDto) bool {
	return input.Group != nil || input.Song != nil || input.ReleaseDate != nil || input.Text != nil || input.Link != nil
}
