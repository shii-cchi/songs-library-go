package dto

// GetSongsDto represents the data transfer object for retrieving songs with filters and pagination.
type GetSongsDto struct {
	Filters          SongParamsDto       `validate:"required"`
	PaginationParams PaginationParamsDto `validate:"required"`
}
