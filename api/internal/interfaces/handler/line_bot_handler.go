package handler

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/service"
	"github.com/sg3t41/api/internal/interfaces/dto"
	"go.uber.org/zap"
)

// LineBotHandler LINE Botハンドラのインターフェース
type LineBotHandler interface {
	Webhook(c *gin.Context)
}

// lineBotHandler LINE Botハンドラの実装
type lineBotHandler struct {
	service service.LineBotService
	logger  *zap.Logger
}

// NewLineBotHandler LINE Botハンドラの生成
func NewLineBotHandler(service service.LineBotService, logger *zap.Logger) LineBotHandler {
	return &lineBotHandler{
		service: service,
		logger:  logger,
	}
}

// Webhook LINE BotのWebhookエンドポイント
func (h *lineBotHandler) Webhook(c *gin.Context) {
	body, err := h.readRequestBody(c)
	if err != nil {
		h.logger.Error("リクエストボディの読み取りに失敗", zap.Error(err))
		respondBadRequest(c, "リクエストボディの読み取りに失敗")
		return
	}

	if !h.validateSignature(c, body) {
		return
	}

	req, err := h.parseRequest(body)
	if err != nil {
		h.logger.Error("JSONパースに失敗", zap.Error(err))
		respondBadRequest(c, "JSONパースに失敗")
		return
	}

	if err := h.service.HandleWebhookEvents(req.Events); err != nil {
		h.logger.Error("イベント処理に失敗", zap.Error(err))
		respondInternalError(c, "イベント処理に失敗")
		return
	}

	respondWithSuccess(c, "OK")
}

// readRequestBody リクエストボディを読み取り
func (h *lineBotHandler) readRequestBody(c *gin.Context) ([]byte, error) {
	return io.ReadAll(c.Request.Body)
}

// validateSignature 署名を検証
func (h *lineBotHandler) validateSignature(c *gin.Context, body []byte) bool {
	signature := c.GetHeader(service.SignatureHeader)
	if signature == "" {
		h.logger.Error("署名ヘッダーが空")
		respondBadRequest(c, "署名ヘッダーが空")
		return false
	}

	if !h.service.ValidateSignature(signature, body) {
		h.logger.Error("署名検証に失敗")
		respondBadRequest(c, "署名検証に失敗")
		return false
	}

	return true
}

// parseRequest JSONリクエストをパース
func (h *lineBotHandler) parseRequest(body []byte) (*dto.WebhookRequest, error) {
	var req dto.WebhookRequest
	err := json.Unmarshal(body, &req)
	return &req, err
}