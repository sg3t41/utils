# DDDアーキテクチャ概要

## プロジェクトの構成

このプロジェクトはDDD（ドメイン駆動設計）の4層アーキテクチャに基づいて設計されています。

```
internal/
├── domain/          # ドメイン層
├── application/     # アプリケーション層
├── infrastructure/  # インフラストラクチャ層
└── interfaces/      # インターフェース層
```

## 各層の責務

### 1. ドメイン層 (`internal/domain/`)

ビジネスロジックの中核を担当する層です。外部の技術的詳細に依存しません。

```
domain/
├── entity/          # エンティティ
│   └── user.go     # Userエンティティ
├── repository/      # リポジトリインターフェース
│   └── user_repository.go
├── service/         # ドメインサービス
│   └── user_service.go
└── module.go       # DIモジュール
```

**エンティティ (Entity)**

- `User`: ビジネスの核となるユーザーオブジェクト
  - ID、Email、Name、作成日時、更新日時を管理
  - ビジネスルール（バリデーション）を内包
  - 不変性を保持（状態変更は専用メソッドを通じて実行）

**リポジトリインターフェース**

- データ永続化の抽象化
- 具体的な実装はインフラストラクチャ層に委譲

**ドメインサービス**

- 複数のエンティティにまたがるビジネスロジック
- 重複チェックなどの複雑な処理を担当

### 2. アプリケーション層 (`internal/application/`)

ユースケースを実装し、ドメイン層を調整する役割を持ちます。

```
application/
├── usecase/
│   ├── create_user_usecase.go  # ユーザー作成ユースケース
│   └── get_user_usecase.go     # ユーザー取得ユースケース
└── module.go
```

**ユースケース (Use Case)**

- 具体的なビジネス操作を表現
- 入力・出力の構造体を定義
- ドメインサービスを組み合わせてアプリケーション固有の処理を実現

### 3. インフラストラクチャ層 (`internal/infrastructure/`)

外部システムとの統合を担当します。

```
infrastructure/
├── persistence/
│   └── memory_user_repository.go  # インメモリ実装
├── external/        # 外部API連携（将来用）
└── module.go
```

**リポジトリ実装**

- ドメイン層で定義されたインターフェースの具体実装
- 現在はインメモリ実装のみ（将来的にDB実装に切り替え可能）

### 4. インターフェース層 (`internal/interfaces/`)

外部からのリクエストを受け付け、適切な形式でレスポンスを返します。

```
interfaces/
├── handler/         # HTTPハンドラー
│   └── user_handler.go
├── middleware/      # ミドルウェア
│   ├── cors.go
│   ├── logger.go
│   └── recovery.go
├── router/          # ルーティング
│   └── router.go
└── module.go
```

**ハンドラー**

- HTTP リクエスト・レスポンスの処理
- リクエストデータの検証
- ユースケースの呼び出し
- レスポンス形式の変換

## 依存関係の方向

```
interfaces → application → domain ← infrastructure
```

- 上位層は下位層に依存可能
- 下位層は上位層に依存してはいけない
- インフラストラクチャ層のみドメイン層のインターフェースに依存

## 現在実装されているAPI

### エンドポイント

- `GET /health` - ヘルスチェック
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users/:id` - ユーザー取得

### ユーザー作成フロー

1. `UserHandler.CreateUser` がリクエストを受信
2. `CreateUserUseCase` を実行
3. `UserService.CreateUser` でビジネスロジック処理
4. `UserRepository.Create` でデータ永続化
5. レスポンス返却

## 設定とDI

- **Uber FX** を使用した依存性注入
- 各層に `module.go` でDIコンテナを定義
- 設定は `pkg/config/` で管理

## 今後の拡張ポイント

1. **データベース実装**: `infrastructure/persistence/` にDB実装を追加
2. **新しいエンティティ**: `domain/entity/` に追加
3. **新しいユースケース**: `application/usecase/` に追加
4. **外部API連携**: `infrastructure/external/` に実装
5. **認証・認可**: ミドルウェアまたは新しい層として追加

## ディレクトリルール

- 各層は他層の具象実装に直接依存しない
- インターフェースを通じた疎結合を維持
- テストは各層で独立して書ける構造
- 新機能追加は既存の層構造に従って配置
