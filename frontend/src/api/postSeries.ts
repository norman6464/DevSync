import client from './client';
import type { PostSeries, PostSeriesItem } from '../types/post';

export const createPostSeries = (data: { title: string; description?: string }) =>
  client.post<PostSeries>('/post-series', data);

export const getPostSeries = (id: number) =>
  client.get<PostSeries>(`/post-series/${id}`);

export const getPostSeriesByUser = (userId: number) =>
  client.get<PostSeries[]>(`/post-series/user/${userId}`);

export const updatePostSeries = (id: number, data: { title?: string; description?: string }) =>
  client.put<PostSeries>(`/post-series/${id}`, data);

export const deletePostSeries = (id: number) =>
  client.delete(`/post-series/${id}`);

export const getSeriesPosts = (seriesId: number) =>
  client.get<PostSeriesItem[]>(`/post-series/${seriesId}/posts`);

export const addPostToSeries = (seriesId: number, postId: number, orderIndex: number) =>
  client.post(`/post-series/${seriesId}/posts`, { post_id: postId, order_index: orderIndex });

export const removePostFromSeries = (seriesId: number, postId: number) =>
  client.delete(`/post-series/${seriesId}/posts/${postId}`);
