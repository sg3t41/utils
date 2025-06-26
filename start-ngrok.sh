#!/bin/bash

# ngrokでローカルのAPIサーバーを外部公開
echo "Starting ngrok to expose local API server..."
echo "API endpoint: http://localhost:8080"
echo ""
echo "ngrokを起動します。表示されるHTTPS URLをLINE Developersに設定してください。"
echo "Webhook URL例: https://xxxx.ngrok.io/api/v1/linebot/webhook"
echo ""

# ngrokがインストールされているか確認
if ! command -v ngrok &> /dev/null; then
    echo "ngrokがインストールされていません。"
    echo "以下のコマンドでインストールしてください："
    echo ""
    echo "# macOS"
    echo "brew install ngrok"
    echo ""
    echo "# Linux"
    echo "curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null && echo \"deb https://ngrok-agent.s3.amazonaws.com buster main\" | sudo tee /etc/apt/sources.list.d/ngrok.list && sudo apt update && sudo apt install ngrok"
    echo ""
    echo "# または公式サイトからダウンロード"
    echo "https://ngrok.com/download"
    exit 1
fi

# ngrokでポート8080を公開
ngrok http 8080