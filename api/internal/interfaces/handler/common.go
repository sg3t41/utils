package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 共通のレスポンス処理

// SendErrorResponse 統一エラーレスポンス
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// ValidationErrorResponse バリデーションエラーレスポンス
func ValidationErrorResponse(c *gin.Context) {
	SendErrorResponse(c, http.StatusBadRequest, "バリデーションが実行されていません")
}

// ParseIDParam IDパラメータの解析
func ParseIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		SendErrorResponse(c, http.StatusBadRequest, "IDが指定されていません")
		return "", false
	}
	return id, true
}

// ParsePaginationParams ページネーションパラメータの解析
func ParsePaginationParams(c *gin.Context) (limit, offset int) {
	limit = 10 // デフォルト値
	offset = 0 // デフォルト値

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}

// GetValidatedBody バリデート済みリクエストボディの取得
func GetValidatedBody[T any](c *gin.Context) (*T, bool) {
	req, exists := c.Get("validated_body")
	if !exists {
		ValidationErrorResponse(c)
		return nil, false
	}

	typedReq, ok := req.(*T)
	if !ok {
		SendErrorResponse(c, http.StatusBadRequest, "リクエスト形式が不正です")
		return nil, false
	}

	return typedReq, true
}