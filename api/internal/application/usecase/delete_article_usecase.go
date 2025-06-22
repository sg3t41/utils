package usecase

import (
	"context"
	"fmt"

	"github.com/sg3t41/api/internal/domain/repository"
)

type DeleteArticleInput struct {
	ID string
}

type DeleteArticleOutput struct {
	Success bool
}

type DeleteArticleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewDeleteArticleUseCase(articleRepo repository.ArticleRepository) *DeleteArticleUseCase {
	return &DeleteArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *DeleteArticleUseCase) Execute(ctx context.Context, input DeleteArticleInput) (*DeleteArticleOutput, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("article ID is required")
	}

	// Check if article exists
	existingArticle, err := uc.articleRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if existingArticle == nil {
		return nil, fmt.Errorf("article not found")
	}

	// Delete the article
	err = uc.articleRepo.Delete(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete article: %w", err)
	}

	return &DeleteArticleOutput{
		Success: true,
	}, nil
}