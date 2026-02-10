package errors

import "net/http"

// AppError represents a standardized application error
type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"` // HTTP status code
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Custom error types
var (
	// Validation errors (400)
	ErrInvalidRequest = &AppError{
		Code:       1001,
		Message:    "Invalid request",
		StatusCode: http.StatusBadRequest,
	}

	ErrMissingUsername = &AppError{
		Code:       1002,
		Message:    "Username is required",
		StatusCode: http.StatusBadRequest,
	}

	ErrMissingPassword = &AppError{
		Code:       1003,
		Message:    "Password is required",
		StatusCode: http.StatusBadRequest,
	}

	ErrInvalidUsername = &AppError{
		Code:       1004,
		Message:    "Username must be between 3 and 50 characters",
		StatusCode: http.StatusBadRequest,
	}

	ErrInvalidPassword = &AppError{
		Code:       1005,
		Message:    "Password must be at least 6 characters",
		StatusCode: http.StatusBadRequest,
	}

	ErrUserAlreadyExists = &AppError{
		Code:       1006,
		Message:    "User already exists",
		StatusCode: http.StatusConflict,
	}

	// Authentication errors (401)
	ErrInvalidCredentials = &AppError{
		Code:       2001,
		Message:    "Invalid credentials",
		StatusCode: http.StatusUnauthorized,
	}

	ErrUserNotFound = &AppError{
		Code:       2002,
		Message:    "User not found",
		StatusCode: http.StatusUnauthorized,
	}

	// Not found errors (404)
	ErrTaskNotFound = &AppError{
		Code:       3001,
		Message:    "Task not found",
		StatusCode: http.StatusNotFound,
	}

	// Server errors (500)
	ErrInternalServer = &AppError{
		Code:       5001,
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrHashingPassword = &AppError{
		Code:       5002,
		Message:    "Error while hashing password",
		StatusCode: http.StatusInternalServerError,
	}

	ErrDatabaseError = &AppError{
		Code:       5003,
		Message:    "Database error",
		StatusCode: http.StatusInternalServerError,
	}
)





// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {  //"Try to convert err into *AppError type."
	appErr, ok := err.(*AppError)
	return appErr, ok
}
