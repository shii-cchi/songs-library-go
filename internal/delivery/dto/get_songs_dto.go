package dto

type GetSongsDto struct {
	Filters          SongParamsDto       `validate:"required,dive"`
	PaginationParams PaginationParamsDto `validate:"required,dive"`
}
