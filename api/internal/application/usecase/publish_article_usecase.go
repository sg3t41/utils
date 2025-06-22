package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type PublishArticleInput struct {
	ID string
}

type PublishArticleOutput struct {
	Article *entity.Article
}

type PublishArticleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewPublishArticleUseCase(articleRepo repository.ArticleRepository) *PublishArticleUseCase {
	return &PublishArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *PublishArticleUseCase) Execute(ctx context.Context, input PublishArticleInput) (*PublishArticleOutput, error) {
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

	// Check if article is already published
	if existingArticle.Status == entity.ArticleStatusPublished {
		return nil, fmt.Errorf("article is already published")
	}

	// Update article status and published_at timestamp
	now := time.Now()
	existingArticle.Status = entity.ArticleStatusPublished
	existingArticle.PublishedAt = &now
	existingArticle.UpdatedAt = now

	// Save the updated article
	err = uc.articleRepo.Update(ctx, existingArticle)
	if err != nil {
		return nil, fmt.Errorf("failed to publish article: %w", err)
	}

	return &PublishArticleOutput{
		Article: existingArticle,
	}, nil
}