import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { recordView, getViewCount, getMostViewed } from '../api/postViews';
import { useAsyncData } from './useAsyncData';
import type { ViewCount } from '../types/post';

export function useViewCount(postId?: number) {
  const { data: viewCount, loading, refetch } = useAsyncData(
    async () => {
      if (!postId) return 0;
      const { data } = await getViewCount(postId);
      return data?.view_count ?? 0;
    },
    { initialData: 0, deps: [postId], enabled: !!postId }
  );

  return { viewCount, loading, refetch };
}

export function useMostViewed() {
  const { data: mostViewed, loading } = useAsyncData(
    async () => {
      const { data } = await getMostViewed();
      return data || [];
    },
    { initialData: [] as ViewCount[] }
  );

  return { mostViewed, loading };
}

export function useRecordView() {
  const { t } = useTranslation();

  const record = useCallback(async (postId: number) => {
    try {
      await recordView(postId);
    } catch {
      toast.error(t('postViews.recordFailed'));
    }
  }, [t]);

  return { recordView: record };
}
