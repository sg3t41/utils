'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useAuth } from '../../contexts/AuthContext';
import { useArticles } from '../../hooks/useArticles';
import { useArticleFilters } from '../../hooks/useArticleFilters';
import SearchFilter from '../../components/SearchFilter';
import ArticleCard from '../../components/ArticleCard';
import Pagination from '../../components/Pagination';
import ErrorMessage from '../../components/ErrorMessage';
import LoadingSpinner from '../../components/LoadingSpinner';

export default function ArticlesPage() {
  const { isAuthenticated, isAdmin } = useAuth();
  const [mounted, setMounted] = useState(false);
  
  // フィルター状態管理
  const {
    statusFilter,
    searchQuery,
    activeSearchQuery,
    page,
    setStatusFilter,
    setSearchQuery,
    setPage,
    handleSearch,
    handleClearSearch,
  } = useArticleFilters();
  
  // 記事データ取得
  const { articles, loading, error, pagination } = useArticles({
    page,
    statusFilter,
    searchQuery: activeSearchQuery,
  });

  useEffect(() => {
    setMounted(true);
  }, []);




  if (loading && articles.length === 0) {
    return <LoadingSpinner message="記事を読み込み中..." />;
  }

  return (
    <div className="min-h-screen bg-gray-50 pt-24 pb-8">
      <div className="max-w-6xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="grid gap-6">
            {/* ヘッダー */}
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold text-gray-900">記事一覧</h1>
              {isAdmin && (
                <Link
                  href="/articles/new"
                  className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg font-medium inline-block"
                >
                  新規作成
                </Link>
              )}
            </div>

            {/* フィルターと検索 */}
            <SearchFilter
              statusFilter={statusFilter}
              searchQuery={searchQuery}
              activeSearchQuery={activeSearchQuery}
              onStatusChange={setStatusFilter}
              onSearchQueryChange={setSearchQuery}
              onSearchSubmit={handleSearch}
              onClearSearch={handleClearSearch}
            />

            {/* エラー表示 */}
            {error && <ErrorMessage message={error} />}

            {/* 記事一覧 */}
            <div className="grid gap-4">
              {articles.length === 0 ? (
                <div className="text-center py-12">
                  <p className="text-gray-500">記事が見つかりませんでした。</p>
                </div>
              ) : (
                articles.map((article) => (
                  <ArticleCard
                    key={article.id}
                    article={article}
                    mounted={mounted}
                  />
                ))
              )}
            </div>

            {/* ページネーション */}
            {pagination && (
              <Pagination
                currentPage={pagination.page}
                totalPages={pagination.total_pages}
                hasNext={pagination.has_next}
                hasPrev={pagination.has_prev}
                onPageChange={setPage}
                totalItems={pagination.total}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}