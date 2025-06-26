package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/interfaces/dto"
	"go.uber.org/zap"
)

// LinkHandler リンク管理ハンドラー
type LinkHandler struct {
	linkRepo repository.LinkRepository
	logger   *zap.Logger
}

// NewLinkHandler リンクハンドラーの新しいインスタンスを作成
func NewLinkHandler(linkRepo repository.LinkRepository, logger *zap.Logger) *LinkHandler {
	return &LinkHandler{
		linkRepo: linkRepo,
		logger:   logger,
	}
}

// GetLinks リンク一覧を取得
func (h *LinkHandler) GetLinks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	// activeパラメータをチェック
	activeOnly := c.Query("active") == "true"
	
	var links []*entity.Link
	var err error
	
	if activeOnly {
		links, err = h.linkRepo.GetActiveLinks(c.Request.Context(), userID)
	} else {
		links, err = h.linkRepo.GetAll(c.Request.Context(), userID)
	}
	
	if err != nil {
		h.logger.Error("リンク一覧の取得に失敗", zap.Error(err))
		respondInternalError(c, "リンク一覧の取得に失敗しました")
		return
	}
	
	// レスポンスに変換
	linkResponses := make([]*dto.LinkResponse, len(links))
	for i, link := range links {
		linkResponses[i] = h.entityToResponse(link)
	}
	
	response := &dto.LinkListResponse{
		Links: linkResponses,
		Total: len(linkResponses),
	}
	
	c.JSON(http.StatusOK, response)
}

// GetLink 特定のリンクを取得
func (h *LinkHandler) GetLink(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "無効なリンクIDです")
		return
	}
	
	link, err := h.linkRepo.GetByID(c.Request.Context(), linkID)
	if err != nil {
		h.logger.Error("リンクの取得に失敗", zap.Error(err))
		respondInternalError(c, "リンクの取得に失敗しました")
		return
	}
	
	c.JSON(http.StatusOK, h.entityToResponse(link))
}

// CreateLink リンクを作成
func (h *LinkHandler) CreateLink(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	var req dto.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストデータが無効です: "+err.Error())
		return
	}
	
	// エンティティに変換
	link := &entity.Link{
		Title:       req.Title,
		URL:         req.URL,
		Description: req.Description,
		Platform:    req.Platform,
		IconName:    req.IconName,
		UserID:      userID,
	}
	
	// デフォルト値を設定
	h.setDefaultValues(link, &req)
	
	// リンクを作成
	if err := h.linkRepo.Create(c.Request.Context(), link); err != nil {
		h.logger.Error("リンクの作成に失敗", zap.Error(err))
		respondInternalError(c, "リンクの作成に失敗しました")
		return
	}
	
	h.logger.Info("リンクを作成しました", 
		zap.Int("link_id", link.ID),
		zap.String("title", link.Title),
		zap.String("user_id", userID.String()))
	
	c.JSON(http.StatusCreated, h.entityToResponse(link))
}

// UpdateLink リンクを更新
func (h *LinkHandler) UpdateLink(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "無効なリンクIDです")
		return
	}
	
	userID := getUserIDFromContext(c)
	
	var req dto.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストデータが無効です: "+err.Error())
		return
	}
	
	// 既存のリンクを取得
	link, err := h.linkRepo.GetByID(c.Request.Context(), linkID)
	if err != nil {
		h.logger.Error("リンクの取得に失敗", zap.Error(err))
		respondInternalError(c, "リンクの取得に失敗しました")
		return
	}
	
	// 権限チェック
	if link.UserID != userID {
		respondBadRequest(c, "このリンクを更新する権限がありません")
		return
	}
	
	// フィールドを更新
	h.updateLinkFields(link, &req)
	
	// リンクを更新
	if err := h.linkRepo.Update(c.Request.Context(), link); err != nil {
		h.logger.Error("リンクの更新に失敗", zap.Error(err))
		respondInternalError(c, "リンクの更新に失敗しました")
		return
	}
	
	h.logger.Info("リンクを更新しました", 
		zap.Int("link_id", link.ID),
		zap.String("user_id", userID.String()))
	
	c.JSON(http.StatusOK, h.entityToResponse(link))
}

