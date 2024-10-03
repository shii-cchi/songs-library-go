package delivery

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

const (
	ErrInvalidPaginationParam = "invalid pagination param"
	ErrParsingParam           = "error parsing pagination param from string to int"
	ErrInvalidFilters         = "invalid filters param"
	ErrInvalidGetSongsParam   = "invalid get songs param"
	ErrInvalidIDInput         = "invalid song id input"
	ErrInvalidUpdateSongInput = "invalid update song input body"
	ErrInvalidJSON            = "invalid JSON body"
	ErrInvalidCreateSongInput = "invalid create song input body"
)

const (
	ErrGettingSongs    = "error getting songs"
	ErrGettingSongText = "error getting song text"
	ErrDeletingSong    = "error deleting song"
	ErrUpdatingSong    = "error updating song"
	ErrCreatingSong    = "error create new song"
)
