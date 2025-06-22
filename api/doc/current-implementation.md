# 現在の実装詳細

## 実装済み機能

### Userエンティティの管理

現在、ユーザーの作成・取得機能が実装されています。

## ファイル別実装詳細

### ドメイン層

#### `internal/domain/entity/user.go`
- **User構造体**: ID, Email, Name, CreatedAt, UpdatedAt
- **NewUser関数**: ユーザー作成時のバリデーション（Email・Name必須）
- **UpdateName関数**: 名前更新時のバリデーションと更新日時自動設定

#### `internal/domain/repository/user_repository.go`
- **UserRepositoryインターフェース**: CRUD操作の抽象化
  - `Create`: ユーザー作成
  - `FindByID`: ID検索
  - `FindByEmail`: Email検索
  - `Update`: ユーザー更新
  - `Delete`: ユーザー削除
  - `List`: ページネーション付きリスト取得

#### `internal/domain/service/user_service.go`
- **UserService**: ビジネスロジックを担当
  - `CreateUser`: Email重複チェック付きユーザー作成
  - `GetUser`: ユーザー取得
  - `UpdateUserName`: ユーザー名更新（未使用）

### アプリケーション層

#### `internal/application/usecase/create_user_usecase.go`
- **CreateUserInput**: Email, Name
- **CreateUserOutput**: User
- **CreateUserUseCase**: ユーザー作成の協調処理

#### `internal/application/usecase/get_user_usecase.go`
- **GetUserInput**: ID
- **GetUserOutput**: User  
- **GetUserUseCase**: ユーザー取得の協調処理

### インフラストラクチャ層

#### `internal/infrastructure/persistence/memory_user_repository.go`
- **MemoryUserRepository**: インメモリ実装
  - `users map[string]*entity.User`: ID→User マッピング
  - `index map[string]*entity.User`: Email→User マッピング
  - goroutine-safe（sync.RWMutex使用）
  - 完全なCRUD操作実装済み

### インターフェース層

#### `internal/interfaces/handler/user_handler.go`
- **CreateUserRequest**: JSON入力バリデーション（email必須・形式チェック、name必須）
- **UserResponse**: JSON出力形式（ISO形式の日時）
- **CreateUser**: POST /api/v1/users
- **GetUser**: GET /api/v1/users/:id

#### `internal/interfaces/router/router.go`
- **Ginフレームワーク**使用
- **ミドルウェア設定**: Logger, Recovery, CORS
- **ルーティング設定**:
  - `GET /health` - ヘルスチェック
  - `POST /api/v1/users` - ユーザー作成
  - `GET /api/v1/users/:id` - ユーザー取得

#### `internal/interfaces/middleware/`
- **CORS**: クロスオリジンリクエスト対応
- **Logger**: リクエスト・レスポンスログ
- **Recovery**: パニック時の復旧処理

### 設定・起動

#### `pkg/config/`
- **Config構造体**: サーバー設定（ポート、モードなど）
- **Logger設定**: Zap logger初期化

#### `cmd/server/main.go`
- **Uber FX**による依存性注入
- **ライフサイクル管理**: 起動・停止フック
- **モジュール構成**: config → domain → application → infrastructure → interfaces

## データフロー

### ユーザー作成
```
POST /api/v1/users
↓
UserHandler.CreateUser (バリデーション)
↓
CreateUserUseCase.Execute
↓
UserService.CreateUser (重複チェック)
↓
entity.NewUser (ビジネスルール検証)
↓
MemoryUserRepository.Create
```

### ユーザー取得
```
GET /api/v1/users/:id
↓
UserHandler.GetUser
↓
GetUserUseCase.Execute
↓
UserService.GetUser
↓
MemoryUserRepository.FindByID
```

## 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: Gin
- **DI**: Uber FX
- **ログ**: Zap
- **バリデーション**: Gin binding
- **UUID**: Google UUID

## 現在の制限事項

1. **データ永続化**: インメモリのみ（再起動で消失）
2. **認証・認可**: 未実装
3. **エラーハンドリング**: 基本的なもののみ
4. **テスト**: 未実装
5. **API仕様書**: 未作成
6. **ログレベル設定**: 固定
7. **ヘルスチェック**: 簡易実装のみ

## 次の実装候補

1. **データベース実装**: PostgreSQL/MySQL対応
2. **ユーザー更新・削除API**: PUT/DELETE エンドポイント
3. **認証機能**: JWT実装
4. **バリデーション強化**: カスタムバリデータ
5. **エラーレスポンス統一**: RFC7807準拠
6. **テスト実装**: 単体・結合テスト
7. **OpenAPI仕様**: Swagger生成
8. **ページネーション**: リスト取得API
9. **検索機能**: 複数条件検索
10. **メトリクス・トレーシング**: 監視機能