// DeleteLink リンクを削除
func (h *LinkHandler) DeleteLink(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "無効なリンクIDです")
		return
	}
	
	userID := getUserIDFromContext(c)
	
	// 既存のリンクを取得して権限チェック
	link, err := h.linkRepo.GetByID(c.Request.Context(), linkID)
	if err != nil {
		h.logger.Error("リンクの取得に失敗", zap.Error(err))
		respondInternalError(c, "リンクの取得に失敗しました")
		return
	}
	
	if link.UserID != userID {
		respondBadRequest(c, "このリンクを削除する権限がありません")
		return
	}
	
	// リンクを削除
	if err := h.linkRepo.Delete(c.Request.Context(), linkID); err != nil {
		h.logger.Error("リンクの削除に失敗", zap.Error(err))
		respondInternalError(c, "リンクの削除に失敗しました")
		return
	}
	
	h.logger.Info("リンクを削除しました", 
		zap.Int("link_id", linkID),
		zap.String("user_id", userID.String()))
	
	respondWithSuccess(c, "リンクを削除しました")
}

// UpdateLinkOrder リンクの表示順序を更新
func (h *LinkHandler) UpdateLinkOrder(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "無効なリンクIDです")
		return
	}
	
	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストデータが無効です: "+err.Error())
		return
	}
	
	if err := h.linkRepo.UpdateOrder(c.Request.Context(), linkID, req.OrderIndex); err != nil {
		h.logger.Error("表示順序の更新に失敗", zap.Error(err))
		respondInternalError(c, "表示順序の更新に失敗しました")
		return
	}
	
	h.logger.Info("表示順序を更新しました", 
		zap.Int("link_id", linkID),
		zap.Int("order_index", req.OrderIndex))
	
	respondWithSuccess(c, "表示順序を更新しました")
}

// Helper methods

func (h *LinkHandler) entityToResponse(link *entity.Link) *dto.LinkResponse {
	return &dto.LinkResponse{
		ID:              link.ID,
		Title:           link.Title,
		URL:             link.URL,
		Description:     link.Description,
		Platform:        link.Platform,
		IconName:        link.IconName,
		BackgroundColor: link.BackgroundColor,
		TextColor:       link.TextColor,
		OrderIndex:      link.OrderIndex,
		IsActive:        link.IsActive,
		UserID:          link.UserID.String(),
		CreatedAt:       link.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       link.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *LinkHandler) setDefaultValues(link *entity.Link, req *dto.CreateLinkRequest) {
	// プラットフォームのデフォルト色を設定
	platform := entity.LinkPlatform(req.Platform)
	defaultBg, defaultText := platform.GetDefaultColors()
	
	if req.BackgroundColor != nil {
		link.BackgroundColor = *req.BackgroundColor
	} else {
		link.BackgroundColor = defaultBg
	}
	
	if req.TextColor != nil {
		link.TextColor = *req.TextColor
	} else {
		link.TextColor = defaultText
	}
	
	if req.OrderIndex != nil {
		link.OrderIndex = *req.OrderIndex
	} else {
		link.OrderIndex = 0
	}
	
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	} else {
		link.IsActive = true
	}
}

func (h *LinkHandler) updateLinkFields(link *entity.Link, req *dto.UpdateLinkRequest) {
	if req.Title != nil {
		link.Title = *req.Title
	}
	if req.URL != nil {
		link.URL = *req.URL
	}
	if req.Description != nil {
		link.Description = req.Description
	}
	if req.Platform != nil {
		link.Platform = *req.Platform
	}
	if req.IconName != nil {
		link.IconName = req.IconName
	}
	if req.BackgroundColor != nil {
		link.BackgroundColor = *req.BackgroundColor
	}
	if req.TextColor != nil {
		link.TextColor = *req.TextColor
	}
	if req.OrderIndex != nil {
		link.OrderIndex = *req.OrderIndex
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}
}

// getUserIDFromContext コンテキストからユーザーIDを取得（仮実装）
func getUserIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: 認証機能実装後に適切に実装
	userID, _ := uuid.Parse("6222a0b4-90c7-4cfa-ab26-4131155ff544")
	return userID // 暫定的に固定値
}