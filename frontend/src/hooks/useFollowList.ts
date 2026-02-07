import { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { getUser, getFollowers, getFollowing } from '../api/users';
import type { User } from '../types/user';
import { useAsyncData } from './useAsyncData';

type Tab = 'followers' | 'following';

export function useFollowList(id: string | undefined) {
  const location = useLocation();
  const userId = id ? parseInt(id) : 0;

  const initialTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
  const [tab, setTab] = useState<Tab>(initialTab);

  useEffect(() => {
    const newTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
    setTab(newTab);
  }, [location.pathname]);

  const { data: profileUser, loading: profileLoading } = useAsyncData(
    async () => {
      const { data } = await getUser(userId);
      return data as User;
    },
    { deps: [userId], enabled: !!userId }
  );

  const { data: users, loading: usersLoading } = useAsyncData(
    async () => {
      const fetcher = tab === 'followers' ? getFollowers : getFollowing;
      const { data } = await fetcher(userId);
      return (data || []) as User[];
    },
    { initialData: [] as User[], deps: [userId, tab], enabled: !!userId }
  );

  return {
    profileUser,
    users,
    tab,
    loading: usersLoading,
    profileLoading,
  };
}
