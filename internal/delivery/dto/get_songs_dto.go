package dto

type GetSongsDto struct {
	Filters          SongParamsDto       `validate:"required"`
	PaginationParams PaginationParamsDto `validate:"required"`
}
