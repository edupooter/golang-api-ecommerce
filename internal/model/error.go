package model

// ErrorResponse is a standard envelope for error responses
type ErrorResponse struct {
	// Message short human-readable error message
	Message string `json:"message" example:"resource not found"`
}
