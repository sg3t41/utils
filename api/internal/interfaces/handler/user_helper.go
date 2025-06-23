package handler

import (
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/interfaces/dto"
)

// User related helper functions

// UserResponse represents user response structure
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Version   int    `json:"version"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// convertUserToResponse converts User entity to response DTO
func convertUserToResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Version:   user.Version,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

// buildUsersPaginatedResponse builds paginated response for users
func buildUsersPaginatedResponse(users []*entity.User, total int, page, limit int) dto.UsersResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = convertUserToResponse(user)
	}

	totalPages := (total + limit - 1) / limit

	return dto.UsersResponse{
		Users: userResponses,
		Pagination: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}
}