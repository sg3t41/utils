/**
 * 日付をフォーマットする共通関数（日付のみ）
 */
export const formatDate = (dateString: string | null | undefined): string => {
  if (!dateString) return '-';
  
  try {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  } catch (error) {
    console.error('日付フォーマットエラー:', error);
    return '無効な日付';
  }
};

/**
 * 日付と時刻をフォーマットする共通関数
 */
export const formatDateTime = (dateString: string | null | undefined): string => {
  if (!dateString) return '-';
  
  try {
    return new Date(dateString).toLocaleString('ja-JP', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  } catch (error) {
    console.error('日付フォーマットエラー:', error);
    return '無効な日付';
  }
};

/**
 * 詳細な日付と時刻をフォーマットする関数（2桁表示）
 */
export const formatDetailedDateTime = (dateString: string | null | undefined): string => {
  if (!dateString) return '-';
  
  try {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (error) {
    console.error('日付フォーマットエラー:', error);
    return '無効な日付';
  }
};

/**
 * 相対的な日時を表示する関数
 */
export const formatRelativeDate = (dateString: string | null | undefined): string => {
  if (!dateString) return '-';
  
  try {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      if (diffHours === 0) {
        const diffMinutes = Math.floor(diffMs / (1000 * 60));
        return `${diffMinutes}分前`;
      }
      return `${diffHours}時間前`;
    } else if (diffDays < 7) {
      return `${diffDays}日前`;
    } else {
      return formatDate(dateString);
    }
  } catch (error) {
    console.error('相対的日付フォーマットエラー:', error);
    return '無効な日付';
  }
};