// ステータス関連のユーティリティ関数

export function getStatusBadge(status: string): { className: string; text: string } {
  switch (status) {
    case 'published':
      return {
        className: 'bg-green-100 text-green-800 px-2 py-1 rounded-full text-sm',
        text: '公開済み'
      };
    case 'draft':
      return {
        className: 'bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-sm',
        text: '下書き'
      };
    default:
      return {
        className: 'bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-sm',
        text: '不明'
      };
  }
}

export function getStatusText(status: string): string {
  return getStatusBadge(status).text;
}