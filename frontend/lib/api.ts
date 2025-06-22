// サーバーサイドとクライアントサイドで異なるURLを使用
const API_BASE_URL = typeof window === 'undefined' 
  ? (process.env.API_URL || 'http://utils_api:8080')  // サーバーサイド（コンテナ内）
  : (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080');  // クライアントサイド（ブラウザ）

export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface UsersResponse {
  data: User[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
  meta: {
    sort: string;
    order: string;
  };
}

export interface CreateUserRequest {
  name: string;
  email: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status} ${response.statusText}`);
      }

      // DELETEリクエストで204 No Contentの場合はレスポンスボディが空なのでJSONパースしない
      if (response.status === 204 || response.headers.get('content-length') === '0') {
        return null as T;
      }

      const text = await response.text();
      return text ? JSON.parse(text) : (null as T);
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
        throw new Error('Unable to connect to API. Please check if the server is running.');
      }
      throw error;
    }
  }

  async getUsers(): Promise<User[]> {
    const response = await this.request<UsersResponse>('/api/v1/users?page=1&limit=100');
    return response.data;
  }

  async getUserById(id: string): Promise<User> {
    return this.request<User>(`/api/v1/users/${id}`);
  }

  async createUser(userData: CreateUserRequest): Promise<User> {
    return this.request<User>('/api/v1/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(id: string, hard: boolean = true): Promise<void> {
    const url = `/api/v1/users/${id}${hard ? '?hard=true' : ''}`;
    await this.request<void>(url, {
      method: 'DELETE',
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);