package dto

type SongsDto struct {
	Songs      []SongDto `json:"songs"`
	TotalPages int       `json:"total_pages"`
}
