// API関連の定数
export const API_ENDPOINTS = {
  UPLOAD_IMAGE: '/api/v1/upload/image',
  UPLOADS: '/api/v1/uploads',
} as const;

// デフォルトのAPI URL
export const DEFAULT_API_URL = 'http://localhost:8080';

// エラーメッセージ
export const ERROR_MESSAGES = {
  UPLOAD_FAILED: '画像のアップロードに失敗しました',
  UPLOAD_SUCCESS: '画像をアップロードしました',
  REQUIRED_FIELDS: 'タイトルと内容は必須です',
  ARTICLE_CREATE_FAILED: '記事の作成に失敗しました',
  ARTICLE_CREATE_SUCCESS: '記事を作成しました！',
} as const;

// スタイルクラス
export const CSS_CLASSES = {
  INPUT: 'w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
  LABEL: 'block text-sm font-medium text-gray-700 mb-2',
  ERROR_TEXT: 'text-xs text-red-600 mt-1',
  INFO_TEXT: 'text-xs text-blue-600 mt-1',
  IMAGE_PREVIEW: 'w-full max-w-md h-48 object-cover rounded-lg shadow-sm',
} as const;