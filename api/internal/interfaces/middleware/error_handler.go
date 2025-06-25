package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse エラーレスポンス構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HandleAuthError 認証エラーのハンドリング
func HandleAuthError(c *gin.Context, err error) {
	status := StatusUnauthorized
	errorCode := ErrCodeInvalidToken
	message := ErrMsgInvalidToken

	if strings.Contains(err.Error(), "revoked") {
		errorCode = ErrCodeTokenRevoked
		message = ErrMsgTokenRevoked
	} else if strings.Contains(err.Error(), "expired") {
		errorCode = ErrCodeTokenExpired
		message = ErrMsgTokenExpired
	}

	c.JSON(status, ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
	c.Abort()
}

// HandleUnauthorizedError 未認証エラーのハンドリング
func HandleUnauthorizedError(c *gin.Context, message string) {
	c.JSON(StatusUnauthorized, ErrorResponse{
		Error:   ErrCodeUnauthorized,
		Message: message,
	})
	c.Abort()
}

// HandleForbiddenError 認可エラーのハンドリング
func HandleForbiddenError(c *gin.Context, message string) {
	c.JSON(StatusForbidden, ErrorResponse{
		Error:   ErrCodeInsufficientPerms,
		Message: message,
	})
	c.Abort()
}