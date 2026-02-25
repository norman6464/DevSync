import { useState, useCallback } from 'react';
import {
  getMyCollections,
  createCollection,
  updateCollection,
  deleteCollection,
  addPostToCollection,
  removePostFromCollection,
  getCollectionPosts,
} from '../api/bookmarkCollections';
import type { BookmarkCollection } from '../types/bookmarkCollection';
import type { Post } from '../types/post';
import { useAsyncData } from './useAsyncData';

export function useBookmarkCollections() {
  const { data: collections, loading, error, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyCollections();
      return data;
    },
    { initialData: [] as BookmarkCollection[] }
  );

  const create = useCallback(async (name: string, description = '', color = 'blue') => {
    const { data } = await createCollection({ name, description, color });
    await refetch();
    return data;
  }, [refetch]);

  const update = useCallback(async (id: number, updates: { name?: string; description?: string; color?: string }) => {
    const { data } = await updateCollection(id, updates);
    await refetch();
    return data;
  }, [refetch]);

  const remove = useCallback(async (id: number) => {
    await deleteCollection(id);
    await refetch();
  }, [refetch]);

  return {
    collections,
    loading,
    error,
    refetch,
    create,
    update,
    remove,
  };
}

export function useCollectionPosts(collectionId: number, limit = 20) {
  const [offset, setOffset] = useState(0);

  const { data, loading, error, refetch } = useAsyncData(
    async () => {
      const { data } = await getCollectionPosts(collectionId, limit, offset);
      return data;
    },
    { initialData: { posts: [] as Post[], total: 0 }, deps: [collectionId, limit, offset] }
  );

  const addPost = useCallback(async (postId: number) => {
    await addPostToCollection(collectionId, postId);
    await refetch();
  }, [collectionId, refetch]);

  const removePost = useCallback(async (postId: number) => {
    await removePostFromCollection(collectionId, postId);
    await refetch();
  }, [collectionId, refetch]);

  return {
    posts: data.posts,
    total: data.total,
    loading,
    error,
    offset,
    setOffset,
    refetch,
    addPost,
    removePost,
  };
}
