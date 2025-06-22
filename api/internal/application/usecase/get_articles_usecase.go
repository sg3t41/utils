package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type GetArticlesInput struct {
	Page     int
	Limit    int
	Sort     string
	Order    string
	Status   string
	Tag      string
	Search   string
	DateFrom *time.Time
	DateTo   *time.Time
}

type GetArticlesOutput struct {
	Articles []*entity.Article
	Total    int
}

type GetArticlesUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewGetArticlesUseCase(articleRepo repository.ArticleRepository) *GetArticlesUseCase {
	return &GetArticlesUseCase{
		articleRepo: articleRepo,
	}
}

func (uc *GetArticlesUseCase) Execute(ctx context.Context, input GetArticlesInput) (*GetArticlesOutput, error) {
	// Set defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Sort == "" {
		input.Sort = "created_at"
	}
	if input.Order == "" {
		input.Order = "desc"
	}

	// Validate sort field
	validSortFields := map[string]bool{
		"created_at":   true,
		"updated_at":   true,
		"published_at": true,
		"title":        true,
	}
	if !validSortFields[input.Sort] {
		return nil, fmt.Errorf("invalid sort field: %s", input.Sort)
	}

	// Validate order
	if input.Order != "asc" && input.Order != "desc" {
		return nil, fmt.Errorf("invalid order: %s", input.Order)
	}

	// Validate status
	if input.Status != "" {
		validStatuses := map[string]bool{
			"draft":     true,
			"published": true,
			"archived":  true,
		}
		if !validStatuses[input.Status] {
			return nil, fmt.Errorf("invalid status: %s", input.Status)
		}
	}

	// Build filter
	// Convert status string to ArticleStatus pointer
	var statusPtr *entity.ArticleStatus
	if input.Status != "" {
		status := entity.ArticleStatus(input.Status)
		statusPtr = &status
	}

	// Convert time.Time pointers to string pointers
	var dateFromStr, dateToStr *string
	if input.DateFrom != nil {
		dateStr := input.DateFrom.Format("2006-01-02")
		dateFromStr = &dateStr
	}
	if input.DateTo != nil {
		dateStr := input.DateTo.Format("2006-01-02")
		dateToStr = &dateStr
	}

	// Convert tag and search to pointers
	var tagPtr, searchPtr *string
	if input.Tag != "" {
		tagPtr = &input.Tag
	}
	if input.Search != "" {
		searchPtr = &input.Search
	}

	filter := repository.ArticleFilter{
		Page:     input.Page,
		Limit:    input.Limit,
		Sort:     input.Sort,
		Order:    input.Order,
		Status:   statusPtr,
		Tag:      tagPtr,
		Search:   searchPtr,
		DateFrom: dateFromStr,
		DateTo:   dateToStr,
	}

	// Get articles with total count
	articles, total, err := uc.articleRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count articles: %w", err)
	}

	return &GetArticlesOutput{
		Articles: articles,
		Total:    total,
	}, nil
}