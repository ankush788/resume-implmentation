package errors

import (
	stderrs "errors"
	"net/http"

	"gorm.io/gorm"
)

// MapDBError maps common DB errors to AppError. Returns nil if input err is nil.
func MapDBError(err error) *AppError {
    if err == nil {
        return nil
    }

    // Record not found -> Not Found response
    if stderrs.Is(err, gorm.ErrRecordNotFound) {
        return &AppError{
            Code:       3001,
            Message:    "Resource not found",
            StatusCode: http.StatusNotFound,
        }
    }

    // Default -> Database internal error
    return NewInternalError(ErrDatabaseError.Message, err.Error())
}
