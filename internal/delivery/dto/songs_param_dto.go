package dto

// SongParamsDto represents the data transfer object for optional song filter parameters.
type SongParamsDto struct {
	Group       *string `json:"group,omitempty" validate:"omitempty,min=1,max=100" example:"Rammstein"`
	Song        *string `json:"song,omitempty" validate:"omitempty,min=1,max=100" example:"Weit Weg"`
	ReleaseDate *string `json:"release_date,omitempty" validate:"omitempty,customDate" example:"17.05.2019"`
	Text        *string `json:"text,omitempty" validate:"omitempty,min=1,max=10000" example:"Niemand kann das Bild beschreiben\nGegen seine Fensterscheibe\nHat er das Gesicht gepresst\nUnd hofft, dass sie das Licht anlässt\nOhne Kleid sah er sie nie\nDie Herrin seiner Fantasie\nEr nimmt die Gläser vom Gesicht\nSingt zitternd eine Melodie\n\nDer Raum wird sich mit Mondlicht füllen\nLässt sie fallen, alle Hüllen\n\n"`
	Link        *string `json:"link,omitempty" validate:"omitempty,url" example:"https://www.youtube.com/watch?v=N9AalJuwLyQ&ab_channel=Rammstein-Topic"`
}
