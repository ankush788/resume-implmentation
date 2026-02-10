package middleware

import (
	"log"
	"net/http"

	"go-user-service/internal/errors"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the standardized error response structure
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}


//----- both under all of this is use for API layer (http format) error throwing ----//

// ErrorHandlerMiddleware is a middleware that catches panics and converts them to proper error responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    5001,
					Message: "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.StatusCode, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		})
		return
	}

	// For unexpected errors
	log.Printf("Unexpected error: %v", err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    5001,
		Message: "Internal server error",
	})
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
