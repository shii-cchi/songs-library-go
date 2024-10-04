package dto

// GetSongsDto represents the data transfer object for retrieving songs with filters and pagination.
type GetSongsDto struct {
	Filters          SongParamsDto       `validate:"required" example:"{\"release_date\":\"2024-10-04\"}"`
	PaginationParams PaginationParamsDto `validate:"required" example:"{\"page\":1, \"limit\":10}"`
}
