import { useState, useCallback } from 'react';

interface UseArticleFiltersReturn {
  statusFilter: string;
  searchQuery: string;
  activeSearchQuery: string;
  page: number;
  setStatusFilter: (status: string) => void;
  setSearchQuery: (query: string) => void;
  setPage: (page: number) => void;
  handleSearch: (e: React.FormEvent) => void;
  handleClearSearch: () => void;
  resetPage: () => void;
}

/**
 * 記事フィルターとページネーションの状態管理を行うカスタムフック
 * 検索クエリ、ステータスフィルター、ページネーションを管理
 */
export function useArticleFilters(): UseArticleFiltersReturn {
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [activeSearchQuery, setActiveSearchQuery] = useState<string>('');
  const [page, setPage] = useState(1);

  // ステータスフィルター変更時にページをリセット
  const handleStatusChange = useCallback((status: string) => {
    setStatusFilter(status);
    setPage(1);
  }, []);

  // 検索実行
  const handleSearch = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    setActiveSearchQuery(searchQuery.trim());
  }, [searchQuery]);

  // 検索クリア
  const handleClearSearch = useCallback(() => {
    setSearchQuery('');
    setActiveSearchQuery('');
    setPage(1);
  }, []);

  // ページリセット
  const resetPage = useCallback(() => {
    setPage(1);
  }, []);

  return {
    statusFilter,
    searchQuery,
    activeSearchQuery,
    page,
    setStatusFilter: handleStatusChange,
    setSearchQuery,
    setPage,
    handleSearch,
    handleClearSearch,
    resetPage
  };
}