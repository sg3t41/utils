#!/bin/bash

echo "Starting ngrok to expose local API server..."
echo "====================================="
echo "APIエンドポイント: http://localhost:8080"
echo ""
echo "ngrokが起動したら、表示される以下のようなHTTPS URLをコピーしてください："
echo "例: https://xxxx-xxx-xxx-xxx.ngrok-free.app"
echo ""
echo "LINE DevelopersのWebhook URLには以下のように設定してください："
echo "https://xxxx-xxx-xxx-xxx.ngrok-free.app/api/v1/linebot/webhook"
echo "====================================="
echo ""

# ngrokでポート8080を公開
./ngrok http 8080