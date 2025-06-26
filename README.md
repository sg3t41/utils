# sg3t41 個人用アプリケーション

Go (Gin) + Next.js + PostgreSQL で構築された個人用ダッシュボードアプリケーション

## 起動方法

### Docker Compose（推奨）

```bash
# 1. 環境変数ファイルをコピー
cp .env.example .env

# 2. 必要な環境変数を設定（下記参照）
# 3. 起動
docker-compose up --build
```

**アクセス先:**
- フロントエンド: http://localhost:3000
- API: http://localhost:8080

### ローカル開発

```bash
# バックエンド
cd api && go run cmd/server/main.go

# フロントエンド  
cd frontend && npm run dev
```

## 必要な環境変数

`.env`ファイルに設定:

```bash
# データベース
POSTGRES_DB=utils_db
POSTGRES_USER=utils_user
POSTGRES_PASSWORD=utils_password
DB_HOST=localhost
DB_PORT=5432

# API
API_PORT=8080

# フロントエンド
NEXT_PUBLIC_API_URL=http://localhost:8080
NODE_ENV=development

# LINE認証（必要に応じて）
LINE_CLIENT_ID=your_line_client_id
LINE_CLIENT_SECRET=your_line_client_secret
LINE_REDIRECT_URL=http://localhost:3000/auth/line/callback

# LINE Bot（必要に応じて）
LINE_BOT_CHANNEL_SECRET=your_channel_secret
LINE_BOT_ACCESS_TOKEN=your_access_token
```

## データベースマイグレーション

```bash
make migrate-up     # マイグレーション実行
make migrate-down   # ロールバック
make migrate-status # ステータス確認
```