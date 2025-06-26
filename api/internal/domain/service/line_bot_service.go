package service

import (
	"github.com/sg3t41/api/internal/interfaces/dto"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/zap"
)

// LineBotService LINE Botサービスのインターフェース
type LineBotService interface {
	ValidateSignature(signature string, body []byte) bool
	HandleWebhookEvents(events []dto.WebhookEvent) error
}

// lineBotService LINE Botサービスの実装
type lineBotService struct {
	logger             *zap.Logger
	signatureValidator SignatureValidator
	messageHandler     MessageHandler
}

// NewLineBotService LINE Botサービスの生成
func NewLineBotService(cfg *config.Config, logger *zap.Logger) LineBotService {
	apiClient := NewLineAPIClient(cfg.LineBotAccessToken)
	return &lineBotService{
		logger:             logger,
		signatureValidator: NewSignatureValidator(cfg.LineBotChannelSecret),
		messageHandler:     NewTextMessageHandler(apiClient),
	}
}

// ValidateSignature 署名を検証
func (s *lineBotService) ValidateSignature(signature string, body []byte) bool {
	return s.signatureValidator.Validate(signature, body)
}

// HandleWebhookEvents webhookイベントを処理
func (s *lineBotService) HandleWebhookEvents(events []dto.WebhookEvent) error {
	for _, event := range events {
		if err := s.processEvent(event); err != nil {
			s.logger.Error("イベント処理に失敗",
				zap.String("event_type", event.Type),
				zap.String("user_id", event.Source.UserID),
				zap.Error(err),
			)
			// エラーが発生しても他のイベントを処理する
		}
	}
	return nil
}

// processEvent 個別のイベントを処理
func (s *lineBotService) processEvent(event dto.WebhookEvent) error {
	if !s.isValidEvent(event) {
		return nil
	}

	switch event.Type {
	case EventTypeMessage:
		return s.processMessageEvent(event)
	default:
		s.logger.Debug("未対応のイベントタイプ", zap.String("type", event.Type))
	}
	return nil
}

// isValidEvent イベントが処理可能か確認
func (s *lineBotService) isValidEvent(event dto.WebhookEvent) bool {
	return event.ReplyToken != "" && event.Source.UserID != ""
}

// processMessageEvent メッセージイベントを処理
func (s *lineBotService) processMessageEvent(event dto.WebhookEvent) error {
	if event.Message == nil {
		return nil
	}

	switch event.Message.Type {
	case MessageTypeText:
		s.logger.Info("テキストメッセージを受信",
			zap.String("text", event.Message.Text),
			zap.String("user_id", event.Source.UserID),
		)
		return s.messageHandler.HandleTextMessage(event.Message.Text, event.ReplyToken)
	default:
		s.logger.Debug("未対応のメッセージタイプ", zap.String("type", event.Message.Type))
	}

	return nil
}