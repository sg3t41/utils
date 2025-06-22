package usecase

import (
	"context"
	"fmt"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type GetArticleInput struct {
	ID string
}

type GetArticleOutput struct {
	Article *entity.Article
}

type GetArticleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewGetArticleUseCase(articleRepo repository.ArticleRepository) *GetArticleUseCase {
	return &GetArticleUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *GetArticleUseCase) Execute(ctx context.Context, input GetArticleInput) (*GetArticleOutput, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("article ID is required")
	}

	article, err := uc.articleRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if article == nil {
		return nil, fmt.Errorf("article not found")
	}

	return &GetArticleOutput{
		Article: article,
	}, nil
}