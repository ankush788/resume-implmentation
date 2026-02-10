package errors

import "net/http"

// for service layer throwing error to api layer

func NewValidationError(message, details string) *AppError {
	return &AppError{
		Code:       1001,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

func NewInternalError(message, details string) *AppError {
	return &AppError{
		Code:       5001,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusInternalServerError,
	}
}