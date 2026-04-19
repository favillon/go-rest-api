package response

import (
	apierrors "backend-productos/api/errors"

	"github.com/gin-gonic/gin"
)

// APIError representa el detalle tecnico de un error de la API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIResponse define un contrato uniforme para todas las respuestas.
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

func RespondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func RespondError(c *gin.Context, statusCode int, message string, code apierrors.Code, detail string) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Error: &APIError{
			Code:    string(code),
			Message: detail,
		},
	})
}
