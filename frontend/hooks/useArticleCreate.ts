import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { CreateArticleRequest } from '../types/article';
import { apiClient } from '../utils/apiClient';

export function useArticleCreate() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createArticle = async (data: CreateArticleRequest) => {
    if (!data.title.trim() || !data.content.trim()) {
      setError('タイトルと内容は必須です');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const article = await apiClient.post('/articles', data);
      alert('記事を作成しました！');
      router.push(`/articles/${article.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '記事の作成に失敗しました');
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