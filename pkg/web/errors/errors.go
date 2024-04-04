package errors

import "net/http"

type APIError struct {
	Message string
	StatusCode int
}

func (e APIError) Error() string {
	return e.Message
}

func NewApiError(status int, message string) *APIError {
	return &APIError{
		StatusCode: status,
		Message: message,
	}
}

func NewForbiddenError(message string) *APIError {
	return NewApiError(http.StatusForbidden, message)
}

func NewValidationError(message string) *APIError {
	return NewApiError(http.StatusUnprocessableEntity, message)
}