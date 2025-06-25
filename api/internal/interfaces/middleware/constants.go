package middleware

import "net/http"

// HTTP ステータスコード定数
const (
	StatusUnauthorized = http.StatusUnauthorized
	StatusForbidden    = http.StatusForbidden
)

// エラーメッセージ定数
const (
	ErrMsgAuthRequired      = "認証が必要です"
	ErrMsgTokenRequired     = "Authorization token required"
	ErrMsgInvalidToken      = "Invalid or expired token"
	ErrMsgTokenRevoked      = "Token has been revoked"
	ErrMsgTokenExpired      = "Token has expired"
	ErrMsgUserNotFound      = "ユーザー情報が見つかりません"
	ErrMsgInvalidUserID     = "ユーザーID形式が不正です"
	ErrMsgAdminRequired     = "管理者権限が必要です"
	ErrMsgInsufficientPerms = "Insufficient permissions to access this resource"
)

// エラーコード定数
const (
	ErrCodeUnauthorized        = "unauthorized"
	ErrCodeInvalidToken        = "invalid_token"
	ErrCodeTokenRevoked        = "token_revoked"
	ErrCodeTokenExpired        = "token_expired"
	ErrCodeInsufficientPerms   = "insufficient_permissions"
)