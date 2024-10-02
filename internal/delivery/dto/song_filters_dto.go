package dto

type SongFiltersDto struct {
	Group       *string `validate:"omitempty,min=1,max=100"`
	Song        *string `validate:"omitempty,min=1,max=100"`
	ReleaseDate *string `validate:"omitempty,customDate"`
	Text        *string `validate:"omitempty,min=1,max=10000"`
	Link        *string `validate:"omitempty,url"`
}
