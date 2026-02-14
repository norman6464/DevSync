import { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { getUser, getUserByUsername, getFollowers, getFollowing } from '../api/users';
import type { User } from '../types/user';
import { useAsyncData } from './useAsyncData';

type Tab = 'followers' | 'following';

export function useFollowList(usernameOrId: string | undefined) {
  const location = useLocation();
  // usernameOrIdが数値の場合はID、そうでない場合はusername
  const isId = usernameOrId && /^\d+$/.test(usernameOrId);
  const userId = isId ? parseInt(usernameOrId) : 0;
  const username = !isId ? usernameOrId : '';

  const initialTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
  const [tab, setTab] = useState<Tab>(initialTab);

  useEffect(() => {
    const newTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
    setTab(newTab);
  }, [location.pathname]);

  const { data: profileUser, loading: profileLoading } = useAsyncData(
    async () => {
      const userRes = username
        ? await getUserByUsername(username)
        : await getUser(userId);
      return userRes.data as User;
    },
    { deps: [usernameOrId], enabled: !!usernameOrId }
  );

  const { data: users, loading: usersLoading } = useAsyncData(
    async () => {
      const actualUserId = profileUser?.id || userId;
      if (!actualUserId) return [];
      const fetcher = tab === 'followers' ? getFollowers : getFollowing;
      const { data } = await fetcher(actualUserId);
      return (data || []) as User[];
    },
    { initialData: [] as User[], deps: [profileUser?.id, tab], enabled: !!profileUser }
  );

  return {
    profileUser,
    users,
    tab,
    loading: usersLoading,
    profileLoading,
  };
}
