// 共通のAPI型定義

export interface Article {
  id: string;
  title: string;
  content: string;
  summary: string;
  status: 'draft' | 'published';
  tags: string[];
  article_image?: string;
  created_at: string;
  updated_at: string;
  published_at?: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  line_user_id?: string;
  profile_image?: string;
  created_at: string;
  updated_at: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface ArticlesResponse {
  data: Article[];
  pagination: PaginationMeta;
  meta: {
    sort: string;
    order: string;
  };
}

export interface ApiError {
  error: string;
  details?: string;
}

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';