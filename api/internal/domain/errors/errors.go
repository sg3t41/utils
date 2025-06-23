package errors

import "errors"

// Domain errors
var (
	// Common errors
	ErrNotFound        = errors.New("リソースが見つかりません")
	ErrAlreadyExists   = errors.New("リソースが既に存在します")
	ErrInvalidInput    = errors.New("入力データが無効です")
	ErrInternalServer  = errors.New("内部サーバーエラーが発生しました")
	ErrUnauthorized    = errors.New("認証が必要です")
	ErrForbidden       = errors.New("アクセス権限がありません")

	// User related errors
	ErrUserNotFound    = errors.New("ユーザーが見つかりません")
	ErrUserExists      = errors.New("ユーザーが既に存在します")
	ErrInvalidPassword = errors.New("パスワードが無効です")
	ErrEmailExists     = errors.New("メールアドレスが既に登録されています")

	// Article related errors
	ErrArticleNotFound     = errors.New("記事が見つかりません")
	ErrArticleExists       = errors.New("記事が既に存在します")
	ErrArticleUnauthorized = errors.New("記事の操作権限がありません")

	// Upload related errors
	ErrInvalidFileType   = errors.New("サポートされていないファイル形式です")
	ErrFileTooLarge      = errors.New("ファイルサイズが大きすぎます")
	ErrUploadFailed      = errors.New("ファイルのアップロードに失敗しました")
	ErrFileNotFound      = errors.New("ファイルが見つかりません")

	// Authentication related errors
	ErrInvalidToken      = errors.New("無効なトークンです")
	ErrTokenExpired      = errors.New("トークンの有効期限が切れています")
	ErrInvalidCredentials = errors.New("認証情報が無効です")
)

// ErrorType represents the type of error for HTTP status mapping
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeAlreadyExists
	ErrorTypeInvalidInput
	ErrorTypeInternalServer
	ErrorTypeUnauthorized
	ErrorTypeForbidden
)

// DomainError represents a domain-specific error with additional context
type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Error constructors
func NewNotFoundError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewAlreadyExistsError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeAlreadyExists,
		Message: message,
		Cause:   cause,
	}
}

func NewInvalidInputError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeInvalidInput,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalServerError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeInternalServer,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Cause:   cause,
	}
}

func NewForbiddenError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Cause:   cause,
	}
}