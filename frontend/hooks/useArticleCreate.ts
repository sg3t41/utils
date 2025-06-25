import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { CreateArticleRequest } from '../types/article';
import { post } from '../utils/apiClient';
import { ERROR_MESSAGES } from '../utils/constants';

/**
 * 記事作成用のカスタムフック
 */
export function useArticleCreate() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const validateArticleData = (data: CreateArticleRequest): boolean => {
    if (!data.title.trim() || !data.content.trim()) {
      setError(ERROR_MESSAGES.REQUIRED_FIELDS);
      return false;
    }
    return true;
  };

  const createArticle = async (data: CreateArticleRequest) => {
    if (!validateArticleData(data)) return;

    try {
      setLoading(true);
      setError(null);
      
      const article = await post('/api/v1/articles', data);
      alert(ERROR_MESSAGES.ARTICLE_CREATE_SUCCESS);
      router.push(`/articles/${article.id}`);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.ARTICLE_CREATE_FAILED;
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return {
    loading,
    error,
    createArticle,
    setError,
  };
}