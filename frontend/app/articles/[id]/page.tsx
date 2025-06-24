'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '../../../contexts/AuthContext';
import { get, post, del, handleApiError } from '../../../utils/apiClient';
import { API_BASE_URL } from '../../../types/api';
import { formatDateTime } from '../../../utils/dateFormat';
import { getStatusBadge, getStatusText } from '../../../utils/statusUtils';

interface Article {
  id: string;
  title: string;
  content: string;
  summary: string;
  status: string;
  author_id: string;
  tags: string[];
  article_image?: string;
  created_at: string;
  updated_at: string;
  published_at: string | null;
}

export default function ArticleDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { isAdmin } = useAuth();
  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const fetchArticle = async () => {
    try {
      setLoading(true);
      const data: Article = await get(`/api/v1/articles/${params.id}`);
      setArticle(data);
    } catch (err) {
      setError(handleApiError(err));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (params.id) {
      fetchArticle();
    }
  }, [params.id]);

  const handlePublish = async () => {
    if (!article) return;
    
    try {
      setActionLoading('publish');
      const updatedArticle: Article = await post(`/api/v1/articles/${article.id}/publish`, {});
      setArticle(updatedArticle);
    } catch (err) {
      alert(handleApiError(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleUnpublish = async () => {
    if (!article) return;
    
    try {
      setActionLoading('unpublish');
      const updatedArticle: Article = await post(`/api/v1/articles/${article.id}/unpublish`);
      setArticle(updatedArticle);
    } catch (err) {
      alert(handleApiError(err));
    } finally {
      setActionLoading(null);
    }
  };

  const handleDelete = async () => {
    if (!article) return;
    
    if (!confirm('この記事を削除しますか？この操作は取り消せません。')) {
      return;
    }

    try {
      setActionLoading('delete');
      await del(`/api/v1/articles/${article.id}`);
      alert('記事を削除しました');
      router.push('/articles');
    } catch (err) {
      alert(handleApiError(err));
    } finally {
      setActionLoading(null);
    }
  };


  if (loading) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">記事を読み込み中...</p>
        </div>
      </div>
    );
  }

  if (error || !article) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">記事が見つかりません</h1>
          <p className="text-gray-600 mb-6">{error}</p>
          <Link
            href="/articles"
            className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium inline-block"
          >
            記事一覧に戻る
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pt-24 pb-8">
      <div className="max-w-6xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          {/* ヘッダー画像 */}
          {article.article_image && (
            <div className="w-full h-64 md:h-96 bg-gray-100">
              <img
                src={mounted ? `${API_BASE_URL}/api/v1/uploads/${article.article_image}` : ''}
                alt={article.title}
                className="w-full h-full object-cover"
                onError={(e) => {
                  e.currentTarget.parentElement.style.display = 'none';
                }}
              />
            </div>
          )}
          
          {/* ヘッダー */}
          <div className="border-b border-gray-200 p-6">
            <div className="grid gap-4">
              <div className="flex justify-between items-start gap-4">
                <div className="flex-1">
                  <h1 className="text-3xl font-bold text-gray-900 mb-2">{article.title}</h1>
                  <div className="flex items-center gap-4">
                    <span className={getStatusBadge(article.status)}>
                      {getStatusText(article.status)}
                    </span>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Link
                    href="/articles"
                    className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-medium inline-block"
                  >
                    一覧に戻る
                  </Link>
                  {isAdmin && (
                    <Link
                      href={`/articles/${article.id}/edit`}
                      className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium inline-block"
                    >
                      編集
                    </Link>
                  )}
                </div>
              </div>

              {/* アクションボタン */}
              {isAdmin && (
                <div className="flex gap-2">
                  {article.status === 'draft' ? (
                    <button
                      onClick={handlePublish}
                      disabled={actionLoading === 'publish'}
                      className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg font-medium disabled:opacity-50"
                    >
                      {actionLoading === 'publish' ? '公開中...' : '公開する'}
                    </button>
                  ) : (
                    <button
                      onClick={handleUnpublish}
                      disabled={actionLoading === 'unpublish'}
                      className="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg font-medium disabled:opacity-50"
                    >
                      {actionLoading === 'unpublish' ? '処理中...' : '下書きに戻す'}
                    </button>
                  )}
                  <button
                    onClick={handleDelete}
                    disabled={actionLoading === 'delete'}
                    className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg font-medium disabled:opacity-50"
                  >
                    {actionLoading === 'delete' ? '削除中...' : '削除'}
                  </button>
                </div>
              )}

              {/* メタ情報 */}
              <div className="grid gap-2 text-sm text-gray-600">
                <div className="flex justify-between">
                  <span>作成日時: {formatDateTime(article.created_at)}</span>
                  <span>更新日時: {formatDateTime(article.updated_at)}</span>
                </div>
                {article.published_at && (
                  <div>
                    <span>公開日時: {formatDateTime(article.published_at)}</span>
                  </div>
                )}
              </div>

              {/* タグ */}
              {article.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {article.tags.map((tag, index) => (
                    <span
                      key={index}
                      className="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              )}


              {/* 概要 */}
              {article.summary && (
                <div className="bg-gray-50 p-4 rounded-lg">
                  <h3 className="font-medium text-gray-900 mb-2">概要</h3>
                  <p className="text-gray-700">{article.summary}</p>
                </div>
              )}
            </div>
          </div>

          {/* 記事内容 */}
          <div className="p-6">
            <div className="prose max-w-none">
              <div className="whitespace-pre-wrap text-gray-900 leading-relaxed">
                {article.content}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}