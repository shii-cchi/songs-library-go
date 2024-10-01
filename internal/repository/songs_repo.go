package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"reflect"
	"songs-library-go/internal/domain"
)

const songsTable = "songs"

type SongsRepo struct {
	goquDb *goqu.Database
}

func NewSongsRepo(db *sql.DB) *SongsRepo {
	return &SongsRepo{
		goquDb: goqu.New("postgres", db),
	}
}

func (r SongsRepo) Create(groupName, songName string) (domain.Song, error) {
	var newSong domain.Song

	insert := r.goquDb.Insert(songsTable).
		Rows(goqu.Record{"group": groupName, "song": songName}).
		Returning("id", "group", "song")

	_, err := insert.Executor().ScanStruct(&newSong)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == domain.CodeUniqueConstraintViolation {
			return domain.Song{}, fmt.Errorf("%w: %s", domain.ErrSongAlreadyExist, err)
		}
		return domain.Song{}, err
	}

	return newSong, nil
}

func (r SongsRepo) AddDetails(params domain.AddDetailsParams) error {
	fieldsToUpdate := map[string]interface{}{}

	for field, value := range map[string]interface{}{
		"release_date": params.ReleaseDate,
		"text":         params.Text,
		"link":         params.Link,
	} {
		if value != nil {
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Ptr && !val.IsNil() {
				fieldsToUpdate[field] = val.Elem().Interface()
			}
		}
	}

	update := r.goquDb.Update(songsTable).
		Set(fieldsToUpdate).
		Where(goqu.Ex{"id": params.ID})

	res, err := update.Executor().Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSongNotFound
	}

	return nil
}

func (r SongsRepo) UpdateSong(params domain.UpdateParams) (domain.Song, error) {
	fieldsToUpdate := map[string]interface{}{}

	for field, value := range map[string]interface{}{
		"group":        params.Group,
		"song":         params.Song,
		"release_date": params.Details.ReleaseDate,
		"text":         params.Details.Text,
		"link":         params.Details.Link,
	} {
		if value != nil {
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Ptr && !val.IsNil() {
				fieldsToUpdate[field] = val.Elem().Interface()
			}
		}
	}

	update := r.goquDb.Update(songsTable).
		Set(fieldsToUpdate).
		Where(goqu.Ex{"id": params.Details.ID}).
		Returning("id", "group", "song", "release_date", "text", "link")

	var updatedSong domain.SongWithNull
	songExistence, err := update.Executor().ScanStruct(&updatedSong)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == domain.CodeUniqueConstraintViolation {
			return domain.Song{}, fmt.Errorf("%w: %s", domain.ErrSongAlreadyExist, err)
		}
		return domain.Song{}, err
	}

	if !songExistence {
		return domain.Song{}, domain.ErrSongNotFound
	}

	normalizedSong := domain.Song{
		ID:    updatedSong.ID,
		Group: updatedSong.Group,
		Song:  updatedSong.Song,
	}

	if updatedSong.ReleaseDate.Valid {
		normalizedSong.ReleaseDate = updatedSong.ReleaseDate.Time
	}
	if updatedSong.Text.Valid {
		normalizedSong.Text = updatedSong.Text.String
	}
	if updatedSong.Link.Valid {
		normalizedSong.Link = updatedSong.Link.String
	}

	return normalizedSong, nil
}

func (r SongsRepo) Delete(songID int32) error {
	de := r.goquDb.Delete(songsTable).Where(goqu.Ex{"id": songID})

	res, err := de.Executor().Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSongNotFound
	}

	return nil
}
