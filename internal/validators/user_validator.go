package validators

import (
	"regexp"
	"strings"

	"go-user-service/internal/errors"
)

// ValidateUsername validates the username field
func ValidateUsername(username string) *errors.AppError {
	username = strings.TrimSpace(username)

	if username == "" {
		return errors.ErrMissingUsername
	}

	if len(username) < 3 || len(username) > 50 {
		return errors.ErrInvalidUsername
	}

	// Only alphanumeric and underscores allowed
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return &errors.AppError{
			Code:       1004,
			Message:    "Username can only contain letters, numbers, and underscores",
			StatusCode: 400,
		}
	}

	return nil
}

// ValidatePassword validates the password field
func ValidatePassword(password string) *errors.AppError {
	if password == "" {
		return errors.ErrMissingPassword
	}

	if len(password) < 6 {
		return errors.ErrInvalidPassword
	}

	if len(password) > 128 {
		return &errors.AppError{
			Code:       1005,
			Message:    "Password must not exceed 128 characters",
			StatusCode: 400,
		}
	}

	return nil
}

// ValidateRegisterRequest validates the registration request
func ValidateRegisterRequest(username, password string) *errors.AppError {
	if err := ValidateUsername(username); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}

// ValidateLoginRequest validates the login request
func ValidateLoginRequest(username, password string) *errors.AppError {
	if username == "" {
		return errors.ErrMissingUsername
	}

	if password == "" {
		return errors.ErrMissingPassword
	}

	return nil
}
