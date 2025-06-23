'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';

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

interface UpdateArticleRequest {
  title?: string;
  content?: string;
  summary?: string;
  tags?: string[];
  article_image?: string;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function EditArticlePage() {
  const params = useParams();
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
    featured_image: '',
  });
  const [tagInput, setTagInput] = useState('');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imageUploadLoading, setImageUploadLoading] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const fetchArticle = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/v1/articles/${params.id}`);
      if (!response.ok) {
        throw new Error('記事が見つかりませんでした');
      }

      const data: Article = await response.json();
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
    if (params.id) {
      fetchArticle();
    }
  }, [params.id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title?.trim() || !formData.content?.trim()) {
      setError('タイトルと内容は必須です');
      return;
    }

    try {
      setSaving(true);
      setError(null);

      const response = await fetch(`${API_BASE_URL}/api/v1/articles/${params.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || '記事の更新に失敗しました');
      }

      const updatedArticle = await response.json();
      alert('記事を更新しました！');
      router.push(`/articles/${updatedArticle.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : '記事の更新に失敗しました');
    } finally {
      setSaving(false);
    }
  };

  const handleAddTag = () => {
    const tag = tagInput.trim();
    if (tag && formData.tags && !formData.tags.includes(tag)) {
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
      tags: formData.tags?.filter(tag => tag !== tagToRemove) || [],
    });
  };

  const handleTagInputKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  const handleImageUpload = async (file: File) => {
    try {
      setImageUploadLoading(true);
      
      const uploadFormData = new FormData();
      uploadFormData.append('image', file);
      
      const response = await fetch(`${API_BASE_URL}/api/v1/upload/image`, {
        method: 'POST',
        body: uploadFormData,
      });
      
      if (!response.ok) {
        throw new Error('画像のアップロードに失敗しました');
      }
      
      const data = await response.json();
      
      // 画像を設定
      setFormData(prev => ({
        ...prev,
        article_image: data.image_path,
      }));
      
      alert('画像をアップロードしました');
    } catch (err) {
      alert(err instanceof Error ? err.message : '画像のアップロードに失敗しました');
    } finally {
      setImageUploadLoading(false);
    }
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setImageFile(file);
      handleImageUpload(file);
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
      <div className="max-w-4xl mx-auto px-4">
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
                  キャンセル
                </Link>
              </div>
            </div>

            {/* 記事情報 */}
            {article && (
              <div className="bg-gray-50 p-4 rounded-lg">
                <div className="grid gap-2 text-sm text-gray-600">
                  <div className="flex justify-between">
                    <span>作成日時: {new Date(article.created_at).toLocaleString('ja-JP')}</span>
                    <span>更新日時: {new Date(article.updated_at).toLocaleString('ja-JP')}</span>
                  </div>
                  <div>ステータス: {article.status === 'published' ? '公開済み' : '下書き'}</div>
                </div>
              </div>
            )}

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
                  value={formData.title || ''}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder="記事のタイトルを入力してください"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  maxLength={500}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">{(formData.title || '').length}/500文字</p>
              </div>

              {/* 概要 */}
              <div>
                <label htmlFor="summary" className="block text-sm font-medium text-gray-700 mb-2">
                  概要
                </label>
                <textarea
                  id="summary"
                  value={formData.summary || ''}
                  onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                  placeholder="記事の概要を入力してください（省略可）"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  rows={3}
                  maxLength={1000}
                />
                <p className="text-xs text-gray-500 mt-1">{(formData.summary || '').length}/1000文字</p>
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
                {formData.tags && formData.tags.length > 0 && (
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

              {/* 記事画像 */}
              <div>
                <label htmlFor="articleImage" className="block text-sm font-medium text-gray-700 mb-2">
                  記事画像（ヘッダー・サムネイルで使用）
                </label>
                <input
                  type="file"
                  id="articleImage"
                  accept="image/*"
                  onChange={handleImageChange}
                  disabled={imageUploadLoading}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                {imageUploadLoading && (
                  <p className="text-xs text-blue-600 mt-1">アップロード中...</p>
                )}
                <div className="mt-2">
                  <img
                    src={mounted ? (formData.article_image 
                      ? `${API_BASE_URL}/api/v1/uploads/${formData.article_image}` 
                      : '/noimage.svg'
                    ) : '/noimage.svg'}
                    alt="記事画像プレビュー"
                    className="w-full max-w-md h-48 object-cover rounded-lg shadow-sm bg-gray-100"
                    onError={(e) => {
                      e.currentTarget.src = '/noimage.svg';
                    }}
                  />
                </div>
              </div>

              {/* 内容 */}
              <div>
                <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
                  内容 <span className="text-red-500">*</span>
                </label>
                <textarea
                  id="content"
                  value={formData.content || ''}
                  onChange={(e) => setFormData({ ...formData, content: e.target.value })}
                  placeholder="記事の内容を入力してください"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  rows={20}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">{(formData.content || '').length}文字</p>
              </div>

              {/* 送信ボタン */}
              <div className="flex justify-end gap-4">
                <Link
                  href={`/articles/${params.id}`}
                  className="bg-gray-500 hover:bg-gray-600 text-white px-6 py-3 rounded-lg font-medium inline-block"
                >
                  キャンセル
                </Link>
                <button
                  type="submit"
                  disabled={saving}
                  className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {saving ? '更新中...' : '記事を更新'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}