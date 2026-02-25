import { useCallback } from 'react';
import { getScheduledPosts, schedulePublish, cancelSchedule } from '../api/posts';
import type { Post } from '../types/post';
import { useAsyncData } from './useAsyncData';

export function useScheduledPosts() {
  const { data: posts, loading, error, refetch } = useAsyncData(
    async () => {
      const { data } = await getScheduledPosts();
      return data;
    },
    { initialData: [] as Post[] }
  );

  const schedule = useCallback(async (postId: number, scheduledAt: string) => {
    const { data } = await schedulePublish(postId, scheduledAt);
    await refetch();
    return data;
  }, [refetch]);

  const cancel = useCallback(async (postId: number) => {
    const { data } = await cancelSchedule(postId);
    await refetch();
    return data;
  }, [refetch]);

  return {
    posts,
    loading,
    error,
    refetch,
    schedule,
    cancel,
  };
}
