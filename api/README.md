## アーキテクチャ

```
api/
├── cmd/
│   └── server/         # アプリケーションエントリポイント
├── internal/
│   ├── domain/         # ドメイン層（エンティティ、リポジトリ、サービス）
│   ├── application/    # アプリケーション層（ユースケース）
│   ├── infrastructure/ # インフラ層（外部サービス、永続化）
│   └── interfaces/     # インターフェース層（ハンドラ、ルーター）
└── pkg/
    └── config/         # 設定管理
```

### 各層の責務

1. **Domain層**: ビジネスロジックとエンティティ

   - Entity: ビジネスオブジェクト
   - Repository: データアクセスのインターフェース
   - Service: ドメインサービス

2. **Application層**: ユースケース

   - UseCase: アプリケーションのビジネスルール

3. **Infrastructure層**: 外部システムとの連携

   - Persistence: リポジトリの実装
   - External: 外部API連携

4. **Interfaces層**: プレゼンテーション
   - Handler: HTTPハンドラ
   - Router: ルーティング
   - Middleware: 共通処理

## Uber FXによる依存性注入

各層は`module.go`でFXモジュールとして定義され、`main.go`で組み立てられます：

```go
fx.New(
    config.Module,
    domain.Module,
    application.Module,
    infrastructure.Module,
    interfaces.Module,
)
```

## 実行方法

```bash
cd /home/sg3t41/workspace/api
go mod download
cd cmd/server
go run main.go
```

## API仕様

### ヘルスチェック

```bash
GET /health
```

### ユーザー作成

```bash
POST /api/v1/users
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe"
}
```

### ユーザー取得

```bash
GET /api/v1/users/{id}
```

## 環境変数

- `SERVER_ADDRESS`: サーバーアドレス（デフォルト: `:8080`）
- `GIN_MODE`: Ginのモード（デフォルト: `debug`）

## 特徴

- **Clean Architecture**: 依存関係が内側から外側への一方向
- **DDD**: ドメイン駆動設計の原則に従った構造
- **Dependency Injection**: Uber FXによる自動依存性注入
- **Testable**: 各層が独立しているためテストが容易
- **Scalable**: 新機能の追加が容易な構造
