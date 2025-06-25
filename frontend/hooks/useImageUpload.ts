import { useState } from 'react';
import { createAuthHeaders, getApiBaseUrl } from '../utils/auth';
import { API_ENDPOINTS, ERROR_MESSAGES } from '../utils/constants';

/**
 * 画像アップロード用のカスタムフック
 */
export function useImageUpload() {
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const uploadImage = async (file: File): Promise<string> => {
    setIsUploading(true);
    setUploadError(null);

    try {
      const formData = new FormData();
      formData.append('image', file);

      const response = await fetch(
        `${getApiBaseUrl()}${API_ENDPOINTS.UPLOAD_IMAGE}`,
        {
          method: 'POST',
          headers: createAuthHeaders(),
          body: formData,
        }
      );

      if (!response.ok) {
        const responseText = await response.text();
        throw new Error(`${ERROR_MESSAGES.UPLOAD_FAILED}: ${response.status} ${responseText}`);
      }

      const data = await response.json();
      return data.image_path;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.UPLOAD_FAILED;
      setUploadError(errorMessage);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  return {
    isUploading,
    uploadError,
    uploadImage,
    setUploadError,
  };
}