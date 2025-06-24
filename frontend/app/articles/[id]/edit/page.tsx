'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { useArticleEdit } from '../../../../hooks/useArticleEdit';
import ArticleForm from '../../../../components/ArticleForm';
import LoadingSpinner from '../../../../components/LoadingSpinner';

export default function EditArticlePage() {
  const params = useParams();
  const [mounted, setMounted] = useState(false);
  
  const {
    article,
    loading,
    saving,
    error,
    formData,
    updateArticle,
  } = useArticleEdit(params.id);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (loading) {
    return <LoadingSpinner message="記事を読み込み中..." />;
  }

  if (error && !article) {
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
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="grid gap-6">
            {/* ヘッダー */}
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold text-gray-900">記事を編集</h1>
              <div className="flex gap-2">
                <Link
                  href={`/articles/${params.id}`}
                  className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-medium"
                >
                  詳細に戻る
                </Link>
              </div>
            </div>

            {/* フォーム */}
            <ArticleForm
              mode="edit"
              initialData={formData}
              onSubmit={updateArticle}
              isSubmitting={saving}
              error={error}
              mounted={mounted}
            />
          </div>
        </div>
      </div>
    </div>
  );
}