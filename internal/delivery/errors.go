package delivery

// JSONError represents the structure for error responses in JSON format.
type JSONError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Error constants for various input validation and parsing issues.
const (
	ErrInvalidPaginationParam = "invalid pagination param"
	ErrParsingParam           = "error parsing pagination param from string to int"
	ErrInvalidFilters         = "invalid filters param"
	ErrInvalidFilter          = "invalid filter param"
	ErrInvalidGetSongsParam   = "invalid get songs param"
	ErrInvalidIDInput         = "invalid song id input"
	ErrInvalidUpdateSongInput = "invalid update song input body"
	ErrInvalidJSON            = "invalid JSON body"
	ErrInvalidCreateSongInput = "invalid create song input body"
)

// Error constants for song-related operations.
const (
	ErrGettingSongs    = "error getting songs"
	ErrGettingSongText = "error getting song text"
	ErrDeletingSong    = "error deleting song"
	ErrUpdatingSong    = "error updating song"
	ErrCreatingSong    = "error create new song"
)
