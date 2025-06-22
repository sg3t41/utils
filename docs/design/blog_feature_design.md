# ブログ機能設計書

## 概要

utilアプリ内のブログ機能として、記事の作成・編集・削除・一覧表示機能を提供する。

## データモデル設計

### Article エンティティ

```go
type Article struct {
    ID          string    `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Content     string    `json:"content" db:"content"`
    Summary     string    `json:"summary" db:"summary"`     // 要約（一覧表示用）
    Status      string    `json:"status" db:"status"`       // draft, published, archived
    AuthorID    string    `json:"author_id" db:"author_id"` // ユーザーID（認証用）
    Tags        []string  `json:"tags" db:"tags"`           // タグ（JSONB）
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    PublishedAt *time.Time `json:"published_at" db:"published_at"` // 公開日時
}
```

### データベーステーブル

```sql
CREATE TABLE articles (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    summary VARCHAR(1000),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    author_id VARCHAR(36) REFERENCES users(id),
    tags JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_published_at ON articles(published_at);
```

## API設計

### エンドポイント一覧

```
GET    /api/v1/articles                 # 記事一覧取得
POST   /api/v1/articles                 # 記事作成
GET    /api/v1/articles/:id             # 記事詳細取得
PUT    /api/v1/articles/:id             # 記事更新
DELETE /api/v1/articles/:id             # 記事削除
POST   /api/v1/articles/:id/publish     # 記事公開
POST   /api/v1/articles/:id/unpublish   # 記事非公開
```

### リクエスト・レスポンス例

#### 記事一覧取得
```
GET /api/v1/articles?page=1&limit=10&status=published&tag=tech

Response:
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Go言語の基礎",
      "summary": "Go言語の基本的な使い方について",
      "status": "published",
      "tags": ["go", "programming"],
      "created_at": "2025-06-22T10:00:00Z",
      "published_at": "2025-06-22T11:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

#### 記事作成
```
POST /api/v1/articles

Request:
{
  "title": "新しい記事",
  "content": "記事の内容...",
  "summary": "記事の要約",
  "tags": ["tech", "blog"]
}

Response:
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "新しい記事",
  "content": "記事の内容...",
  "summary": "記事の要約",
  "status": "draft",
  "tags": ["tech", "blog"],
  "created_at": "2025-06-22T10:00:00Z"
}
```

## フロントエンド設計

### ページ構成

```
/                          # ホーム（記事一覧）
/articles                  # 記事一覧（管理画面）
/articles/new              # 新規記事作成
/articles/:id              # 記事詳細
/articles/:id/edit         # 記事編集
```

### コンポーネント設計

```
components/
├── article/
│   ├── ArticleCard.tsx      # 記事カード
│   ├── ArticleList.tsx      # 記事一覧
│   ├── ArticleForm.tsx      # 記事作成・編集フォーム
│   ├── ArticleDetail.tsx    # 記事詳細表示
│   └── ArticleStatus.tsx    # ステータス表示
├── ui/
│   ├── Button.tsx
│   ├── Input.tsx
│   └── TextArea.tsx
└── layout/
    ├── Header.tsx
    └── Navigation.tsx
```

## 機能要件

### 基本機能
1. **記事作成**: タイトル、内容、要約、タグの入力
2. **記事編集**: 既存記事の修正
3. **記事削除**: 記事の削除（確認ダイアログ付き）
4. **記事一覧**: ページネーション、フィルタリング
5. **記事詳細**: 個別記事の表示

### ステータス管理
- **Draft**: 下書き状態
- **Published**: 公開状態
- **Archived**: アーカイブ状態

### フィルタリング・ソート
- ステータス別フィルタ
- タグ別フィルタ
- 作成日・更新日・公開日でソート

## アーキテクチャ

### Clean Architecture維持

```
internal/
├── domain/
│   ├── entity/
│   │   └── article.go
│   └── repository/
│       └── article_repository.go
├── application/
│   └── usecase/
│       ├── create_article_usecase.go
│       ├── get_articles_usecase.go
│       ├── update_article_usecase.go
│       └── delete_article_usecase.go
├── infrastructure/
│   └── persistence/
│       └── postgres_article_repository.go
└── interfaces/
    ├── handler/
    │   └── article_handler.go
    └── dto/
        └── article_dto.go
```

## マイグレーション戦略

1. **Phase 1**: 記事テーブル作成
2. **Phase 2**: 記事CRUD API実装
3. **Phase 3**: フロントエンド実装
4. **Phase 4**: 不要なユーザー機能削除

## 将来の拡張予定

- Markdown対応
- 画像アップロード
- コメント機能
- カテゴリ機能
- 検索機能
- RSS配信