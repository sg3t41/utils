package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sg3t41/api/internal/interfaces/dto"
)

// LineAPIClient LINE APIクライアントのインターフェース
type LineAPIClient interface {
	SendReply(replyToken string, messages []dto.MessageBody) error
}

// lineAPIClient LINE APIクライアントの実装
type lineAPIClient struct {
	accessToken string
	httpClient  *http.Client
}

// NewLineAPIClient LINE APIクライアントの生成
func NewLineAPIClient(accessToken string) LineAPIClient {
	return &lineAPIClient{
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendReply 返信メッセージを送信
func (c *lineAPIClient) SendReply(replyToken string, messages []dto.MessageBody) error {
	replyMessage := dto.ReplyMessage{
		ReplyToken: replyToken,
		Messages:   messages,
	}

	jsonData, err := json.Marshal(replyMessage)
	if err != nil {
		return fmt.Errorf("メッセージのJSON変換に失敗: %w", err)
	}

	url := LineAPIBaseURL + LineReplyEndpoint
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("リクエスト作成に失敗: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("APIリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE APIがエラーを返しました: status=%d", resp.StatusCode)
	}

	return nil
}

// setHeaders リクエストヘッダーを設定
func (c *lineAPIClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", ContentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
}