package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response struct untuk format response API yang konsisten
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse untuk response sukses
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// CreatedResponse untuk response ketika data berhasil dibuat
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse untuk response error
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
	c.JSON(statusCode, response)
}

// BadRequestError untuk error 400 Bad Request
func BadRequestError(c *gin.Context, message string, err interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// NotFoundError untuk error 404 Not Found
func NotFoundError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerError untuk error 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err interface{}) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}

// ValidationError untuk error validasi
func ValidationError(c *gin.Context, err interface{}) {
	ErrorResponse(c, http.StatusBadRequest, "Validation failed", err)
}
