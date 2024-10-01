package delivery

type ContextKey string

const (
	CreateSongInputKey ContextKey = "createSongInput"
	UpdateSongInputKey ContextKey = "updateSongInput"
	IDInputKey         ContextKey = "idInputKey"
)

const (
	MesInvalidCreateSongInput = "fields group and song are required and must have at least 1 character and can have at most 100 characters"
	MesInvalidUpdateSongInput = "at least one field must be provided for update. fields group and song must have at least 1 character and can have at most 100 characters, field release_date must be a valid date in the format `dd.mm.yyyy`, field text must have at least 1 character and can have at most 10,000 characters, field link must be a valid URL"
	MesInvalidIDInput         = "id must be a positive integer"
)
