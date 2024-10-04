package dto

// PaginationParamsDto represents the data transfer object for pagination parameters.
type PaginationParamsDto struct {
	Page  int `validate:"required,gte=1" example:"1"`
	Limit int `validate:"required,gte=1,lte=100" example:"3"`
}
