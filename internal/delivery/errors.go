package delivery

import "errors"

const (
	ErrInvalidCreateSongInput = "invalid create song input body"
	ErrInvalidUpdateSongInput = "invalid update song input body"
	ErrInvalidIDInput         = "invalid song id input"
	ErrInvalidJSON            = "invalid JSON body"
	ErrInvalidGetSongsParam   = "invalid get songs param"
	ErrInvalidPaginationParam = "invalid pagination param"
	ErrInvalidFilters         = "invalid filters param"
)

const (
	ErrCreatingSong = "error adding new song"
	ErrUpdatingSong = "error updating song info"
	ErrDeletingSong = "error deleting song"
	ErrGettingSongs = "error getting songs"
	ErrGettingSong  = "error getting song"
)

var (
	ErrParsingParam      = errors.New("error parsing pagination param from string to int")
	ErrInvalidFilterName = errors.New("invalid filter name")
)
