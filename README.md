# sg3t41 個人用アプリケーション

Go (Gin)、Next.js、PostgreSQLで構築された個人用ダッシュボードアプリケーションです。

## 概要

個人の情報管理とコンテンツ作成を効率化するためのWebアプリケーションです。ブログ、制作物管理、収支表、ランキング、ファイル共有、写真管理、メモ、リンク集など、様々な機能を統合したダッシュボードを提供します。

## 技術構成

- **バックエンド**: Go + Gin フレームワーク（Clean Architecture/DDD）
- **フロントエンド**: Next.js + TypeScript + Tailwind CSS
- **データベース**: PostgreSQL
- **認証**: LINE Login対応
- **コンテナ化**: Docker & Docker Compose

## セットアップ

### 必要な環境

- Docker と Docker Compose
- Node.js（ローカル開発時）
- Go 1.23+（ローカル開発時）

### Docker Composeでの起動

1. リポジトリをクローン
2. 環境設定ファイルをコピー:
   ```bash
   cp .env.example .env
   cp frontend/.env.local.example frontend/.env.local
   ```
3. 全サービスを起動:
   ```bash
   docker-compose up --build
   ```

アプリケーションは以下のURLでアクセス可能です:
- フロントエンド: http://localhost:3000
- API: http://localhost:8080
- データベース: localhost:5432

### ローカル開発

#### バックエンド（API）
```bash
cd api
go mod download
go run cmd/server/main.go
```

#### フロントエンド
```bash
cd frontend
npm install
npm run dev
```

## プロジェクト構成

```
utils/
├── api/                    # Goバックエンド
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── frontend/               # Next.jsフロントエンド
│   ├── app/
│   ├── components/
│   └── hooks/
├── db/                     # データベース初期化
│   └── init/
├── docs/                   # ドキュメント
│   ├── dev_diary/          # 開発日誌
│   └── dev_ticket/         # チケット管理
└── docker-compose.yml
```

## 主な機能

### 実装済み機能
- **ブログシステム**: 記事の作成・編集・公開管理
- **ユーザー管理**: LINE認証による認証システム
- **ファイルアップロード**: 画像アップロード機能
- **レスポンシブUI**: モバイル・タブレット・デスクトップ対応

### 計画中の機能
- **制作物管理**: プロジェクト・ポートフォリオ管理
- **収支表**: 家計簿・収支管理機能
- **ランキング**: 各種ランキング管理
- **共有**: ファイル共有システム
- **写真**: フォトギャラリー
- **メモ**: 簡易メモ機能
- **リンク**: ブックマーク・リンク集

## API エンドポイント

### 認証
- `POST /api/v1/auth/login` - ログイン
- `GET /api/v1/auth/line/url` - LINE認証URL取得
- `POST /api/v1/auth/line/callback` - LINE認証コールバック

### 記事管理
- `GET /api/v1/articles` - 記事一覧取得
- `GET /api/v1/articles/:id` - 記事詳細取得
- `POST /api/v1/articles` - 記事作成（管理者のみ）
- `PUT /api/v1/articles/:id` - 記事更新（管理者のみ）
- `DELETE /api/v1/articles/:id` - 記事削除（管理者のみ）

### ユーザー管理
- `GET /api/v1/users` - ユーザー一覧取得
- `GET /api/v1/users/:id` - ユーザー詳細取得
- `POST /api/v1/users` - ユーザー作成

## 開発メモ

詳細な開発ログと進捗管理は `docs/dev_diary/` と `docs/dev_ticket/` をご覧ください。

## 特徴

- **個人使用特化**: 単一ユーザー（st）による管理を前提とした設計
- **モダンUI**: Windows 10風タイルデザインによる直感的なダッシュボード
- **Clean Architecture**: 保守性・拡張性を重視したアーキテクチャ
- **完全レスポンシブ**: あらゆるデバイスサイズに対応