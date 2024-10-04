package dto

// SongsDto represents the data transfer object for a collection of songs and total page count.
type SongsDto struct {
	Songs      []SongDto `json:"songs"`
	TotalPages int       `json:"total_pages" example:"1"`
}
