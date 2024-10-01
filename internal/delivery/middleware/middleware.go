package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"songs-library-go/internal/delivery"
	"songs-library-go/internal/delivery/dto"
	"strconv"
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
				delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidUpdateSongInput, Message: delivery.MesInvalidUpdateSongInput})
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
		if err != nil {
			log.WithError(err).Error(delivery.ErrInvalidIDInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidIDInput})
			return
		}

		if songID <= 0 {
			log.Error(delivery.ErrInvalidIDInput)
			delivery.RespondWithJSON(w, http.StatusBadRequest, delivery.JsonError{Error: delivery.ErrInvalidIDInput, Message: delivery.MesInvalidIDInput})
			return
		}

		ctx := context.WithValue(r.Context(), delivery.IDInputKey, int32(songID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isAnyFieldProvided(input dto.UpdateSongDto) bool {
	return input.Group != nil || input.Song != nil || input.ReleaseDate != nil || input.Text != nil || input.Link != nil
}
