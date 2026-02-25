import client from './client';
import type { BookmarkCollection } from '../types/bookmarkCollection';
import type { Post } from '../types/post';

export const getMyCollections = () =>
  client.get<BookmarkCollection[]>('/bookmark-collections');

export const createCollection = (data: { name: string; description?: string; color?: string }) =>
  client.post<BookmarkCollection>('/bookmark-collections', data);

export const updateCollection = (id: number, data: { name?: string; description?: string; color?: string }) =>
  client.put<BookmarkCollection>(`/bookmark-collections/${id}`, data);

export const deleteCollection = (id: number) =>
  client.delete(`/bookmark-collections/${id}`);

export const addPostToCollection = (collectionId: number, postId: number) =>
  client.post(`/bookmark-collections/${collectionId}/posts/${postId}`);

export const removePostFromCollection = (collectionId: number, postId: number) =>
  client.delete(`/bookmark-collections/${collectionId}/posts/${postId}`);

export const getCollectionPosts = (collectionId: number, limit = 20, offset = 0) =>
  client.get<{ posts: Post[]; total: number }>(`/bookmark-collections/${collectionId}/posts`, {
    params: { limit, offset },
  });
