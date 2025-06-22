package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type UnpublishArticleInput struct {
	ID string
}

type UnpublishArticleOutput struct {
	Article *entity.Article
}

type UnpublishArticleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewUnpublishArticleUseCase(articleRepo repository.ArticleRepository) *UnpublishArticleUseCase {
	return &UnpublishArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *UnpublishArticleUseCase) Execute(ctx context.Context, input UnpublishArticleInput) (*UnpublishArticleOutput, error) {
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

	// Check if article is not published
	if existingArticle.Status != entity.ArticleStatusPublished {
		return nil, fmt.Errorf("article is not published")
	}

	// Update article status to draft and clear published_at timestamp
	existingArticle.Status = entity.ArticleStatusDraft
	existingArticle.PublishedAt = nil
	existingArticle.UpdatedAt = time.Now()

	// Save the updated article
	updatedArticle, err := uc.articleRepo.Update(ctx, existingArticle)
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish article: %w", err)
	}

	return &UnpublishArticleOutput{
		Article: updatedArticle,
	}, nil
}