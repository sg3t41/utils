'use client';

import { useState } from 'react';
import { apiClient } from '../utils/apiClient';

interface ImageUploadProps {
  currentImage?: string;
  onUpload: (imagePath: string) => void;
  mounted?: boolean;
  label?: string;
}

export default function ImageUpload({ 
  currentImage, 
  onUpload,
  mounted = false,
  label = '記事画像（ヘッダー・サムネイルで使用）'
}: ImageUploadProps) {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const handleImageUpload = async (file: File) => {
    try {
      setIsUploading(true);
      setUploadError(null);
      
      const formData = new FormData();
      formData.append('image', file);
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/upload/image`, {
        method: 'POST',
        body: formData,
      });
      
      if (!response.ok) {
        throw new Error('画像のアップロードに失敗しました');
      }
      
      const data = await response.json();
      onUpload(data.image_path);
      alert('画像をアップロードしました');
    } catch (err) {
      setUploadError(err instanceof Error ? err.message : '画像のアップロードに失敗しました');
    } finally {
      setIsUploading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      handleImageUpload(file);
    }
  };

  return (
    <div>
      <label htmlFor="articleImage" className="block text-sm font-medium text-gray-700 mb-2">
        {label}
      </label>
      <input
        type="file"
        id="articleImage"
        accept="image/*"
        onChange={handleChange}
        disabled={isUploading}
        className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
      />
      {isUploading && (
        <p className="text-xs text-blue-600 mt-1">アップロード中...</p>
      )}
      {uploadError && (
        <p className="text-xs text-red-600 mt-1">{uploadError}</p>
      )}
      {currentImage && mounted && (
        <div className="mt-2">
          <img
            src={`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/uploads/${currentImage}`}
            alt="記事画像プレビュー"
            className="w-full max-w-md h-48 object-cover rounded-lg shadow-sm"
            onError={(e) => {
              e.currentTarget.style.display = 'none';
            }}
          />
        </div>
      )}
    </div>
  );
}