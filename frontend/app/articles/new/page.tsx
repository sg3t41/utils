'use client';

import { useState, useEffect } from 'react';
import { useArticleCreate } from '../../../hooks/useArticleCreate';
import ArticleForm from '../../../components/ArticleForm';

export default function NewArticlePage() {
  const [mounted, setMounted] = useState(false);
  const { loading, error, createArticle } = useArticleCreate();

  useEffect(() => {
    setMounted(true);
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 pt-24 pb-8">
      <div className="max-w-6xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="grid gap-6">
            {/* ヘッダー */}
            <h1 className="text-3xl font-bold text-gray-900">新しい記事を作成</h1>

            {/* フォーム */}
            <ArticleForm
              mode="create"
              onSubmit={createArticle}
              isSubmitting={loading}
              error={error}
              mounted={mounted}
            />
          </div>
        </div>
      </div>
    </div>
  );
}