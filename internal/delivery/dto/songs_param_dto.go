package dto

type SongParamsDto struct {
	Group       *string `json:"group,omitempty" validate:"omitempty,min=1,max=100"`
	Song        *string `json:"song,omitempty" validate:"omitempty,min=1,max=100"`
	ReleaseDate *string `json:"release_date,omitempty" validate:"omitempty,customDate"`
	Text        *string `json:"text,omitempty" validate:"omitempty,min=1,max=10000"`
	Link        *string `json:"link,omitempty" validate:"omitempty,url"`
}
