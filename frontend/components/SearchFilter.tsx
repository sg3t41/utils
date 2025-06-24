import { getStatusText } from '../utils/statusUtils';

interface SearchFilterProps {
  statusFilter: string;
  searchQuery: string;
  activeSearchQuery: string;
  onStatusChange: (value: string) => void;
  onSearchQueryChange: (value: string) => void;
  onSearchSubmit: (e: React.FormEvent) => void;
  onClearSearch: () => void;
}

/**
 * 記事の検索・フィルタリング用コンポーネント
 * ステータスフィルターと検索フォームを提供
 */
export default function SearchFilter({
  statusFilter,
  searchQuery,
  activeSearchQuery,
  onStatusChange,
  onSearchQueryChange,
  onSearchSubmit,
  onClearSearch
}: SearchFilterProps) {
  return (
    <div className="grid gap-4">
      {/* フィルターと検索 */}
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-2">
            ステータス
          </label>
          <select
            id="status"
            value={statusFilter}
            onChange={(e) => onStatusChange(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">すべて</option>
            <option value="published">公開済み</option>
            <option value="draft">下書き</option>
            <option value="archived">アーカイブ</option>
          </select>
        </div>
        <form onSubmit={onSearchSubmit}>
          <label htmlFor="search" className="block text-sm font-medium text-gray-700 mb-2">
            検索
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              id="search"
              value={searchQuery}
              onChange={(e) => onSearchQueryChange(e.target.value)}
              placeholder="タイトルや内容を検索..."
              className="flex-1 border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white p-2 rounded-lg transition-colors"
              aria-label="検索"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </button>
            {activeSearchQuery.trim() && (
              <button
                type="button"
                onClick={onClearSearch}
                className="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg"
              >
                クリア
              </button>
            )}
          </div>
        </form>
      </div>

      {/* 現在の検索条件表示 */}
      {(activeSearchQuery.trim() || statusFilter) && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <p className="text-blue-800">
            検索条件: 
            {activeSearchQuery.trim() && <span className="ml-2 font-medium">「{activeSearchQuery.trim()}」</span>}
            {statusFilter && <span className="ml-2 font-medium">ステータス: {getStatusText(statusFilter)}</span>}
          </p>
        </div>
      )}
    </div>
  );
}