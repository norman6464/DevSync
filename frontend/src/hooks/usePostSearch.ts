import { useState, useCallback } from 'react';
import { searchPosts, type PostSearchFilters } from '../api/posts';
import type { Post } from '../types/post';

export interface PostSearchState {
  query: string;
  results: Post[];
  total: number;
  loading: boolean;
  searched: boolean;
  filters: PostSearchFilters;
}

export function usePostSearch() {
  const [query, setQueryState] = useState('');
  const [results, setResults] = useState<Post[]>([]);
  const [total, setTotal] = useState(0);
  const [searched, setSearched] = useState(false);
  const [loading, setLoading] = useState(false);
  const [filters, setFiltersState] = useState<PostSearchFilters>({
    sortBy: 'latest',
  });

  const handleSearch = useCallback(async (searchQuery?: string, searchFilters?: PostSearchFilters) => {
    const q = (searchQuery ?? query).trim();
    if (!q) {
      setResults([]);
      setTotal(0);
      setSearched(false);
      return;
    }
    setLoading(true);
    const activeFilters = searchFilters ?? filters;
    try {
      const { data } = await searchPosts(q, 20, 0, activeFilters);
      setResults(data?.posts || []);
      setTotal(data?.total || 0);
      setSearched(true);
    } catch {
      setResults([]);
      setTotal(0);
      setSearched(true);
    } finally {
      setLoading(false);
    }
  }, [query, filters]);

  const setQuery = useCallback((value: string) => {
    setQueryState(value);
    if (!value.trim()) {
      setSearched(false);
      setResults([]);
      setTotal(0);
    }
  }, []);

  const setFilters = useCallback((newFilters: PostSearchFilters) => {
    setFiltersState(newFilters);
  }, []);

  return {
    query,
    setQuery,
    results,
    total,
    loading,
    searched,
    filters,
    setFilters,
    handleSearch,
  };
}
