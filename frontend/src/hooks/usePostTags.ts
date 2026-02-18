import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { setPostTags, getPostTags, searchPostsByTag, getPopularTags } from '../api/postTags';
import { useAsyncData } from './useAsyncData';
import type { Post, TagCount } from '../types/post';

export function usePostTags(postId?: number) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: tags, loading, refetch } = useAsyncData(
    async () => {
      if (!postId) return [];
      const { data } = await getPostTags(postId);
      return data?.tags || [];
    },
    { initialData: [] as string[], deps: [postId], enabled: !!postId }
  );

  const handleSetTags = useCallback(async (newTags: string[]) => {
    if (!postId) return false;
    setSaving(true);
    try {
      await setPostTags(postId, newTags);
      await refetch();
      toast.success(t('tags.updateSuccess'));
      return true;
    } catch {
      toast.error(t('tags.updateFailed'));
      return false;
    } finally {
      setSaving(false);
    }
  }, [postId, t, refetch]);

  return { tags, loading, saving, setTags: handleSetTags, refetch };
}

export function useTagSearch(tag: string, page = 1) {
  const { data, loading } = useAsyncData(
    async () => {
      if (!tag) return { posts: [] as Post[], total: 0 };
      const { data } = await searchPostsByTag(tag, page);
      return { posts: data?.data || [], total: data?.total || 0 };
    },
    { initialData: { posts: [] as Post[], total: 0 }, deps: [tag, page], enabled: !!tag }
  );

  return { posts: data.posts, total: data.total, loading };
}

export function usePopularTags() {
  const { data: tags, loading } = useAsyncData(
    async () => {
      const { data } = await getPopularTags();
      return data?.tags || [];
    },
    { initialData: [] as TagCount[], deps: [] }
  );

  return { tags, loading };
}
