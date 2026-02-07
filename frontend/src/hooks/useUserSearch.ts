import { useState, useCallback } from 'react';
import { getUsers } from '../api/users';
import { useAuthStore } from '../store/authStore';
import type { User } from '../types/user';
import { useAsyncData } from './useAsyncData';

export function useUserSearch() {
  const currentUser = useAuthStore((s) => s.user);
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<User[]>([]);
  const [searched, setSearched] = useState(false);
  const [searchLoading, setSearchLoading] = useState(false);

  const { data: allUsers, loading: initialLoading } = useAsyncData(
    async () => {
      const { data } = await getUsers();
      return data || [];
    },
    { initialData: [] as User[] }
  );

  const handleSearch = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) {
      setResults([]);
      setSearched(false);
      return;
    }
    setSearchLoading(true);
    try {
      const { data } = await getUsers(query);
      setResults(data || []);
      setSearched(true);
    } catch {
      setResults([]);
    } finally {
      setSearchLoading(false);
    }
  }, [query]);

  const handleQueryChange = useCallback((value: string) => {
    setQuery(value);
    if (!value.trim()) {
      setSearched(false);
      setResults([]);
    }
  }, []);

  const displayUsers = searched ? results : allUsers;
  const filteredUsers = displayUsers.filter(u => u.id !== currentUser?.id);
  const loading = initialLoading || searchLoading;

  return {
    query,
    setQuery: handleQueryChange,
    filteredUsers,
    loading,
    searched,
    handleSearch,
  };
}
