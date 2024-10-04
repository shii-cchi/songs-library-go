package dto

// SongDto represents the data transfer object for a song with its details.
type SongDto struct {
	ID          int32  `json:"id" example:"1"`
	Group       string `json:"group" example:"Rammstein"`
	Song        string `json:"song" example:"Weit Weg"`
	ReleaseDate string `json:"release_date,omitempty" example:"17.05.2019"`
	Text        string `json:"text,omitempty" example:"Niemand kann das Bild beschreiben\nGegen seine Fensterscheibe\nHat er das Gesicht gepresst\nUnd hofft, dass sie das Licht anlässt\nOhne Kleid sah er sie nie\nDie Herrin seiner Fantasie\nEr nimmt die Gläser vom Gesicht\nSingt zitternd eine Melodie\n\nDer Raum wird sich mit Mondlicht füllen\nLässt sie fallen, alle Hüllen\n\n"`
	Link        string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=N9AalJuwLyQ&ab_channel=Rammstein-Topic"`
}
