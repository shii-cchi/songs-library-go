package delivery

const (
	DefaultPage       = 1
	DefaultSongsLimit = 10
	DefaultVerseLimit = 2
	LimitSongsPerPage = 100
)

const (
	MesInvalidCreateSongInput = "fields group and song are required and must have at least 1 character and can have at most 100 characters"
	MesEmptyUpdateSongInput   = "at least one field must be provided for update"
	MesInvalidUpdateSongInput = "fields group and song must have at least 1 character and can have at most 100 characters, field release_date must be a valid date in the format `dd.mm.yyyy`, field text must have at least 1 character and can have at most 10,000 characters, field link must be a valid URL"
	MesInvalidIDInput         = "id must be a positive integer"
	MesInvalidGetSongsParam   = "page must be a positive integer, limit must be a positive integer and can't be greater than 100, group and song must have at least 1 character and can have at most 100 characters, field release_date must be a valid date in the format `dd.mm.yyyy`, field text must have at least 1 character and can have at most 100 characters, field link must be a valid URL"
	MesInvalidPaginationParam = "page must be a positive integer, limit must be a positive integer and can't be greater than 100"
	MesInvalidFilterName      = "filters can be only group, song, release_date, text or link"
)
