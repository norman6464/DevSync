import { getBookmarkedPosts } from '../api/posts';
import type { Post } from '../types/post';
import { useAsyncData } from './useAsyncData';

export function useBookmarks(page = 1, limit = 20) {
  const { data, loading, error, refetch } = useAsyncData(
    async () => {
      const { data } = await getBookmarkedPosts(page, limit);
      return data;
    },
    { initialData: { posts: [] as Post[], total: 0 }, deps: [page, limit] }
  );

  return {
    posts: data.posts,
    total: data.total,
    loading,
    error,
    refetch,
  };
}
