import { ArticleStatus } from '@/types/article';

/**
 * ステータスに応じたバッジスタイルを返す
 */
export const getStatusBadge = (status: ArticleStatus): string => {
  switch (status) {
    case 'published':
      return 'bg-green-100 text-green-800 px-2 py-1 rounded-full text-sm';
    case 'draft':
      return 'bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-sm';
    case 'archived':
      return 'bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-sm';
    default:
      return 'bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-sm';
  }
};

/**
 * ステータスの日本語表示名を返す
 */
export const getStatusText = (status: ArticleStatus): string => {
  switch (status) {
    case 'published':
      return '公開済み';
    case 'draft':
      return '下書き';
    case 'archived':
      return 'アーカイブ';
    default:
      return '不明';
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