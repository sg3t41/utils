package service

import (
	"strings"

	"github.com/sg3t41/api/internal/interfaces/dto"
)

// MessageHandler メッセージ処理のインターフェース
type MessageHandler interface {
	HandleTextMessage(text string, replyToken string) error
}

// textMessageHandler テキストメッセージ処理の実装
type textMessageHandler struct {
	apiClient LineAPIClient
}

// NewTextMessageHandler テキストメッセージハンドラーの生成
func NewTextMessageHandler(apiClient LineAPIClient) MessageHandler {
	return &textMessageHandler{
		apiClient: apiClient,
	}
}

// HandleTextMessage テキストメッセージを処理
func (h *textMessageHandler) HandleTextMessage(text string, replyToken string) error {
	replyText := h.generateReplyText(text)
	if replyText == "" {
		return nil
	}

	message := dto.MessageBody{
		Type: MessageTypeText,
		Text: replyText,
	}
	
	return h.apiClient.SendReply(replyToken, []dto.MessageBody{message})
}

// generateReplyText 返信テキストを生成
func (h *textMessageHandler) generateReplyText(text string) string {
	lowerText := strings.ToLower(text)
	
	// コマンドマッピング（将来的に拡張可能）
	commands := map[string]string{
		"ハロー": "ハロー！",
		"hello": "ハロー！",
	}
	
	for cmd, reply := range commands {
		if lowerText == strings.ToLower(cmd) {
			return reply
		}
	}
	
	return ""
}