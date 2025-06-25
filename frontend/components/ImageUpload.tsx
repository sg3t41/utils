'use client';

import { useImageUpload } from '../hooks/useImageUpload';
import { CSS_CLASSES, ERROR_MESSAGES } from '../utils/constants';
import ImagePreview from './ui/ImagePreview';
import StatusMessage from './ui/StatusMessage';

interface ImageUploadProps {
  currentImage?: string;
  onUpload: (imagePath: string) => void;
  mounted?: boolean;
  label?: string;
}

/**
 * 画像アップロードコンポーネント
 */
export default function ImageUpload({ 
  currentImage, 
  onUpload,
  mounted = false,
  label = '記事画像（ヘッダー・サムネイルで使用）'
}: ImageUploadProps) {
  const { isUploading, uploadError, uploadImage } = useImageUpload();

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const imagePath = await uploadImage(file);
      onUpload(imagePath);
      alert(ERROR_MESSAGES.UPLOAD_SUCCESS);
    } catch (error) {
      // エラーはuseImageUploadフック内で処理済み
    }
  };

  return (
    <div>
      <label htmlFor="articleImage" className={CSS_CLASSES.LABEL}>
        {label}
      </label>
      <input
        type="file"
        id="articleImage"
        accept="image/*"
        onChange={handleFileChange}
        disabled={isUploading}
        className={CSS_CLASSES.INPUT}
      />
      {isUploading && <StatusMessage message="アップロード中..." type="info" />}
      {uploadError && <StatusMessage message={uploadError} type="error" />}
      {currentImage && mounted && (
        <ImagePreview imagePath={currentImage} alt="記事画像プレビュー" />
      )}
    </div>
  );
}