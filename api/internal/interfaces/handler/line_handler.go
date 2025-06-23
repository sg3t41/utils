package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/service"
	"github.com/sg3t41/api/internal/interfaces/dto"
	"github.com/sg3t41/api/pkg/config"
)

// LineHandler LINEログイン機能のハンドラー
type LineHandler struct {
	lineService           *service.LineService
	authService           *service.AuthenticationService
	createLineUserUseCase *usecase.CreateLineUserUseCase
}

// NewLineHandler LineHandlerのコンストラクタ
func NewLineHandler(
	config *config.Config,
	authService *service.AuthenticationService,
	createLineUserUseCase *usecase.CreateLineUserUseCase,
) *LineHandler {
	lineConfig := dto.LineOAuthConfig{
		ClientID:     config.LineClientID,
		ClientSecret: config.LineClientSecret,
		RedirectURL:  config.LineRedirectURL,
		Scope:        []string{"profile", "openid"},
	}

	lineService := service.NewLineService(lineConfig)

	return &LineHandler{
		lineService:           lineService,
		authService:           authService,
		createLineUserUseCase: createLineUserUseCase,
	}
}

// GetAuthURL LINE認証URL取得
func (h *LineHandler) GetAuthURL(c *gin.Context) {
	state := generateRandomState()
	authURL := h.lineService.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// CallbackPost LINEコールバック処理
func (h *LineHandler) CallbackPost(c *gin.Context) {
	var req dto.LineLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}


	// アクセストークン取得
	accessToken, err := h.getAccessToken(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get access token",
			"details": err.Error(),
		})
		return
	}

	// プロフィール取得
	profile, err := h.getProfile(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get profile",
			"details": err.Error(),
		})
		return
	}

	// ユーザー作成または取得
	user, err := h.createOrGetUser(c.Request.Context(), profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create/get user",
			"details": err.Error(),
		})
		return
	}

	// JWTトークン生成
	tokens, err := h.generateTokens(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate tokens",
			"details": err.Error(),
		})
		return
	}

	// 成功レスポンス
	h.sendSuccessResponse(c, tokens, user)
}

// getAccessToken アクセストークン取得（分離して可読性向上）
func (h *LineHandler) getAccessToken(code string) (string, error) {
	return h.lineService.GetAccessToken(code)
}

// getProfile プロフィール取得（分離して可読性向上）
func (h *LineHandler) getProfile(accessToken string) (*dto.LineProfile, error) {
	return h.lineService.GetProfile(accessToken)
}

// createOrGetUser ユーザー作成または取得（分離して可読性向上）
func (h *LineHandler) createOrGetUser(ctx context.Context, profile *dto.LineProfile) (*entity.User, error) {
	return h.createLineUserUseCase.ExecuteFromProfile(ctx, profile)
}

// generateTokens JWTトークン生成（分離して可読性向上）
func (h *LineHandler) generateTokens(c *gin.Context, user *entity.User) (*entity.TokenPair, error) {
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	return h.authService.GenerateTokensForUser(c.Request.Context(), user, ipAddress, userAgent)
}

// sendSuccessResponse 成功レスポンス送信（分離して可読性向上）
func (h *LineHandler) sendSuccessResponse(c *gin.Context, tokens *entity.TokenPair, user *entity.User) {
	response := dto.LineLoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         dto.ToUserResponse(user),
	}
	c.JSON(http.StatusOK, response)
}

// generateRandomState ランダムなstate生成
func generateRandomState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}