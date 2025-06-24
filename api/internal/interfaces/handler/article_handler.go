package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
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
	createArticleReq, ok := GetValidatedBody[dto.CreateArticleRequest](c)
	if !ok {
		return
	}


	// Convert empty strings to nil for optional fields
	articleImage := h.processOptionalString(createArticleReq.ArticleImage)

	input := usecase.CreateArticleInput{
		Title:    createArticleReq.Title,
		Content:  createArticleReq.Content,
		Summary:  createArticleReq.Summary,
		Tags:     createArticleReq.Tags,
		ArticleImage:    articleImage,
	}

	output, err := h.createArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "記事")
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
		HandleUseCaseError(c, err, "記事")
		return
	}

	response := buildArticlesPaginatedResponse(output.Articles, output.Total, params)
	c.JSON(http.StatusOK, response)
}

// GetArticle retrieves a single article by ID
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	input := usecase.GetArticleInput{
		ID: id,
	}

	output, err := h.getArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "記事")
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// UpdateArticle updates an existing article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	updateArticleReq, ok := GetValidatedBody[dto.UpdateArticleRequest](c)
	if !ok {
		return
	}

	// Convert empty strings to nil for optional fields  
	articleImage := h.processOptionalString(updateArticleReq.ArticleImage)

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
		HandleUseCaseError(c, err, "記事")
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// DeleteArticle deletes an article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	input := usecase.DeleteArticleInput{
		ID: id,
	}

	_, err := h.deleteArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "記事")
		return
	}

	c.Status(http.StatusNoContent)
}

// PublishArticle publishes an article
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	input := usecase.PublishArticleInput{
		ID: id,
	}

	output, err := h.publishArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "記事")
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}

// UnpublishArticle unpublishes an article
func (h *ArticleHandler) UnpublishArticle(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	input := usecase.UnpublishArticleInput{
		ID: id,
	}

	output, err := h.unpublishArticleUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "記事")
		return
	}

	response := convertArticleToResponse(output.Article)
	c.JSON(http.StatusOK, response)
}



// processOptionalString オプショナルな文字列の処理
func (h *ArticleHandler) processOptionalString(value *string) *string {
	if value != nil && *value != "" {
		return value
	}
	return nil
}
