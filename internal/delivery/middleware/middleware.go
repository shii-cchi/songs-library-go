package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
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

// ValidateGetSongsParam validates pagination and filter parameters for getting songs.
func ValidateGetSongsParam(v *validator.Validate, next func(http.ResponseWriter, *http.Request, dto.GetSongsDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := getPaginationParam(w, r, "page", delivery.DefaultPage)
		if err != nil {
			return
		}

		limit, err := getPaginationParam(w, r, "limit", delivery.DefaultSongsLimit)
		if err != nil {
			return
		}

		filters, err := getFilters(w, r)
		if err != nil {
			return
		}

		getSongsDto := dto.GetSongsDto{
			Filters: filters,
			PaginationParams: dto.PaginationParamsDto{
				Page:  page,
				Limit: limit,
			},
		}

		if err := v.Struct(getSongsDto); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidGetSongsParam, Message: delivery.MesInvalidGetSongsParam})
			return
		}

		next(w, r, getSongsDto)
	}
}

// ValidateGetSongParam validates the song ID and pagination parameters for retrieving a specific song.
func ValidateGetSongParam(v *validator.Validate, next func(http.ResponseWriter, *http.Request, int, dto.PaginationParamsDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		page, err := getPaginationParam(w, r, "page", delivery.DefaultPage)
		if err != nil {
			return
		}

		limit, err := getPaginationParam(w, r, "limit", delivery.DefaultVerseLimit)
		if err != nil {
			return
		}

		paginationParams := dto.PaginationParamsDto{
			Page:  page,
			Limit: limit,
		}

		if err := v.Struct(paginationParams); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidPaginationParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidPaginationParam, Message: delivery.MesInvalidPaginationParam})
			return
		}

		next(w, r, songID, paginationParams)
	}
}

// ValidateIDInput validates the song ID extracted from the request for further processing.
func ValidateIDInput(next func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		next(w, r, songID)
	}
}

// ValidateUpdateSongInput validates the song ID and update input for modifying a song's details.
func ValidateUpdateSongInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, int, dto.SongParamsDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songID, err := extractAndValidateID(w, r)
		if err != nil {
			return
		}

		var updateSongInput dto.SongParamsDto

		if err := json.NewDecoder(r.Body).Decode(&updateSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.ErrInvalidJSON})
			return
		}

		if !isAnyFieldProvided(updateSongInput) {
			log.Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.MesEmptyUpdateSongInput})
			return
		}

		trimSpace(&updateSongInput)

		if err := v.Struct(updateSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidUpdateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.MesInvalidUpdateSongInput})
			return
		}

		next(w, r, songID, updateSongInput)
	}
}

// ValidateCreateSongInput validates the input for creating a new song.
func ValidateCreateSongInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, dto.CreateSongDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createSongInput dto.CreateSongDto

		if err := json.NewDecoder(r.Body).Decode(&createSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidCreateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidCreateSongInput, Message: delivery.ErrInvalidJSON})
			return
		}

		trimSpace(&createSongInput)

		if err := v.Struct(createSongInput); err != nil {
			log.WithError(err).Error(delivery.ErrInvalidCreateSongInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidCreateSongInput, Message: delivery.MesInvalidCreateSongInput})
			return
		}

		next(w, r, createSongInput)
	}
}

func getPaginationParam(w http.ResponseWriter, r *http.Request, paramName string, defaultValue int) (int, error) {
	paramStr := r.URL.Query().Get(paramName)
	if paramStr != "" {
		paramValue, err := strconv.Atoi(paramStr)
		if err != nil {
			log.WithError(err).Error(delivery.ErrInvalidPaginationParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidPaginationParam, Message: fmt.Sprintf("%s (param: %s, value: %s)", delivery.ErrParsingParam, paramName, paramStr)})
			return 0, errors.New(delivery.ErrParsingParam)
		}
		return paramValue, nil
	}

	return defaultValue, nil
}

func getFilters(w http.ResponseWriter, r *http.Request) (dto.SongParamsDto, error) {
	validFilters := map[string]bool{
		"group":        true,
		"song":         true,
		"release_date": true,
		"text":         true,
		"link":         true,
	}

	var dtoFilters dto.SongParamsDto

	for filter, values := range r.URL.Query() {
		if filter == "page" || filter == "limit" {
			continue
		}

		if !validFilters[filter] {
			log.Error(fmt.Sprintf("%s (filter name: %s)", delivery.ErrInvalidFilters, filter))
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidFilters, Message: delivery.MesInvalidFilterName})
			return dto.SongParamsDto{}, errors.New(delivery.ErrInvalidFilters)
		}

		value := values[0]

		if value == "" {
			log.Error(fmt.Sprintf("%s (filter name: %s)", delivery.ErrInvalidFilter, filter))
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidFilter, Message: delivery.MesEmptyFilter})
			return dto.SongParamsDto{}, errors.New(delivery.ErrInvalidFilter)
		}

		if filter == "release_date" {
			filter = "releaseDate"
		}

		reflect.ValueOf(&dtoFilters).Elem().FieldByName(strings.Title(filter)).Set(reflect.ValueOf(&value))
	}

	trimSpace(&dtoFilters)

	return dtoFilters, nil
}

func extractAndValidateID(w http.ResponseWriter, r *http.Request) (int, error) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil || songID <= 0 {
		log.WithError(err).Error(delivery.ErrInvalidIDInput)
		delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JSONError{Error: delivery.ErrInvalidIDInput, Message: delivery.MesInvalidIDInput})
		return 0, errors.New(delivery.MesInvalidIDInput)
	}

	return songID, nil
}

func isAnyFieldProvided(input dto.SongParamsDto) bool {
	return input.Group != nil || input.Song != nil || input.ReleaseDate != nil || input.Text != nil || input.Link != nil
}

func trimSpace(input interface{}) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			field.SetString(strings.TrimSpace(field.String()))
		}

		if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String {
			if !field.IsNil() {
				strVal := field.Elem()
				strVal.SetString(strings.TrimSpace(strVal.String()))
			}
		}
	}
}
