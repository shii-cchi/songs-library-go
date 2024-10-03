package domain

import (
	"database/sql"
	"time"
)

// Song represents the data model for a song.
type Song struct {
	ID          int32     `db:"id"`
	Group       string    `db:"group"`
	Song        string    `db:"song"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
}

// SongWithNull represents the data model for a song with nullable fields to handle optional details.
type SongWithNull struct {
	ID          int32          `db:"id"`
	Group       string         `db:"group"`
	Song        string         `db:"song"`
	ReleaseDate sql.NullTime   `db:"release_date"`
	Text        sql.NullString `db:"text"`
	Link        sql.NullString `db:"link"`
}
