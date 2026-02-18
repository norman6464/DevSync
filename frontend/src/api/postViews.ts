import client from './client';
import type { ViewCount } from '../types/post';

export const recordView = (postId: number) =>
  client.post(`/post-views/posts/${postId}`);

export const getViewCount = (postId: number) =>
  client.get<{ post_id: number; view_count: number }>(`/post-views/posts/${postId}`);

export const getMostViewed = () =>
  client.get<ViewCount[]>('/post-views/popular');
