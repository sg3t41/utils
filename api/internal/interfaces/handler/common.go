package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/domain/errors"
)

// 共通のレスポンス処理

// SendErrorResponse 統一エラーレスポンス
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// SendErrorWithCodeResponse エラーコード付きエラーレスポンス
func SendErrorWithCodeResponse(c *gin.Context, statusCode int, message string, code string) {
	c.JSON(statusCode, StandardErrorResponse{
		Error: message,
		Code:  code,
	})
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

// HandleUseCaseError ユースケースエラーの統一処理
func HandleUseCaseError(c *gin.Context, err error, resourceName string) {
	if err == nil {
		return
	}

	// まず特定のユースケースエラーをチェック
	switch err {
	case usecase.ErrUserNotFound:
		SendErrorResponse(c, http.StatusNotFound, fmt.Sprintf("%sが見つかりません", resourceName))
	case usecase.ErrEmailAlreadyTaken:
		SendErrorResponse(c, http.StatusConflict, "メールアドレスは既に使用されています")
	case usecase.ErrVersionConflict:
		SendErrorResponse(c, http.StatusConflict, "更新が競合しました。再度お試しください")
	case usecase.ErrInvalidOldPassword:
		SendErrorResponse(c, http.StatusBadRequest, "現在のパスワードが正しくありません")
	case usecase.ErrPasswordMismatch:
		SendErrorResponse(c, http.StatusBadRequest, "新しいパスワードが一致しません")
	case usecase.ErrWeakPassword:
		SendErrorResponse(c, http.StatusBadRequest, "パスワードが弱すぎます")
	default:
		// エラーメッセージから判断
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			SendErrorResponse(c, http.StatusNotFound, fmt.Sprintf("%sが見つかりません", resourceName))
		} else if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "already taken") {
			SendErrorResponse(c, http.StatusConflict, err.Error())
		} else if strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "validation") {
			SendErrorResponse(c, http.StatusBadRequest, err.Error())
		} else if strings.Contains(errMsg, "unauthorized") {
			SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else if strings.Contains(errMsg, "forbidden") {
			SendErrorResponse(c, http.StatusForbidden, err.Error())
		} else {
			// ドメインエラーの場合は専用ハンドラーを使用
			if domainErr, ok := err.(*errors.DomainError); ok {
				HandleError(c, domainErr)
			} else {
				SendErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		}
	}
}

// HandleAuthError 認証関連エラーの処理
func HandleAuthError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	
	// レート制限エラー
	if strings.Contains(errMsg, "too many") {
		c.JSON(http.StatusTooManyRequests, StandardErrorResponse{
			Error: errMsg,
			Code:  "RATE_LIMIT_EXCEEDED",
		})
		return
	}
	
	// トークンの再利用検出
	if strings.Contains(errMsg, "reuse detected") {
		c.JSON(http.StatusUnauthorized, StandardErrorResponse{
			Error: errMsg,
			Code:  "TOKEN_REUSE_DETECTED",
		})
		return
	}
	
	// セッション期限切れ
	if strings.Contains(errMsg, "session expired") {
		c.JSON(http.StatusUnauthorized, StandardErrorResponse{
			Error: errMsg,
			Code:  "SESSION_EXPIRED",
		})
		return
	}
	
	// その他の認証エラー
	c.JSON(http.StatusUnauthorized, StandardErrorResponse{
		Error: errMsg,
		Code:  "INVALID_CREDENTIALS",
	})
}

// ValidateUserPermission ユーザー権限の検証
func ValidateUserPermission(c *gin.Context, targetUserID string) bool {
	currentUser, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, http.StatusUnauthorized, "認証されていません")
		return false
	}

	// ユーザー情報の型アサーション（実際のユーザー構造体に合わせて調整）
	user, ok := currentUser.(map[string]interface{})
	if !ok {
		SendErrorResponse(c, http.StatusInternalServerError, "ユーザー情報の取得に失敗しました")
		return false
	}

	userID, hasID := user["id"].(string)
	isAdmin, hasAdmin := user["is_admin"].(bool)

	if !hasID {
		SendErrorResponse(c, http.StatusInternalServerError, "ユーザーIDが見つかりません")
		return false
	}

	// 自分自身または管理者の場合はアクセス許可
	if userID == targetUserID || (hasAdmin && isAdmin) {
		return true
	}

	SendErrorResponse(c, http.StatusForbidden, "アクセス権限がありません")
	return false
}

// ParseOptionalStringParam オプショナルな文字列パラメータの解析
func ParseOptionalStringParam(c *gin.Context, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}