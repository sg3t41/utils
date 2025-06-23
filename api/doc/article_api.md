# ブログ記事管理API仕様書

## 概要

ブログ記事の作成、読取、更新、削除（CRUD）および公開管理を行うRESTful APIです。

## ベースURL

```
http://localhost:8080/api/v1
```

## 認証

現在は開発環境のため認証を無効化しています。本番環境ではJWTトークンによる認証が必要です。

## エンドポイント一覧

### 1. 記事一覧取得

```
GET /articles
```

#### クエリパラメータ

| パラメータ | 型 | 必須 | デフォルト | 説明 |
|------------|----|----|-----------|------|
| page | integer | No | 1 | ページ番号 |
| limit | integer | No | 10 | 1ページあたりの記事数（最大100） |
| sort | string | No | created_at | ソート項目（created_at, updated_at, published_at, title） |
| order | string | No | desc | ソート順序（asc, desc） |
| status | string | No | - | ステータス絞り込み（draft, published, archived） |
| search | string | No | - | タイトル・内容での検索 |
| tag | string | No | - | タグでの絞り込み |
| date_from | string | No | - | 作成日の開始日（YYYY-MM-DD） |
| date_to | string | No | - | 作成日の終了日（YYYY-MM-DD） |

#### レスポンス例

```json
{
  "data": [
    {
      "id": "4622b865-0acc-4ff0-b990-05b2e3986973",
      "title": "React開発のベストプラクティス",
      "summary": "React開発の効率化",
      "status": "draft",
      "tags": ["React", "JavaScript", "フロントエンド"],
      "thumbnail_image": null,
      "created_at": "2025-06-22T21:04:10Z",
      "published_at": null
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1,
    "has_next": false,
    "has_prev": false
  },
  "meta": {
    "sort": "created_at",
    "order": "desc"
  }
}
```

### 2. 記事作成

```
POST /articles
```

#### リクエストボディ

```json
{
  "title": "記事のタイトル",
  "content": "記事の内容",
  "summary": "記事の概要",
  "tags": ["タグ1", "タグ2"],
  "featured_image": "/uploads/images/featured.jpg"
}
```

#### レスポンス例

```json
{
  "id": "c7f7ccdf-72b7-4072-b8d4-e88bee10e306",
  "title": "記事のタイトル",
  "content": "記事の内容",
  "summary": "記事の概要",
  "status": "draft",
  "author_id": "b7c987cb-781a-4027-bb71-30d0a9d7cf14",
  "tags": ["タグ1", "タグ2"],
  "featured_image": "/uploads/images/featured.jpg",
  "thumbnail_image": null,
  "created_at": "2025-06-23T17:38:56Z",
  "updated_at": "2025-06-23T17:38:56Z",
  "published_at": null
}
```

### 3. 記事詳細取得

```
GET /articles/{id}
```

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|------------|----|----|------|
| id | string | Yes | 記事ID（UUID） |

### 4. 記事更新

```
PUT /articles/{id}
```

#### リクエストボディ

```json
{
  "title": "更新されたタイトル",
  "content": "更新された内容",
  "summary": "更新された概要",
  "tags": ["新しいタグ"]
}
```

### 5. 記事削除

```
DELETE /articles/{id}
```

### 6. 記事公開

```
POST /articles/{id}/publish
```

記事のステータスを`published`に変更し、`published_at`に現在時刻を設定します。

### 7. 記事非公開

```
POST /articles/{id}/unpublish
```

記事のステータスを`draft`に変更し、`published_at`をnullに設定します。

## エラーレスポンス

### バリデーションエラー

```json
{
  "code": "VALIDATION_ERROR",
  "error": "バリデーションエラー",
  "errors": [
    {
      "code": "required",
      "field": "Title",
      "message": "Title is required",
      "value": null
    }
  ]
}
```

### 記事が見つからない

```json
{
  "error": "記事が見つかりません"
}
```

### 内部サーバーエラー

```json
{
  "error": "内部サーバーエラーが発生しました"
}
```

## ステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | バリデーションエラー |
| 404 | リソースが見つからない |
| 500 | 内部サーバーエラー |

## データ型

### Article

| フィールド | 型 | 説明 |
|------------|----|----|
| id | string | 記事ID（UUID） |
| title | string | タイトル（必須、最大500文字） |
| content | string | 本文（必須） |
| summary | string | 概要（最大1000文字） |
| status | string | ステータス（draft, published, archived） |
| author_id | string | 作成者ID（UUID） |
| tags | string[] | タグ配列 |
| featured_image | string | アイキャッチ画像パス |
| thumbnail_image | string | サムネイル画像パス |
| created_at | string | 作成日時（ISO 8601） |
| updated_at | string | 更新日時（ISO 8601） |
| published_at | string | 公開日時（ISO 8601、未公開の場合null） |