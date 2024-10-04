package dto

// VersesDto represents the data transfer object for a collection of verses and total verse count.
type VersesDto struct {
	Verses     []string `json:"verses" example:"[\"Niemand kann das Bild beschreiben\\nGegen seine Fensterscheibe\\nHat er das Gesicht gepresst\\nUnd hofft, dass sie das Licht anlässt\\nOhne Kleid sah er sie nie\\nDie Herrin seiner Fantasie\\nEr nimmt die Gläser vom Gesicht\\nSingt zitternd eine Melodie\", \"Der Raum wird sich mit Mondlicht füllen\\nLässt sie fallen, alle Hüllen\"]"`
	TotalPages int      `json:"total_pages" example:"1"`
}
