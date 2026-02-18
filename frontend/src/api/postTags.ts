import client from './client';
import type { Post, TagCount } from '../types/post';

export const setPostTags = (postId: number, tags: string[]) =>
  client.put(`/post-tags/posts/${postId}`, { tags });

export const getPostTags = (postId: number) =>
  client.get<{ tags: string[] }>(`/post-tags/posts/${postId}`);

export const searchPostsByTag = (tag: string, page = 1, limit = 20) =>
  client.get<{ data: Post[]; total: number; page: number; limit: number }>(
    `/post-tags/search?tag=${encodeURIComponent(tag)}&page=${page}&limit=${limit}`
  );

export const getPopularTags = () =>
  client.get<{ tags: TagCount[] }>('/post-tags/popular');
