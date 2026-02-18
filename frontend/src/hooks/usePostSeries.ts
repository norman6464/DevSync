import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { PostSeries, PostSeriesItem } from '../types/post';
import {
  getPostSeriesByUser,
  createPostSeries,
  updatePostSeries,
  deletePostSeries as apiDeleteSeries,
  getSeriesPosts,
  addPostToSeries,
  removePostFromSeries,
} from '../api/postSeries';
import { useAsyncData } from './useAsyncData';

export function usePostSeries(userId?: number) {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const targetUserId = userId ?? user?.id;
  const [saving, setSaving] = useState(false);

  const { data: series, loading, refetch } = useAsyncData(
    async () => {
      if (!targetUserId) return [];
      const { data } = await getPostSeriesByUser(targetUserId);
      return data || [];
    },
    { initialData: [] as PostSeries[], deps: [targetUserId], enabled: !!targetUserId }
  );

  const handleCreate = useCallback(async (title: string, description: string) => {
    setSaving(true);
    try {
      const { data } = await createPostSeries({ title, description });
      await refetch();
      toast.success(t('series.createSuccess'));
      return data;
    } catch {
      toast.error(t('series.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleUpdate = useCallback(async (id: number, title: string, description: string) => {
    setSaving(true);
    try {
      const { data } = await updatePostSeries(id, { title, description });
      await refetch();
      toast.success(t('series.updateSuccess'));
      return data;
    } catch {
      toast.error(t('series.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await apiDeleteSeries(id);
      await refetch();
      toast.success(t('series.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('series.deleteFailed'));
      return false;
    }
  }, [t, refetch]);

  return {
    series,
    loading,
    saving,
    createSeries: handleCreate,
    updateSeries: handleUpdate,
    deleteSeries: handleDelete,
    refetch,
  };
}

export function useSeriesPosts(seriesId: number) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: items, loading, refetch } = useAsyncData(
    async () => {
      if (!seriesId) return [];
      const { data } = await getSeriesPosts(seriesId);
      return data || [];
    },
    { initialData: [] as PostSeriesItem[], deps: [seriesId], enabled: !!seriesId }
  );

  const handleAddPost = useCallback(async (postId: number, orderIndex: number) => {
    setSaving(true);
    try {
      await addPostToSeries(seriesId, postId, orderIndex);
      await refetch();
      toast.success(t('series.postAdded'));
      return true;
    } catch {
      toast.error(t('series.postAddFailed'));
      return false;
    } finally {
      setSaving(false);
    }
  }, [seriesId, t, refetch]);

  const handleRemovePost = useCallback(async (postId: number) => {
    try {
      await removePostFromSeries(seriesId, postId);
      await refetch();
      toast.success(t('series.postRemoved'));
      return true;
    } catch {
      toast.error(t('series.postRemoveFailed'));
      return false;
    }
  }, [seriesId, t, refetch]);

  return {
    items,
    loading,
    saving,
    addPost: handleAddPost,
    removePost: handleRemovePost,
    refetch,
  };
}
