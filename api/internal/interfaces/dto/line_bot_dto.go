package dto

// Package dto LINE Bot Webhook関連のデータ転送オブジェクト

// WebhookRequest LINE Botからのwebhookリクエスト
type WebhookRequest struct {
	Events []WebhookEvent `json:"events"`
}

// WebhookEvent webhookイベント
type WebhookEvent struct {
	Type       string   `json:"type"`
	Timestamp  int64    `json:"timestamp"`
	Source     Source   `json:"source"`
	Message    *Message `json:"message,omitempty"`
	ReplyToken string   `json:"replyToken,omitempty"`
}

// Source メッセージ送信者情報
type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

// Message メッセージ内容
type Message struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ReplyMessage 返信メッセージ
type ReplyMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []MessageBody `json:"messages"`
}

// MessageBody メッセージ本体
type MessageBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
}