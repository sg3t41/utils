package repository

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
)

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	// Create creates a new article
	Create(ctx context.Context, article *entity.Article) error

	// FindByID finds an article by its ID
	FindByID(ctx context.Context, id string) (*entity.Article, error)

	// FindAll finds articles with optional filtering and pagination
	FindAll(ctx context.Context, filter ArticleFilter) ([]*entity.Article, int, error)

	// Update updates an existing article
	Update(ctx context.Context, article *entity.Article) error

	// Delete deletes an article by its ID
	Delete(ctx context.Context, id string) error

	// FindByStatus finds articles by status with pagination
	FindByStatus(ctx context.Context, status entity.ArticleStatus, limit, offset int) ([]*entity.Article, int, error)

	// FindByTag finds articles that have a specific tag
	FindByTag(ctx context.Context, tag string, limit, offset int) ([]*entity.Article, int, error)

	// FindByAuthor finds articles by author ID
	FindByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Article, int, error)
}

// ArticleFilter represents filtering options for article queries
type ArticleFilter struct {
	Status     *entity.ArticleStatus `json:"status"`
	AuthorID   *string               `json:"author_id"`
	Tag        *string               `json:"tag"`
	Search     *string               `json:"search"`     // Search in title and content
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	Sort       string                `json:"sort"`       // created_at, updated_at, published_at, title
	Order      string                `json:"order"`      // asc, desc
	DateFrom   *string               `json:"date_from"`  // Filter by created_at >= date
	DateTo     *string               `json:"date_to"`    // Filter by created_at <= date
}

// GetOffset calculates the offset for pagination
func (f *ArticleFilter) GetOffset() int {
	if f.Page <= 0 {
		f.Page = 1
	}
	return (f.Page - 1) * f.Limit
}

// GetLimit returns the limit with a default value
func (f *ArticleFilter) GetLimit() int {
	if f.Limit <= 0 {
		return 10 // Default limit
	}
	if f.Limit > 100 {
		return 100 // Maximum limit
	}
	return f.Limit
}

// GetSort returns the sort field with a default value
func (f *ArticleFilter) GetSort() string {
	switch f.Sort {
	case "created_at", "updated_at", "published_at", "title":
		return f.Sort
	default:
		return "created_at"
	}
}

// GetOrder returns the sort order with a default value
func (f *ArticleFilter) GetOrder() string {
	switch f.Order {
	case "asc", "desc":
		return f.Order
	default:
		return "desc"
	}
}