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
	songExists, err := update.Executor().ScanStruct(&updatedSong)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == domain.CodeUniqueConstraintViolation {
			return domain.Song{}, fmt.Errorf("%w: %s", domain.ErrSongAlreadyExist, err)
		}
		return domain.Song{}, err
	}

	if !songExists {
		return domain.Song{}, domain.ErrSongNotFound
	}

	return r.toSong(updatedSong), nil
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

func (r SongsRepo) GetSongs(page int, limit int, filtersMap map[string]string) ([]domain.Song, error) {
	query := r.goquDb.From(songsTable)

	conditions := goqu.Ex{}
	if filtersMap != nil && len(filtersMap) > 0 {
		for field, value := range filtersMap {
			conditions[field] = value
		}

		query = query.Where(conditions)
	}

	totalCount, err := r.getTotalCount(conditions)
	if err != nil {
		return nil, err
	}

	totalPages := totalCount/limit + 1
	if page > totalPages {
		return nil, fmt.Errorf("%w (page: %d, total pages: %d)", domain.ErrPageDoesntExist, page, totalPages)
	}

	query = query.Limit(uint(limit)).Offset(uint((page - 1) * limit))

	var songs []domain.SongWithNull
	if err := query.Executor().ScanStructs(&songs); err != nil {
		return nil, err
	}

	return r.toSongs(songs), nil
}

func (r SongsRepo) GetSong(songID int32) (string, error) {
	query := r.goquDb.Select("text").From(songsTable).Where(goqu.Ex{"id": songID})

	var text sql.NullString
	songExists, err := query.Executor().ScanVal(&text)
	if err != nil {
		return "", err
	}

	if !songExists {
		return "", domain.ErrSongNotFound
	}

	if !text.Valid {
		return "", domain.ErrSongTextNotFound
	}

	return text.String, nil
}

func (r SongsRepo) toSong(song domain.SongWithNull) domain.Song {
	normalizedSong := domain.Song{
		ID:    song.ID,
		Group: song.Group,
		Song:  song.Song,
	}

	if song.ReleaseDate.Valid {
		normalizedSong.ReleaseDate = song.ReleaseDate.Time
	}
	if song.Text.Valid {
		normalizedSong.Text = song.Text.String
	}
	if song.Link.Valid {
		normalizedSong.Link = song.Link.String
	}

	return normalizedSong
}

func (r SongsRepo) toSongs(songs []domain.SongWithNull) []domain.Song {
	normalizedSongs := make([]domain.Song, len(songs))
	for i, song := range songs {
		normalizedSong := domain.Song{
			ID:    song.ID,
			Group: song.Group,
			Song:  song.Song,
		}

		if song.ReleaseDate.Valid {
			normalizedSong.ReleaseDate = song.ReleaseDate.Time
		}
		if song.Text.Valid {
			normalizedSong.Text = song.Text.String
		}
		if song.Link.Valid {
			normalizedSong.Link = song.Link.String
		}

		normalizedSongs[i] = normalizedSong
	}

	return normalizedSongs
}

func (r SongsRepo) getTotalCount(conditions goqu.Ex) (int, error) {
	var totalCount int

	query := r.goquDb.Select(goqu.COUNT("id")).From(songsTable).Where(conditions)

	if _, err := query.Executor().ScanVal(&totalCount); err != nil {
		return 0, err
	}

	return totalCount, nil
}
