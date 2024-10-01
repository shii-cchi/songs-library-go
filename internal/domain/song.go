package domain

import (
	"database/sql"
	"time"
)

type Song struct {
	ID          int32     `db:"id"`
	Group       string    `db:"group"`
	Song        string    `db:"song"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
}

type SongWithNull struct {
	ID          int32          `db:"id"`
	Group       string         `db:"group"`
	Song        string         `db:"song"`
	ReleaseDate sql.NullTime   `db:"release_date"`
	Text        sql.NullString `db:"text"`
	Link        sql.NullString `db:"link"`
}

type AddDetailsParams struct {
	ID          int32
	ReleaseDate *time.Time
	Text        *string
	Link        *string
}

type UpdateParams struct {
	Details AddDetailsParams
	Group   *string
	Song    *string
}
