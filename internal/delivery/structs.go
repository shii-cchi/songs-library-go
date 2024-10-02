package delivery

type JsonError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type PaginationParams struct {
	Page  int
	Limit int
}
