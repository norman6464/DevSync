import client from './client';
import type { PostPin } from '../types/post';

export const pinPost = (postId: number) =>
  client.post(`/post-pins/posts/${postId}`);

export const unpinPost = (postId: number) =>
  client.delete(`/post-pins/posts/${postId}`);

export const getPinnedPosts = (userId: number) =>
  client.get<{ pins: PostPin[] }>(`/post-pins/users/${userId}`);

export const reorderPins = (postIds: number[]) =>
  client.put('/post-pins/reorder', { post_ids: postIds });
