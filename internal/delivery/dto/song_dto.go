package dto

// SongDto represents the data transfer object for a song with its details.
type SongDto struct {
	ID          int32  `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}
