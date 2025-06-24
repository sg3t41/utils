# フロントエンド

[Next.js](https://nextjs.org)で構築された個人用ダッシュボードのフロントエンドアプリケーションです。

## 技術スタック

- **Framework**: Next.js 15.3.4 (App Router)
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS
- **UI コンポーネント**: React 19
- **認証**: カスタム認証コンテキスト（LINE Login対応）
- **開発環境**: ESLint, Prettier

## 開発環境のセットアップ

開発サーバーを起動:

```bash
npm run dev
# または
yarn dev
# または
pnpm dev
# または
bun dev
```

ブラウザで [http://localhost:3000](http://localhost:3000) を開いて結果を確認してください。

`app/page.tsx` を編集することでページを変更できます。ファイルを編集すると自動的にページが更新されます。

## プロジェクト構成

```
frontend/
├── app/                    # App Router ページ
│   ├── articles/           # ブログ記事関連ページ
│   ├── auth/              # 認証関連ページ
│   ├── users/             # ユーザー管理ページ
│   ├── layout.tsx         # ルートレイアウト
│   └── page.tsx           # ホームページ（タイルダッシュボード）
├── components/            # 再利用可能なコンポーネント
│   ├── ArticleCard.tsx    # 記事カード
│   ├── ArticleForm.tsx    # 記事作成・編集フォーム
│   ├── FloatingMenuButton.tsx # ハンバーガーメニュー
│   ├── ImageUpload.tsx    # 画像アップロード
│   ├── Pagination.tsx     # ページネーション
│   └── ...
├── contexts/              # React Context
│   └── AuthContext.tsx    # 認証状態管理
├── hooks/                 # カスタムフック
│   ├── useArticles.ts     # 記事管理フック
│   ├── useScrollHeader.ts # スクロール状態管理
│   └── ...
├── types/                 # TypeScript型定義
│   ├── api.ts             # API関連型
│   └── article.ts         # 記事関連型
└── utils/                 # ユーティリティ関数
    ├── apiClient.ts       # API クライアント
    ├── dateFormat.ts      # 日付フォーマット
    └── statusUtils.ts     # ステータス関連
```

## 主な機能

### 実装済み機能

1. **タイルダッシュボード**
   - Windows 10風のタイル型レイアウト
   - レスポンシブデザイン（モバイル・タブレット・デスクトップ対応）
   - ホバーアニメーション

2. **ブログシステム**
   - 記事一覧・詳細・作成・編集ページ
   - 画像アップロード機能
   - 公開・下書き状態管理
   - 検索・フィルタリング機能

3. **認証システム**
   - LINE Login 対応
   - JWT トークンベース認証
   - 管理者権限制御

4. **UI/UX**
   - ハンバーガーメニューナビゲーション
   - ページネーション
   - エラーハンドリング
   - ローディング状態表示

### 計画中の機能

- 制作物管理ページ
- 収支表管理ページ
- ランキング管理ページ
- ファイル共有ページ
- フォトギャラリーページ
- メモ機能ページ
- リンク集ページ

## 開発ガイドライン

### コーディング規約

- **スタイル**: Tailwind CSS を使用、grid-layout を優先
- **コンポーネント**: 関数コンポーネント + TypeScript
- **状態管理**: React Context + カスタムフック
- **ファイル命名**: PascalCase（コンポーネント）、camelCase（その他）

### ディレクトリ構成ルール

- `app/`: ページコンポーネント（App Router）
- `components/`: 再利用可能なUIコンポーネント
- `hooks/`: カスタムフック
- `contexts/`: グローバル状態管理
- `utils/`: ユーティリティ関数
- `types/`: TypeScript型定義

## APIとの連携

バックエンドAPI（http://localhost:8080）との通信は`utils/apiClient.ts`を通じて行います。

- 自動JWT認証ヘッダー付与
- エラーハンドリング（401エラー時の自動ログアウト）
- レスポンスの型安全性確保

## 参考資料

Next.jsについて詳しく学ぶには以下のリソースをご覧ください:

- [Next.js Documentation](https://nextjs.org/docs) - Next.jsの機能とAPI
- [Learn Next.js](https://nextjs.org/learn) - インタラクティブなNext.jsチュートリアル
- [Tailwind CSS Documentation](https://tailwindcss.com/docs) - Tailwind CSSの使用方法
