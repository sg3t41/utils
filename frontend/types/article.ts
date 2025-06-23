// 記事関連の型定義
export interface Article {
  id: string;
  title: string;
  summary: string;
  content?: string;
  status: 'draft' | 'published' | 'archived';
  tags: string[];
  article_image?: string;
  created_at: string;
  updated_at?: string;
  published_at?: string;
}

export interface ArticleListResponse {
  data: Article[];
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

export interface CreateArticleRequest {
  title: string;
  summary: string;
  content: string;
  status: 'draft' | 'published';
  tags: string[];
  article_image?: string;
}

export interface UpdateArticleRequest extends CreateArticleRequest {
  id: string;
}

export type ArticleStatus = 'draft' | 'published' | 'archived';