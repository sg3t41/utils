package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse エラーレスポンスの構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse 成功レスポンスの構造体
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// respondWithError エラーレスポンスを返す
func respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error: message,
	})
}

// respondWithSuccess 成功レスポンスを返す
func respondWithSuccess(c *gin.Context, message string, data ...interface{}) {
	response := SuccessResponse{
		Message: message,
	}
	
	if len(data) > 0 {
		response.Data = data[0]
	}
	
	c.JSON(http.StatusOK, response)
}

// respondBadRequest 400エラーを返す
func respondBadRequest(c *gin.Context, message string) {
	respondWithError(c, http.StatusBadRequest, message)
}

// respondInternalError 500エラーを返す
func respondInternalError(c *gin.Context, message string) {
	respondWithError(c, http.StatusInternalServerError, message)
}