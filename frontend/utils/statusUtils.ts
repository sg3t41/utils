import { ArticleStatus } from '@/types/article';

/**
 * ステータスに応じたバッジスタイルを返す
 */
export const getStatusBadge = (status: ArticleStatus | string): string => {
  const baseClasses = 'px-2 py-1 text-xs font-medium rounded-full';
  switch (status) {
    case 'published':
      return `${baseClasses} bg-green-100 text-green-800`;
    case 'draft':
      return `${baseClasses} bg-yellow-100 text-yellow-800`;
    case 'archived':
      return `${baseClasses} bg-gray-100 text-gray-800`;
    default:
      return `${baseClasses} bg-gray-100 text-gray-800`;
  }
};

/**
 * ステータスの日本語表示名を返す
 */
export const getStatusText = (status: ArticleStatus | string): string => {
  switch (status) {
    case 'published':
      return '公開済み';
    case 'draft':
      return '下書き';
    case 'archived':
      return 'アーカイブ';
    default:
      return status || '不明';
  }
};

/**
 * ステータス選択用のオプション配列
 */
export const statusOptions = [
  { value: 'draft', label: '下書き' },
  { value: 'published', label: '公開済み' },
  { value: 'archived', label: 'アーカイブ' },
] as const;