package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"reflect"
	"songs-library-go/internal/delivery"
	"songs-library-go/internal/delivery/dto"
	"strconv"
	"strings"
)

func ValidateCreateSongInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			ctx := context.WithValue(r.Context(), delivery.CreateSongInputKey, createSongInput)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateUpdateSongInput(v *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			ctx := context.WithValue(r.Context(), delivery.UpdateSongInputKey, updateSongInput)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateIDInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		songIDStr := chi.URLParam(r, "id")
		songID, err := strconv.Atoi(songIDStr)
		if err != nil || songID <= 0 {
			log.WithError(err).Error(delivery.ErrInvalidIDInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidIDInput, Message: delivery.MesInvalidIDInput})
			return
		}

		ctx := context.WithValue(r.Context(), delivery.IDInputKey, int32(songID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateGetSongsParam(v *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			query := r.URL.Query()

			page, err := validatePage(query)
			if err != nil {
				log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: err.Error()})
				return
			}
			ctx = context.WithValue(ctx, delivery.PageKey, page)

			limit, err := validateLimit(query, delivery.DefaultSongsLimit)
			if err != nil {
				log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: err.Error()})
				return
			}
			ctx = context.WithValue(ctx, delivery.LimitKey, limit)

			filters, err := validateFilters(query, v)
			if err != nil {
				log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: err.Error()})
				return
			}
			ctx = context.WithValue(ctx, delivery.FiltersKey, filters)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateGetSongParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := r.URL.Query()

		page, err := validatePage(query)
		if err != nil {
			log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: err.Error()})
			return
		}
		ctx = context.WithValue(ctx, delivery.PageKey, page)

		limit, err := validateLimit(query, delivery.DefaultVerseLimit)
		if err != nil {
			log.WithError(err).Error(delivery.ErrInvalidGetSongsParam)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidGetSongsParam, Message: err.Error()})
			return
		}
		ctx = context.WithValue(ctx, delivery.LimitKey, limit)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isAnyFieldProvided(input dto.UpdateSongDto) bool {
	return input.Group != nil || input.Song != nil || input.ReleaseDate != nil || input.Text != nil || input.Link != nil
}

func validateLimit(query url.Values, defaultLimit int) (int, error) {
	limitStr := query.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > delivery.LimitSongsPerPage {
			return 0, errors.New(delivery.MesInvalidLimit)
		}
		return limit, nil
	}

	return defaultLimit, nil
}

func validatePage(query url.Values) (int, error) {
	pageStr := query.Get("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return 0, errors.New(delivery.MesInvalidPage)
		}
		return page, nil
	}

	return delivery.DefaultPage, nil
}

func validateFilters(query url.Values, v *validator.Validate) (map[string]string, error) {
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
		if value := query.Get(filter); value != "" {
			if !isValid {
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
			return nil, errors.New(delivery.MesInvalidFilters)
		}
	}

	return filters, nil
}
