import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { PostCollection, PostCollectionItem } from '../types/post';
import {
  getPostCollectionsByUser,
  createPostCollection,
  updatePostCollection,
  deletePostCollection as apiDeleteCollection,
  getCollectionPosts,
  addPostToCollection,
  removePostFromCollection,
} from '../api/postCollections';
import { useAsyncData } from './useAsyncData';

export function usePostCollections(userId?: number) {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const targetUserId = userId ?? user?.id;
  const [saving, setSaving] = useState(false);

  const { data: collections, loading, refetch } = useAsyncData(
    async () => {
      if (!targetUserId) return [];
      const { data } = await getPostCollectionsByUser(targetUserId);
      return data || [];
    },
    { initialData: [] as PostCollection[], deps: [targetUserId], enabled: !!targetUserId }
  );

  const handleCreate = useCallback(async (title: string, description: string, isPublic: boolean) => {
    setSaving(true);
    try {
      const { data } = await createPostCollection({ title, description, is_public: isPublic });
      await refetch();
      toast.success(t('collections.createSuccess'));
      return data;
    } catch {
      toast.error(t('collections.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleUpdate = useCallback(async (id: number, title: string, description: string, isPublic: boolean) => {
    setSaving(true);
    try {
      const { data } = await updatePostCollection(id, { title, description, is_public: isPublic });
      await refetch();
      toast.success(t('collections.updateSuccess'));
      return data;
    } catch {
      toast.error(t('collections.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleDelete = useCallback(async (id: number) => {
    try {
      await apiDeleteCollection(id);
      await refetch();
      toast.success(t('collections.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('collections.deleteFailed'));
      return false;
    }
  }, [t, refetch]);

  return {
    collections,
    loading,
    saving,
    createCollection: handleCreate,
    updateCollection: handleUpdate,
    deleteCollection: handleDelete,
    refetch,
  };
}

export function useCollectionPosts(collectionId: number) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: items, loading, refetch } = useAsyncData(
    async () => {
      if (!collectionId) return [];
      const { data } = await getCollectionPosts(collectionId);
      return data || [];
    },
    { initialData: [] as PostCollectionItem[], deps: [collectionId], enabled: !!collectionId }
  );

  const handleAddPost = useCallback(async (postId: number, note?: string) => {
    setSaving(true);
    try {
      await addPostToCollection(collectionId, postId, note);
      await refetch();
      toast.success(t('collections.postAdded'));
      return true;
    } catch {
      toast.error(t('collections.postAddFailed'));
      return false;
    } finally {
      setSaving(false);
    }
  }, [collectionId, t, refetch]);

  const handleRemovePost = useCallback(async (postId: number) => {
    try {
      await removePostFromCollection(collectionId, postId);
      await refetch();
      toast.success(t('collections.postRemoved'));
      return true;
    } catch {
      toast.error(t('collections.postRemoveFailed'));
      return false;
    }
  }, [collectionId, t, refetch]);

  return {
    items,
    loading,
    saving,
    addPost: handleAddPost,
    removePost: handleRemovePost,
    refetch,
  };
}
