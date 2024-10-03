package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"math"
	"songs-library-go/internal/domain"
	"strings"
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

func (r SongsRepo) GetSongs(page int, limit int, filtersMap map[string]string) ([]domain.Song, int, error) {
	query := r.goquDb.From(songsTable)

	conditions := goqu.Ex{}

	if len(filtersMap) > 0 {
		for field, value := range filtersMap {
			if field == "text" && value != "" {
				words := strings.Fields(value)

				var textConditions []goqu.Expression
				for _, word := range words {
					textConditions = append(textConditions, goqu.I("text").ILike("%"+word+"%"))
				}

				query = query.Where(goqu.And(textConditions...))
			} else {
				conditions[field] = value
			}
		}

		query = query.Where(conditions)
	}

	totalCount, err := r.getTotalCount(conditions)
	if err != nil {
		return nil, 0, err
	}

	query = query.Limit(uint(limit)).Offset(uint((page - 1) * limit))

	var songs []domain.SongWithNull
	if err := query.Executor().ScanStructs(&songs); err != nil {
		return nil, 0, err
	}

	return r.toSongs(songs), int(math.Ceil(float64(totalCount) / float64(limit))), nil
}

func (r SongsRepo) GetSongText(songID int32) (string, error) {
	query := r.goquDb.Select("text").From(songsTable).Where(goqu.Ex{"id": songID})

	var text sql.NullString
	songExists, err := query.Executor().ScanVal(&text)
	if err != nil {
		return "", err
	}

	if !songExists {
		return "", fmt.Errorf("%w (id: %d)", domain.ErrSongNotFound, songID)
	}

	if !text.Valid {
		return "", nil
	}

	return text.String, nil
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
		return fmt.Errorf("%w (id: %d)", domain.ErrSongNotFound, songID)
	}

	return nil
}

func (r SongsRepo) UpdateSong(songID int32, paramsMap map[string]string) (domain.Song, error) {
	update := r.goquDb.Update(songsTable).
		Set(paramsMap).
		Where(goqu.Ex{"id": songID}).
		Returning("id", "group", "song", "release_date", "text", "link")

	var updatedSong domain.SongWithNull
	songExists, err := update.Executor().ScanStruct(&updatedSong)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == domain.CodeUniqueConstraintViolation {
			return domain.Song{}, fmt.Errorf("%w (id: %d): %s", domain.ErrSongAlreadyExist, songID, err)
		}
		return domain.Song{}, err
	}

	if !songExists {
		return domain.Song{}, fmt.Errorf("%w (id: %d)", domain.ErrSongNotFound, songID)
	}

	return r.toSong(updatedSong), nil
}

func (r SongsRepo) Create(groupName, songName string) (domain.Song, error) {
	insert := r.goquDb.Insert(songsTable).
		Rows(goqu.Record{"group": groupName, "song": songName}).
		Returning("id", "group", "song")

	var newSong domain.Song
	_, err := insert.Executor().ScanStruct(&newSong)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == domain.CodeUniqueConstraintViolation {
			return domain.Song{}, fmt.Errorf("%w (group name: %s, song name: %s): %s", domain.ErrSongAlreadyExist, groupName, songName, err)
		}
		return domain.Song{}, err
	}

	return newSong, nil
}

func (r SongsRepo) AddDetails(songID int32, paramsMap map[string]string) error {
	update := r.goquDb.Update(songsTable).
		Set(paramsMap).
		Where(goqu.Ex{"id": songID})

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

func (r SongsRepo) getTotalCount(conditions goqu.Ex) (int, error) {
	query := r.goquDb.Select(goqu.COUNT("id")).From(songsTable).Where(conditions)

	var totalCount int
	if _, err := query.Executor().ScanVal(&totalCount); err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r SongsRepo) toSongs(songs []domain.SongWithNull) []domain.Song {
	normalizedSongs := make([]domain.Song, len(songs))
	for i, song := range songs {
		normalizedSongs[i] = r.toSong(song)
	}

	return normalizedSongs
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
