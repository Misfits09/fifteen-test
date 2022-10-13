package structs

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Success     bool   `json:"success"`
	Description string `json:"description"`
}
