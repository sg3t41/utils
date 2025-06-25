/**
 * 認証関連のユーティリティ関数
 */

/**
 * ローカルストレージからアクセストークンを取得
 */
export const getAccessToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('accessToken');
};

/**
 * 認証ヘッダーを作成
 */
export const createAuthHeaders = (): HeadersInit => {
  const token = getAccessToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
};

/**
 * APIのベースURLを取得
 */
export const getApiBaseUrl = (): string => {
  return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
};