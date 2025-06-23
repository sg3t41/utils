package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type CreateArticleInput struct {
	Title    string   `json:"title" validate:"required,min=1,max=500"`
	Content  string   `json:"content" validate:"required,min=1"`
	Summary  string   `json:"summary" validate:"max=1000"`
	AuthorID string   `json:"author_id" validate:"required"`
	Tags     []string `json:"tags" validate:"dive,min=1,max=50"`
	ArticleImage    *string  `json:"article_image" validate:"omitempty,max=500"`
}

type CreateArticleOutput struct {
	Article *entity.Article `json:"article"`
}

type CreateArticleUseCase struct {
	articleRepository repository.ArticleRepository
}

func NewCreateArticleUseCase(articleRepository repository.ArticleRepository) *CreateArticleUseCase {
	return &CreateArticleUseCase{
		articleRepository: articleRepository,
	}
}

func (uc *CreateArticleUseCase) Execute(ctx context.Context, input CreateArticleInput) (*CreateArticleOutput, error) {
	// Generate new UUID for the article
	articleID := uuid.New().String()

	// Create new article entity
	article := &entity.Article{
		ID:        articleID,
		Title:     input.Title,
		Content:   input.Content,
		Summary:   input.Summary,
		Status:    entity.ArticleStatusDraft, // New articles start as draft
		AuthorID:  input.AuthorID,
		Tags:      input.Tags,
		ArticleImage:     input.ArticleImage,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save article to repository
	err := uc.articleRepository.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return &CreateArticleOutput{
		Article: article,
	}, nil
}