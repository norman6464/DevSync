import { useState, useCallback } from 'react';
import { searchPosts } from '../api/posts';
import type { Post } from '../types/post';

export function usePostSearch() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Post[]>([]);
  const [searched, setSearched] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSearch = useCallback(async () => {
    if (!query.trim()) {
      setResults([]);
      setSearched(false);
      return;
    }
    setLoading(true);
    try {
      const { data } = await searchPosts(query);
      setResults(data || []);
      setSearched(true);
    } catch {
      setResults([]);
      setSearched(true);
    } finally {
      setLoading(false);
    }
  }, [query]);

  const handleQueryChange = useCallback((value: string) => {
    setQuery(value);
    if (!value.trim()) {
      setSearched(false);
      setResults([]);
    }
  }, []);

  return {
    query,
    setQuery: handleQueryChange,
    results,
    loading,
    searched,
    handleSearch,
  };
}
