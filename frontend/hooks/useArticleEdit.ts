import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Article, UpdateArticleRequest } from '../types/article';
import { get, put } from '../utils/apiClient';
import { ERROR_MESSAGES } from '../utils/constants';

export function useArticleEdit(articleId: string | string[]) {
  const router = useRouter();
  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const [formData, setFormData] = useState<UpdateArticleRequest>({
    title: '',
    content: '',
    summary: '',
    tags: [],
    article_image: '',
  });

  const fetchArticle = async () => {
    try {
      setLoading(true);
      const data = await get<Article>(`/api/v1/articles/${articleId}`);
      setArticle(data);
      setFormData({
        title: data.title,
        content: data.content,
        summary: data.summary,
        tags: data.tags,
        article_image: data.article_image,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (articleId) {
      fetchArticle();
    }
  }, [articleId]);

  const updateArticle = async (data: UpdateArticleRequest) => {
    if (!data.title?.trim() || !data.content?.trim()) {
      setError(ERROR_MESSAGES.REQUIRED_FIELDS);
      return;
    }

    try {
      setSaving(true);
      setError(null);
      
      await put(`/api/v1/articles/${articleId}`, data);
      alert('記事を更新しました！');
      router.push(`/articles/${articleId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '記事の更新に失敗しました');
    } finally {
      setSaving(false);
    }
  };

  return {
    article,
    loading,
    saving,
    error,
    formData,
    setFormData,
    updateArticle,
    setError,
  };
}