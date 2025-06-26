#!/bin/bash

# LINE Bot Webhook テスト用スクリプト

echo "LINE Bot Webhook テスト開始..."

# テスト用のJSONペイロード（署名は無効だがローカルテスト用）
TEST_PAYLOAD='{
  "events": [
    {
      "type": "message",
      "replyToken": "test-reply-token",
      "source": {
        "userId": "test-user-id",
        "type": "user"
      },
      "message": {
        "id": "test-message-id",
        "type": "text",
        "text": "ハロー"
      },
      "timestamp": 1640995200000,
      "mode": "active"
    }
  ]
}'

echo "テストペイロード："
echo "$TEST_PAYLOAD" | jq .

echo ""
echo "API サーバーが動作中か確認..."
curl -s http://localhost:8080/health | jq .

echo ""
echo "LINE Bot Webhook エンドポイントをテスト..."
# 署名を無効にするため、Channel Secretを空にしてテスト
curl -X POST http://localhost:8080/api/v1/linebot/webhook \
  -H "Content-Type: application/json" \
  -H "X-Line-Signature: sha256=" \
  -d "$TEST_PAYLOAD" \
  -v