package dto

import (
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
)

// UserResponse はユーザー情報のレスポンス用DTO
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Version   int    `json:"version"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToUserResponse はUserエンティティをUserResponseに変換する
func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Version:   user.Version,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateUserRequest はユーザー作成リクエスト用DTO
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// UpdateUserRequest はユーザー更新リクエスト用DTO
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
}

// UpdatePasswordRequest はパスワード更新リクエスト用DTO
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// GetUserQuery はユーザー取得クエリ用DTO
type GetUserQuery struct {
	ID string `form:"id" validate:"required,uuid"`
}

// ListUsersQuery はユーザー一覧取得クエリ用DTO
type ListUsersQuery struct {
	Page   int    `form:"page" validate:"min=1"`
	Limit  int    `form:"limit" validate:"min=1,max=100"`
	Search string `form:"search" validate:"max=100"`
}