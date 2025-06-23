package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/interfaces/dto"
)

type ArticleHandler struct {
	createArticleUseCase    *usecase.CreateArticleUseCase
	getArticlesUseCase      *usecase.GetArticlesUseCase
	getArticleUseCase       *usecase.GetArticleUseCase
	updateArticleUseCase    *usecase.UpdateArticleUseCase
	deleteArticleUseCase    *usecase.DeleteArticleUseCase
	publishArticleUseCase   *usecase.PublishArticleUseCase
	unpublishArticleUseCase *usecase.UnpublishArticleUseCase
}

func NewArticleHandler(
	createArticleUseCase *usecase.CreateArticleUseCase,
	getArticlesUseCase *usecase.GetArticlesUseCase,
	getArticleUseCase *usecase.GetArticleUseCase,
	updateArticleUseCase *usecase.UpdateArticleUseCase,
	deleteArticleUseCase *usecase.DeleteArticleUseCase,
	publishArticleUseCase *usecase.PublishArticleUseCase,
	unpublishArticleUseCase *usecase.UnpublishArticleUseCase,
) *ArticleHandler {
	return &ArticleHandler{
		createArticleUseCase:    createArticleUseCase,
		getArticlesUseCase:      getArticlesUseCase,
		getArticleUseCase:       getArticleUseCase,
		updateArticleUseCase:    updateArticleUseCase,
		deleteArticleUseCase:    deleteArticleUseCase,
		publishArticleUseCase:   publishArticleUseCase,
		unpublishArticleUseCase: unpublishArticleUseCase,
	}
}

// CreateArticle creates a new article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	req, exists := c.Get("validated_body")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "バリデーションが実行されていません"})
		return
	}

	createArticleReq, ok := req.(*dto.CreateArticleRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}


	// Convert empty strings to nil for optional fields
	var articleImage *string
	if createArticleReq.ArticleImage != nil && *createArticleReq.ArticleImage != "" {
		articleImage = createArticleReq.ArticleImage
	}

	input := usecase.CreateArticleInput{
		Title:    createArticleReq.Title,
		Content:  createArticleReq.Content,
		Summary:  createArticleReq.Summary,
		Tags:     createArticleReq.Tags,
		ArticleImage:    articleImage,
	}

	output, err := h.createArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusCreated, response)
}

// GetArticles retrieves a list of articles with filtering and pagination
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	params := parseArticleQueryParams(c)

	// Parse date filters
	var dateFrom, dateTo *time.Time
	if params.DateFrom != "" {
		if from, err := time.Parse("2006-01-02", params.DateFrom); err == nil {
			dateFrom = &from
		}
	}
	if params.DateTo != "" {
		if to, err := time.Parse("2006-01-02", params.DateTo); err == nil {
			dateTo = &to
		}
	}

	input := usecase.GetArticlesInput{
		Page:     params.Page,
		Limit:    params.Limit,
		Sort:     params.Sort,
		Order:    params.Order,
		Status:   params.Status,
		Tag:      params.Tag,
		Search:   params.Search,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	output, err := h.getArticlesUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := buildArticlesPaginatedResponse(output.Articles, output.Total, params)
	c.JSON(http.StatusOK, response)
}

// GetArticle retrieves a single article by ID
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := parseArticleIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "記事IDが無効です"})
		return
	}

	input := usecase.GetArticleInput{
		ID: id,
	}

	output, err := h.getArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// UpdateArticle updates an existing article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := parseArticleIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "記事IDが無効です"})
		return
	}

	req, exists := c.Get("validated_body")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "バリデーションが実行されていません"})
		return
	}

	updateArticleReq, ok := req.(*dto.UpdateArticleRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	// Convert empty strings to nil for optional fields  
	var articleImage *string
	if updateArticleReq.ArticleImage != nil && *updateArticleReq.ArticleImage != "" {
		articleImage = updateArticleReq.ArticleImage
	}

	input := usecase.UpdateArticleInput{
		ID:      id,
		Title:   updateArticleReq.Title,
		Content: updateArticleReq.Content,
		Summary: updateArticleReq.Summary,
		Tags:    updateArticleReq.Tags,
		ArticleImage:   articleImage,
	}

	output, err := h.updateArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// DeleteArticle deletes an article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := parseArticleIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "記事IDが無効です"})
		return
	}

	input := usecase.DeleteArticleInput{
		ID: id,
	}

	_, err = h.deleteArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// PublishArticle publishes an article
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	id, err := parseArticleIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "記事IDが無効です"})
		return
	}

	input := usecase.PublishArticleInput{
		ID: id,
	}

	output, err := h.publishArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// UnpublishArticle unpublishes an article
func (h *ArticleHandler) UnpublishArticle(c *gin.Context) {
	id, err := parseArticleIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "記事IDが無効です"})
		return
	}

	input := usecase.UnpublishArticleInput{
		ID: id,
	}

	output, err := h.unpublishArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "記事が見つかりません"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// Helper functions

func convertArticleToResponse(article *entity.Article) *dto.ArticleResponse {
	var publishedAt *string
	if article.PublishedAt != nil {
		published := article.PublishedAt.Format("2006-01-02T15:04:05Z")
		publishedAt = &published
	}

	return &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		Summary:     article.Summary,
		Status:      string(article.Status),
		Tags:        article.Tags,
		ArticleImage:       article.ArticleImage,
		CreatedAt:   article.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   article.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		PublishedAt: publishedAt,
	}
}

func convertArticleToListResponse(article *entity.Article) dto.ArticleListResponse {
	var publishedAt *string
	if article.PublishedAt != nil {
		published := article.PublishedAt.Format("2006-01-02T15:04:05Z")
		publishedAt = &published
	}

	return dto.ArticleListResponse{
		ID:          article.ID,
		Title:       article.Title,
		Summary:     article.Summary,
		Status:      string(article.Status),
		Tags:        article.Tags,
		ArticleImage:       article.ArticleImage,
		CreatedAt:   article.CreatedAt.Format("2006-01-02T15:04:05Z"),
		PublishedAt: publishedAt,
	}
}

func buildArticlesPaginatedResponse(articles []*entity.Article, total int, params dto.ListArticlesQuery) *dto.ArticlesResponse {
	data := make([]dto.ArticleListResponse, len(articles))
	for i, article := range articles {
		data[i] = convertArticleToListResponse(article)
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	totalPages := (total + params.Limit - 1) / params.Limit
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	return &dto.ArticlesResponse{
		Data: data,
		Pagination: dto.PaginationMeta{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Meta: dto.SortMeta{
			Sort:  params.Sort,
			Order: params.Order,
		},
	}
}

func parseArticleIDParam(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("missing article ID")
	}
	return id, nil
}

func parseArticleQueryParams(c *gin.Context) dto.ListArticlesQuery {
	params := dto.ListArticlesQuery{}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}

	params.Sort = c.Query("sort")
	params.Order = c.Query("order")
	params.Status = c.Query("status")
	params.Tag = c.Query("tag")
	params.Search = c.Query("search")
	params.DateFrom = c.Query("date_from")
	params.DateTo = c.Query("date_to")

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Sort == "" {
		params.Sort = "created_at"
	}
	if params.Order == "" {
		params.Order = "desc"
	}

	return params
}