package dto

type PaginationParamsDto struct {
	Page  int `validate:"required,gte=1"`
	Limit int `validate:"required,gte=1,lte=100"`
}
