package delivery

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

const (
	ErrInvalidCreateSongInput = "invalid create song input body"
	ErrInvalidUpdateSongInput = "invalid update song input body"
	ErrInvalidIDInput         = "invalid song id input"
	ErrInvalidJSON            = "invalid JSON body"
)

const (
	ErrCreatingSong = "error adding new song"
	ErrUpdatingSong = "error updating song info"
	ErrDeletingSong = "error deleting song"
)
