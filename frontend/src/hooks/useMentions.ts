import { getMyMentions, getPostMentions } from '../api/mentions';
import { useAsyncData } from './useAsyncData';
import type { Mention } from '../types/post';

export function useMentions(page = 1, limit = 20) {
  const { data: mentions, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyMentions(page, limit);
      return data || [];
    },
    { initialData: [] as Mention[], deps: [page, limit] }
  );

  return { mentions, loading, refetch };
}

export function usePostMentions(postId?: number) {
  const { data: mentions, loading } = useAsyncData(
    async () => {
      if (!postId) return [];
      const { data } = await getPostMentions(postId);
      return data || [];
    },
    { initialData: [] as Mention[], deps: [postId], enabled: !!postId }
  );

  return { mentions, loading };
}
