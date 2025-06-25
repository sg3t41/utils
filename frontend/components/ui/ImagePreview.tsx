'use client';

import { getApiBaseUrl } from '../../utils/auth';
import { API_ENDPOINTS, CSS_CLASSES } from '../../utils/constants';

interface ImagePreviewProps {
  imagePath: string;
  alt: string;
  className?: string;
}

/**
 * 画像プレビューコンポーネント
 */
export default function ImagePreview({ 
  imagePath, 
  alt, 
  className = CSS_CLASSES.IMAGE_PREVIEW 
}: ImagePreviewProps) {
  const handleImageError = (e: React.SyntheticEvent<HTMLImageElement>) => {
    e.currentTarget.style.display = 'none';
  };

  return (
    <div className="mt-2">
      <img
        src={`${getApiBaseUrl()}${API_ENDPOINTS.UPLOADS}/${imagePath}`}
        alt={alt}
        className={className}
        onError={handleImageError}
      />
    </div>
  );
}