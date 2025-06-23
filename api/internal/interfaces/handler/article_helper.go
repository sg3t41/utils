package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/interfaces/dto"
)

// Article related helper functions

// parseArticleIDParam parses and validates article ID parameter
func parseArticleIDParam(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("記事IDが指定されていません")
	}
	return id, nil
}

// ArticleQueryParams represents query parameters for article listing
type ArticleQueryParams struct {
	Page     int
	Limit    int
	Sort     string
	Order    string
	Status   string
	Tag      string
	Search   string
	DateFrom string
	DateTo   string
}

// parseArticleQueryParams parses and validates query parameters
func parseArticleQueryParams(c *gin.Context) ArticleQueryParams {
	params := ArticleQueryParams{
		Page:     1,
		Limit:    10,
		Sort:     "created_at",
		Order:    "desc",
		Status:   c.Query("status"),
		Tag:      c.Query("tag"),
		Search:   c.Query("search"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	if order := c.Query("order"); order != "" {
		params.Order = order
	}

	return params
}

// convertArticleToResponse converts Article entity to response DTO
func convertArticleToResponse(article *entity.Article) dto.ArticleResponse {
	response := dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Summary:   article.Summary,
		Status:    string(article.Status),
		Tags:      article.Tags,
		CreatedAt: article.CreatedAt.Format(time.RFC3339),
		UpdatedAt: article.UpdatedAt.Format(time.RFC3339),
	}

	if article.ArticleImage != nil {
		response.ArticleImage = article.ArticleImage
	}

	if article.PublishedAt != nil {
		publishedAtStr := article.PublishedAt.Format(time.RFC3339)
		response.PublishedAt = &publishedAtStr
	}

	return response
}

// convertArticleToListResponse converts Article entity to list response DTO
func convertArticleToListResponse(article *entity.Article) dto.ArticleListResponse {
	response := dto.ArticleListResponse{
		ID:        article.ID,
		Title:     article.Title,
		Summary:   article.Summary,
		Status:    string(article.Status),
		Tags:      article.Tags,
		CreatedAt: article.CreatedAt.Format(time.RFC3339),
	}

	if article.ArticleImage != nil {
		response.ArticleImage = article.ArticleImage
	}

	if article.PublishedAt != nil {
		publishedAtStr := article.PublishedAt.Format(time.RFC3339)
		response.PublishedAt = &publishedAtStr
	}

	return response
}

// buildArticlesPaginatedResponse builds paginated response for articles
func buildArticlesPaginatedResponse(articles []*entity.Article, total int, params ArticleQueryParams) *dto.ArticlesResponse {
	data := make([]dto.ArticleListResponse, len(articles))
	for i, article := range articles {
		data[i] = convertArticleToListResponse(article)
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &dto.ArticlesResponse{
		Data: data,
		Pagination: dto.PaginationMeta{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
		Meta: dto.SortMeta{
			Sort:  params.Sort,
			Order: params.Order,
		},
	}
}