package dto

// VersesDto represents the data transfer object for a collection of verses and total verse count.
type VersesDto struct {
	Verses      []string `json:"verses"`
	TotalVerses int      `json:"total_verses"`
}
