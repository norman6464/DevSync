import { getRecommendedUsers, getTrendingPosts, getTrendingResources } from '../api/recommendations';
import type { RecommendedUser, TrendingPost, TrendingResource } from '../types/recommendation';
import { useAsyncData } from './useAsyncData';

/** おすすめユーザーを取得するフック */
export function useRecommendedUsers() {
  const { data, loading, error, refetch } = useAsyncData<RecommendedUser[]>(
    async () => {
      const res = await getRecommendedUsers();
      return res.data || [];
    },
    { initialData: [] as RecommendedUser[] }
  );

  return { users: data, loading, error, refetch };
}

/** トレンド投稿を取得するフック */
export function useTrendingPosts() {
  const { data, loading, error, refetch } = useAsyncData<TrendingPost[]>(
    async () => {
      const res = await getTrendingPosts();
      return res.data || [];
    },
    { initialData: [] as TrendingPost[] }
  );

  return { posts: data, loading, error, refetch };
}

/** トレンド学習リソースを取得するフック */
export function useTrendingResources() {
  const { data, loading, error, refetch } = useAsyncData<TrendingResource[]>(
    async () => {
      const res = await getTrendingResources();
      return res.data || [];
    },
    { initialData: [] as TrendingResource[] }
  );

  return { resources: data, loading, error, refetch };
}
