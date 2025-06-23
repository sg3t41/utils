package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/errors"
)

// StandardErrorResponse represents the enhanced error response format
type StandardErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// HandleError handles domain errors and maps them to appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.DomainError:
		handleDomainError(c, e)
	default:
		// Handle unknown errors as internal server errors
		c.JSON(http.StatusInternalServerError, StandardErrorResponse{
			Error: "内部サーバーエラーが発生しました",
			Code:  "INTERNAL_SERVER_ERROR",
		})
	}
}

func handleDomainError(c *gin.Context, err *errors.DomainError) {
	switch err.Type {
	case errors.ErrorTypeNotFound:
		c.JSON(http.StatusNotFound, StandardErrorResponse{
			Error: err.Message,
			Code:  "NOT_FOUND",
		})
	case errors.ErrorTypeAlreadyExists:
		c.JSON(http.StatusConflict, StandardErrorResponse{
			Error: err.Message,
			Code:  "ALREADY_EXISTS",
		})
	case errors.ErrorTypeInvalidInput:
		c.JSON(http.StatusBadRequest, StandardErrorResponse{
			Error: err.Message,
			Code:  "INVALID_INPUT",
		})
	case errors.ErrorTypeUnauthorized:
		c.JSON(http.StatusUnauthorized, StandardErrorResponse{
			Error: err.Message,
			Code:  "UNAUTHORIZED",
		})
	case errors.ErrorTypeForbidden:
		c.JSON(http.StatusForbidden, StandardErrorResponse{
			Error: err.Message,
			Code:  "FORBIDDEN",
		})
	case errors.ErrorTypeInternalServer:
		c.JSON(http.StatusInternalServerError, StandardErrorResponse{
			Error: err.Message,
			Code:  "INTERNAL_SERVER_ERROR",
		})
	default:
		c.JSON(http.StatusInternalServerError, StandardErrorResponse{
			Error: "内部サーバーエラーが発生しました",
			Code:  "INTERNAL_SERVER_ERROR",
		})
	}
}

// HandleValidationError handles validation errors from gin-validator
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, StandardErrorResponse{
		Error: "バリデーションエラー",
		Code:  "VALIDATION_ERROR",
		Details: err.Error(),
	})
}