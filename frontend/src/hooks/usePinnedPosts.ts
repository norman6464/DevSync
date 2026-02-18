import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { pinPost, unpinPost, getPinnedPosts, reorderPins } from '../api/postPins';
import { useAsyncData } from './useAsyncData';
import type { PostPin } from '../types/post';

export function usePinnedPosts(userId?: number) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: pins, loading, refetch } = useAsyncData(
    async () => {
      if (!userId) return [];
      const { data } = await getPinnedPosts(userId);
      return data?.pins || [];
    },
    { initialData: [] as PostPin[], deps: [userId], enabled: !!userId }
  );

  const handlePin = useCallback(async (postId: number) => {
    setSaving(true);
    try {
      await pinPost(postId);
      await refetch();
      toast.success(t('pinnedPosts.pinSuccess'));
      return true;
    } catch {
      toast.error(t('pinnedPosts.pinFailed'));
      return false;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleUnpin = useCallback(async (postId: number) => {
    setSaving(true);
    try {
      await unpinPost(postId);
      await refetch();
      toast.success(t('pinnedPosts.unpinSuccess'));
      return true;
    } catch {
      toast.error(t('pinnedPosts.unpinFailed'));
      return false;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleReorder = useCallback(async (postIds: number[]) => {
    try {
      await reorderPins(postIds);
      await refetch();
      return true;
    } catch {
      toast.error(t('pinnedPosts.reorderFailed'));
      return false;
    }
  }, [t, refetch]);

  return { pins, loading, saving, pin: handlePin, unpin: handleUnpin, reorder: handleReorder, refetch };
}
