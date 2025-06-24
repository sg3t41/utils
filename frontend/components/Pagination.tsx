interface PaginationProps {
  currentPage: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
  onPageChange: (page: number) => void;
  totalItems?: number;
}

/**
 * 汎用的なページネーションコンポーネント
 * 前へ/次へボタンとページ情報を表示
 */
export default function Pagination({
  currentPage,
  totalPages,
  hasNext,
  hasPrev,
  onPageChange,
  totalItems
}: PaginationProps) {
  const handlePrevious = () => {
    if (hasPrev) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (hasNext) {
      onPageChange(currentPage + 1);
    }
  };

  if (totalPages <= 1) {
    return null;
  }

  return (
    <div className="flex flex-col items-center gap-4">
      {/* ページネーションコントロール */}
      <div className="flex justify-center items-center gap-4">
        <button
          onClick={handlePrevious}
          disabled={!hasPrev}
          className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
        >
          前へ
        </button>
        <span className="text-gray-600 min-w-[120px] text-center">
          {currentPage} / {totalPages} ページ
        </span>
        <button
          onClick={handleNext}
          disabled={!hasNext}
          className="px-4 py-2 border border-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
        >
          次へ
        </button>
      </div>

      {/* 統計情報 */}
      {totalItems !== undefined && (
        <div className="text-center text-sm text-gray-500">
          総件数: {totalItems}件
        </div>
      )}
    </div>
  );
}