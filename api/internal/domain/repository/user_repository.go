package repository

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
)

type PaginationFilter struct {
	Search      string
	Status      string
	CreatedFrom string
	CreatedTo   string
}

type SortOption struct {
	Field string
	Order string
}

type PaginationResult struct {
	Users []*entity.User
	Total int
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByLineUserID(ctx context.Context, lineUserID string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	FindWithOffsetPagination(ctx context.Context, limit, offset int, filter PaginationFilter, sort SortOption) (*PaginationResult, error)
	FindWithCursorPagination(ctx context.Context, limit int, cursor string, filter PaginationFilter, sort SortOption) ([]*entity.User, error)
}