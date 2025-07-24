package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
} 