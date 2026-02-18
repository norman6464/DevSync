import client from './client';
import type { PostCollection, PostCollectionItem } from '../types/post';

export const createPostCollection = (data: { title: string; description?: string; is_public?: boolean }) =>
  client.post<PostCollection>('/post-collections', data);

export const getPostCollection = (id: number) =>
  client.get<PostCollection>(`/post-collections/${id}`);

export const getPostCollectionsByUser = (userId: number) =>
  client.get<PostCollection[]>(`/post-collections/user/${userId}`);

export const updatePostCollection = (id: number, data: { title: string; description?: string; is_public?: boolean }) =>
  client.put<PostCollection>(`/post-collections/${id}`, data);

export const deletePostCollection = (id: number) =>
  client.delete(`/post-collections/${id}`);

export const getCollectionPosts = (collectionId: number) =>
  client.get<PostCollectionItem[]>(`/post-collections/${collectionId}/posts`);

export const addPostToCollection = (collectionId: number, postId: number, note?: string) =>
  client.post(`/post-collections/${collectionId}/posts`, { post_id: postId, note: note || '' });

export const removePostFromCollection = (collectionId: number, postId: number) =>
  client.delete(`/post-collections/${collectionId}/posts/${postId}`);
