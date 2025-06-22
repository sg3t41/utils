'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface CreateArticleRequest {
  title: string;
  content: string;
  summary: string;
  tags: string[];
}

export default function NewArticlePage() {
  const router = useRouter();
  const [formData, setFormData] = useState<CreateArticleRequest>({
    title: '',
    content: '',
    summary: '',
    tags: [],
  });
  const [tagInput, setTagInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title.trim() || !formData.content.trim()) {
      setError('タイトルと内容は必須です');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await fetch('http://localhost:8080/api/v1/articles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '記事の作成に失敗しました');
      }

      const article = await response.json();
      alert('記事を作成しました！');
      router.push(`/articles/${article.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '記事の作成に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleAddTag = () => {
    const tag = tagInput.trim();
    if (tag && !formData.tags.includes(tag)) {
      setFormData({
        ...formData,
        tags: [...formData.tags, tag],
      });
      setTagInput('');
    }
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setFormData({
      ...formData,
      tags: formData.tags.filter(tag => tag !== tagToRemove),
    });
  };

  const handleTagInputKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="bg-white rounded-lg shadow-sm p-6">
          <div className="grid gap-6">
            {/* ヘッダー */}
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold text-gray-900">新しい記事を作成</h1>
              <Link
                href="/articles"
                className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-medium inline-block"
              >
                キャンセル
              </Link>
            </div>

            {/* エラー表示 */}
            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <p className="text-red-800">{error}</p>
              </div>
            )}

            {/* フォーム */}
            <form onSubmit={handleSubmit} className="grid gap-6">
              {/* タイトル */}
              <div>
                <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                  タイトル <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  id="title"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder="記事のタイトルを入力してください"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  maxLength={500}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">{formData.title.length}/500文字</p>
              </div>

              {/* 概要 */}
              <div>
                <label htmlFor="summary" className="block text-sm font-medium text-gray-700 mb-2">
                  概要
                </label>
                <textarea
                  id="summary"
                  value={formData.summary}
                  onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                  placeholder="記事の概要を入力してください（省略可）"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  rows={3}
                  maxLength={1000}
                />
                <p className="text-xs text-gray-500 mt-1">{formData.summary.length}/1000文字</p>
              </div>

              {/* タグ */}
              <div>
                <label htmlFor="tagInput" className="block text-sm font-medium text-gray-700 mb-2">
                  タグ
                </label>
                <div className="flex gap-2 mb-3">
                  <input
                    type="text"
                    id="tagInput"
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    onKeyPress={handleTagInputKeyPress}
                    placeholder="タグを入力してEnterキーで追加"
                    className="flex-1 border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    maxLength={50}
                  />
                  <button
                    type="button"
                    onClick={handleAddTag}
                    className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg"
                  >
                    追加
                  </button>
                </div>
                {formData.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2">
                    {formData.tags.map((tag, index) => (
                      <span
                        key={index}
                        className="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full flex items-center gap-2"
                      >
                        {tag}
                        <button
                          type="button"
                          onClick={() => handleRemoveTag(tag)}
                          className="text-blue-600 hover:text-blue-800 font-medium"
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                )}
              </div>

              {/* 内容 */}
              <div>
                <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
                  内容 <span className="text-red-500">*</span>
                </label>
                <textarea
                  id="content"
                  value={formData.content}
                  onChange={(e) => setFormData({ ...formData, content: e.target.value })}
                  placeholder="記事の内容を入力してください"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  rows={20}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">{formData.content.length}文字</p>
              </div>

              {/* 送信ボタン */}
              <div className="flex justify-end gap-4">
                <Link
                  href="/articles"
                  className="bg-gray-500 hover:bg-gray-600 text-white px-6 py-3 rounded-lg font-medium inline-block"
                >
                  キャンセル
                </Link>
                <button
                  type="submit"
                  disabled={loading}
                  className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? '作成中...' : '記事を作成'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}