package dto

// CreateSongDto represents the data transfer object for creating a new song.
type CreateSongDto struct {
	Group string `json:"group" validate:"required,max=100"`
	Song  string `json:"song" validate:"required,max=100"`
}
