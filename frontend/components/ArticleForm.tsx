'use client';

import { useState } from 'react';
import Link from 'next/link';
import { CreateArticleRequest, UpdateArticleRequest } from '../types/article';
import TagInput from './TagInput';
import ImageUpload from './ImageUpload';

interface ArticleFormProps {
  mode: 'create' | 'edit';
  initialData?: UpdateArticleRequest;
  onSubmit: (
    data: CreateArticleRequest | UpdateArticleRequest
  ) => Promise<void>;
  isSubmitting?: boolean;
  error?: string | null;
  mounted?: boolean;
}

export default function ArticleForm({
  mode,
  initialData,
  onSubmit,
  isSubmitting = false,
  error,
  mounted = false,
}: ArticleFormProps) {
  const [formData, setFormData] = useState<
    CreateArticleRequest | UpdateArticleRequest
  >(
    initialData || {
      title: '',
      content: '',
      summary: '',
      status: 'draft', // デフォルトでdraft
      tags: [],
      article_image: '',
    }
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
  };

  const handleImageUpload = (imagePath: string) => {
    setFormData({ ...formData, article_image: imagePath });
  };

  const handleTagsChange = (tags: string[]) => {
    setFormData({ ...formData, tags });
  };

  return (
    <form onSubmit={handleSubmit} className="grid gap-6">
      {/* タイトル */}
      <div>
        <label
          htmlFor="title"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
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
        <p className="text-xs text-gray-500 mt-1">
          {formData.title.length}/500文字
        </p>
      </div>

      {/* 概要 */}
      <div>
        <label
          htmlFor="summary"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          概要
        </label>
        <textarea
          id="summary"
          value={formData.summary}
          onChange={(e) =>
            setFormData({ ...formData, summary: e.target.value })
          }
          placeholder="記事の概要を入力してください（省略可）"
          className="w-fll border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          rows={3}
          maxLength={1000}
        />
        <p className="text-xs text-gray-500 mt-1">
          {formData.summary.length}/1000文字
        </p>
      </div>

      {/* タグ */}
      <TagInput tags={formData.tags} onChange={handleTagsChange} />

      {/* 記事画像 */}
      <ImageUpload
        currentImage={formData.article_image}
        onUpload={handleImageUpload}
        mounted={mounted}
      />

      {/* 内容 */}
      <div>
        <label
          htmlFor="content"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          内容 <span className="text-red-500">*</span>
        </label>
        <textarea
          id="content"
          value={formData.content}
          onChange={(e) =>
            setFormData({ ...formData, content: e.target.value })
          }
          placeholder="記事の内容を入力してください"
          className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          rows={20}
          required
        />
        <p className="text-xs text-gray-500 mt-1">
          {formData.content.length}文字
        </p>
      </div>

      {/* エラー表示 */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">{error}</p>
        </div>
      )}

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
          disabled={isSubmitting}
          className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isSubmitting
            ? mode === 'create'
              ? '作成中...'
              : '更新中...'
            : mode === 'create'
              ? '記事を作成'
              : '記事を更新'}
        </button>
      </div>
    </form>
  );
}
