import { useState, useCallback } from 'react';
import { searchCircles } from '../api/studyCircles';
import type { StudyCircle } from '../types/studyCircle';

export function useCircleSearch() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<StudyCircle[]>([]);
  const [searched, setSearched] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSearch = useCallback(async (searchQuery?: string) => {
    const q = (searchQuery ?? query).trim();
    if (!q) {
      setResults([]);
      setSearched(false);
      return;
    }
    setLoading(true);
    try {
      const { data } = await searchCircles(q);
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
