'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';

interface Article {
  id: string;
  title: string;
  summary: string;
  status: string;
  tags: string[];
  article_image?: string;
  created_at: string;
  published_at: string | null;
}

interface ArticlesResponse {
  data: Article[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}

export default function ArticlesPage() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [pagination, setPagination] = useState<ArticlesResponse['pagination'] | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [activeSearchQuery, setActiveSearchQuery] = useState<string>('');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const fetchArticles = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '10',
      });
      
      if (statusFilter) {
        params.append('status', statusFilter);
      }
      
      if (activeSearchQuery.trim()) {
        params.append('search', activeSearchQuery.trim());
      }

      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/articles?${params}`);
      if (!response.ok) {
        throw new Error('Failed to fetch articles');
      }

      const data: ArticlesResponse = await response.json();
      setArticles(data.data);
      setPagination(data.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArticles();
  }, [page, statusFilter, activeSearchQuery]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    setActiveSearchQuery(searchQuery.trim());
  };

  const handleClearSearch = () => {
    setSearchQuery('');
    setActiveSearchQuery('');
    setPage(1);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusBadge = (status: string) => {
    const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
    switch (status) {
      case 'published':
        return `${baseClasses} bg-green-100 text-green-800`;
      case 'draft':
        return `${baseClasses} bg-yellow-100 text-yellow-800`;
      case 'archived':
        return `${baseClasses} bg-gray-100 text-gray-800`;
      default:
        return `${baseClasses} bg-gray-100 text-gray-800`;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'published':
        return '公開済み';
      case 'draft':
        return '下書き';
      case 'archived':
        return 'アーカイブ';
      default:
        return status;
    }
  };

  if (loading && articles.length === 0) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">記事を読み込み中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="grid gap-6">
            {/* ヘッダー */}
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold text-gray-900">記事一覧</h1>
              <Link
                href="/articles/new"
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium inline-block"
              >
                新規作成
              </Link>
            </div>

            {/* フィルターと検索 */}
            <div className="grid gap-4 md:grid-cols-2">
              <div>
                <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-2">
                  ステータス
                </label>
                <select
                  id="status"
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="">すべて</option>
                  <option value="published">公開済み</option>
                  <option value="draft">下書き</option>
                  <option value="archived">アーカイブ</option>
                </select>
              </div>
              <form onSubmit={handleSearch}>
                <label htmlFor="search" className="block text-sm font-medium text-gray-700 mb-2">
                  検索
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    id="search"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="タイトルや内容を検索..."
                    className="flex-1 border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <button
                    type="submit"
                    className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg"
                  >
                    検索
                  </button>
                  {activeSearchQuery.trim() && (
                    <button
                      type="button"
                      onClick={handleClearSearch}
                      className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg"
                    >
                      クリア
                    </button>
                  )}
                </div>
              </form>
            </div>

            {/* 現在の検索条件表示 */}
            {(activeSearchQuery.trim() || statusFilter) && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-blue-800">
                  検索条件: 
                  {activeSearchQuery.trim() && <span className="ml-2 font-medium">「{activeSearchQuery.trim()}」</span>}
                  {statusFilter && <span className="ml-2 font-medium">ステータス: {getStatusText(statusFilter)}</span>}
                </p>
              </div>
            )}

            {/* エラー表示 */}
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-red-800">エラー: {error}</p>
              </div>
            )}

            {/* 記事一覧 */}
            <div className="grid gap-4">
              {articles.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-gray-500">記事が見つかりませんでした。</p>
                </div>
              ) : (
                articles.map((article) => (
                  <div
                    key={article.id}
                    className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow"
                  >
                    <div className="flex gap-4">
                      {article.article_image && (
                        <div className="w-40 h-28 flex-shrink-0">
                          <img
                            src={mounted ? `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/uploads/${article.article_image}` : ''}
                            alt={article.title}
                            className="w-full h-full object-cover rounded-lg shadow-sm"
                            onError={(e) => {
                              e.currentTarget.style.display = 'none';
                            }}
                          />
                        </div>
                      )}
                      <div className="flex-1 grid gap-3">
                        <div className="flex justify-between items-start gap-4">
                          <h2 className="text-xl font-semibold text-gray-900 hover:text-blue-600">
                            <Link href={`/articles/${article.id}`}>{article.title}</Link>
                          </h2>
                          <span className={getStatusBadge(article.status)}>
                            {getStatusText(article.status)}
                          </span>
                        </div>
                      
                      {article.summary && (
                        <p className="text-gray-600 line-clamp-2">{article.summary}</p>
                      )}
                      
                      {article.tags.length > 0 && (
                        <div className="flex flex-wrap gap-2">
                          {article.tags.map((tag, index) => (
                            <span
                              key={index}
                              className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                            >
                              {tag}
                            </span>
                          ))}
                        </div>
                      )}
                      
                        <div className="flex justify-between items-center text-sm text-gray-500">
                          <span>作成: {formatDate(article.created_at)}</span>
                          {article.published_at && (
                            <span>公開: {formatDate(article.published_at)}</span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>

            {/* ページネーション */}
            {pagination && pagination.total_pages > 1 && (
              <div className="flex justify-center items-center gap-4">
                <button
                  onClick={() => setPage(page - 1)}
                  disabled={!pagination.has_prev}
                  className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
                >
                  前へ
                </button>
                <span className="text-gray-600">
                  {pagination.page} / {pagination.total_pages} ページ
                </span>
                <button
                  onClick={() => setPage(page + 1)}
                  disabled={!pagination.has_next}
                  className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
                >
                  次へ
                </button>
              </div>
            )}

            {/* 統計情報 */}
            {pagination && (
              <div className="text-center text-sm text-gray-500">
                総件数: {pagination.total}件
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}