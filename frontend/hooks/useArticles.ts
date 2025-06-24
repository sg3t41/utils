import { useState, useEffect, useCallback } from 'react';
import { Article, ArticleListResponse } from '../types/article';
import { get, handleApiError } from '../utils/apiClient';

interface UseArticlesParams {
  page: number;
  statusFilter: string;
  searchQuery: string;
}

interface UseArticlesReturn {
  articles: Article[];
  loading: boolean;
  error: string | null;
  pagination: ArticleListResponse['pagination'] | null;
  refetch: () => void;
}

/**
 * 記事データの取得とフィルタリングを管理するカスタムフック
 * ページネーション、ステータスフィルター、検索機能を含む
 */
export function useArticles({ page, statusFilter, searchQuery }: UseArticlesParams): UseArticlesReturn {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState<ArticleListResponse['pagination'] | null>(null);

  const fetchArticles = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '10',
      });
      
      if (statusFilter) {
        params.append('status', statusFilter);
      }
      
      if (searchQuery.trim()) {
        params.append('search', searchQuery.trim());
      }

      const data: ArticleListResponse = await get(`/api/v1/articles?${params}`);
      setArticles(data.data);
      setPagination(data.pagination);
    } catch (err) {
      setError(handleApiError(err));
      setArticles([]);
      setPagination(null);
    } finally {
      setLoading(false);
    }
  }, [page, statusFilter, searchQuery]);

  useEffect(() => {
    fetchArticles();
  }, [fetchArticles]);

  return {
    articles,
    loading,
    error,
    pagination,
    refetch: fetchArticles
  };
}