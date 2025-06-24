/**
 * 共通APIクライアント
 * エラーハンドリングとレスポンス処理を統一化
 */

// APIエラーレスポンスの型定義
export interface ApiErrorResponse {
  error: string;
  code?: string;
  details?: any;
}

// APIエラークラス
export class ApiError extends Error {
  public status: number;
  public code?: string;
  public details?: any;

  constructor(message: string, status: number, code?: string, details?: any) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

// APIクライアント設定
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * APIリクエストのオプション
 */
export interface ApiRequestOptions extends RequestInit {
  includeAuth?: boolean;
  timeout?: number;
}

/**
 * 共通APIリクエスト関数
 */
export async function apiRequest<T = any>(
  endpoint: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const {
    includeAuth = true,
    timeout = 10000,
    headers: customHeaders = {},
    ...fetchOptions
  } = options;

  // デフォルトヘッダー
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...customHeaders,
  };

  // 認証トークンの追加
  if (includeAuth && typeof window !== 'undefined') {
    const token = localStorage.getItem('accessToken');
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
  }

  // タイムアウト設定
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...fetchOptions,
      headers,
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    // レスポンスボディを取得
    const responseText = await response.text();
    let responseData: any;

    try {
      responseData = responseText ? JSON.parse(responseText) : null;
    } catch (parseError) {
      console.warn('JSONパースエラー:', parseError);
      responseData = responseText;
    }

    // エラーレスポンスの処理
    if (!response.ok) {
      // 401エラーの場合は認証切れとしてログアウト処理
      if (response.status === 401 && typeof window !== 'undefined') {
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        localStorage.removeItem('user');
        // ページをリロードして認証状態をリセット
        setTimeout(() => window.location.reload(), 100);
      }
      
      const errorData = responseData as ApiErrorResponse;
      throw new ApiError(
        errorData?.error || `HTTP ${response.status}: ${response.statusText}`,
        response.status,
        errorData?.code,
        errorData?.details
      );
    }

    return responseData;
  } catch (error) {
    clearTimeout(timeoutId);

    if (error instanceof ApiError) {
      throw error;
    }

    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new ApiError('リクエストがタイムアウトしました', 408);
      }
      throw new ApiError(
        `ネットワークエラー: ${error.message}`,
        0
      );
    }

    throw new ApiError('不明なエラーが発生しました', 500);
  }
}

/**
 * GET リクエスト
 */
export async function get<T = any>(
  endpoint: string,
  options: Omit<ApiRequestOptions, 'method' | 'body'> = {}
): Promise<T> {
  return apiRequest<T>(endpoint, { ...options, method: 'GET' });
}

/**
 * POST リクエスト
 */
export async function post<T = any>(
  endpoint: string,
  data?: any,
  options: Omit<ApiRequestOptions, 'method' | 'body'> = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'POST',
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * PUT リクエスト
 */
export async function put<T = any>(
  endpoint: string,
  data?: any,
  options: Omit<ApiRequestOptions, 'method' | 'body'> = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'PUT',
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * DELETE リクエスト
 */
export async function del<T = any>(
  endpoint: string,
  options: Omit<ApiRequestOptions, 'method' | 'body'> = {}
): Promise<T> {
  return apiRequest<T>(endpoint, { ...options, method: 'DELETE' });
}

/**
 * PATCH リクエスト
 */
export async function patch<T = any>(
  endpoint: string,
  data?: any,
  options: Omit<ApiRequestOptions, 'method' | 'body'> = {}
): Promise<T> {
  return apiRequest<T>(endpoint, {
    ...options,
    method: 'PATCH',
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * エラーハンドリング用のヘルパー関数
 */
export function handleApiError(error: unknown): string {
  if (error instanceof ApiError) {
    switch (error.status) {
      case 400:
        return error.message || '入力データが無効です';
      case 401:
        return 'ログインが必要です';
      case 403:
        return 'アクセス権限がありません';
      case 404:
        return 'データが見つかりません';
      case 408:
        return 'リクエストがタイムアウトしました';
      case 409:
        return 'データが競合しています';
      case 422:
        return 'データの形式が正しくありません';
      case 429:
        return 'リクエスト数が上限を超えています';
      case 500:
        return 'サーバーでエラーが発生しました';
      case 502:
        return 'サーバーが一時的に利用できません';
      case 503:
        return 'サービスが一時的に利用できません';
      default:
        return error.message || '不明なエラーが発生しました';
    }
  }

  if (error instanceof Error) {
    return error.message;
  }

  return '不明なエラーが発生しました';
}

/**
 * レスポンスの型安全なキャスト
 */
export function castResponse<T>(data: unknown): T {
  return data as T;
}