package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type UpdateArticleInput struct {
	ID      string
	Title   *string
	Content *string
	Summary *string
	Tags    []string
	ArticleImage   *string
}

type UpdateArticleOutput struct {
	Article *entity.Article
}

type UpdateArticleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewUpdateArticleUseCase(articleRepo repository.ArticleRepository) *UpdateArticleUseCase {
	return &UpdateArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *UpdateArticleUseCase) Execute(ctx context.Context, input UpdateArticleInput) (*UpdateArticleOutput, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("article ID is required")
	}

	// Get existing article
	existingArticle, err := uc.articleRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if existingArticle == nil {
		return nil, fmt.Errorf("article not found")
	}

	// Update fields if provided
	if input.Title != nil {
		if *input.Title == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
		existingArticle.Title = *input.Title
	}

	if input.Content != nil {
		if *input.Content == "" {
			return nil, fmt.Errorf("content cannot be empty")
		}
		existingArticle.Content = *input.Content
	}

	if input.Summary != nil {
		existingArticle.Summary = *input.Summary
	}

	if input.Tags != nil {
		// Validate tags
		for _, tag := range input.Tags {
			if tag == "" {
				return nil, fmt.Errorf("tag cannot be empty")
			}
		}
		existingArticle.Tags = input.Tags
	}

	if input.ArticleImage != nil {
		existingArticle.ArticleImage = input.ArticleImage
	}

	// Update timestamp
	existingArticle.UpdatedAt = time.Now()

	// Save the updated article
	err = uc.articleRepo.Update(ctx, existingArticle)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return &UpdateArticleOutput{
		Article: existingArticle,
	}, nil
}