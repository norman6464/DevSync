import client from './client';
import type { Post, Comment, ReactionResponse } from '../types/post';

export const getPosts = (page = 1, limit = 20) =>
  client.get<Post[]>('/posts', { params: { page, limit } });

export const getTimeline = (page = 1, limit = 20) =>
  client.get<Post[]>('/posts/timeline', { params: { page, limit } });

export const getPost = (id: number) =>
  client.get<Post>(`/posts/${id}`);

export const getUserPosts = (userId: number) =>
  client.get<Post[]>(`/users/${userId}/posts`);

export const createPost = (data: {
  title: string;
  content: string;
  image_urls?: string;
  is_draft?: boolean;
  code_snippets?: { language: string; file_name?: string; code: string }[];
}) => client.post<Post>('/posts', data);

export const uploadImage = async (file: File): Promise<{ url: string }> => {
  const formData = new FormData();
  formData.append('image', file);
  const response = await client.post<{ url: string }>('/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return response.data;
};

export const uploadImages = async (files: File[]): Promise<{ urls: string[] }> => {
  const formData = new FormData();
  files.forEach((file) => formData.append('images', file));
  const response = await client.post<{ urls: string[] }>('/upload/images', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return response.data;
};

export const updatePost = (id: number, data: { title: string; content: string; image_urls?: string }) =>
  client.put<Post>(`/posts/${id}`, data);

export const deletePost = (id: number) =>
  client.delete(`/posts/${id}`);

export const likePost = (id: number) =>
  client.post(`/posts/${id}/likes`);

export const unlikePost = (id: number) =>
  client.delete(`/posts/${id}/likes`);

export const bookmarkPost = (id: number) =>
  client.post(`/posts/${id}/bookmark`);

export const unbookmarkPost = (id: number) =>
  client.delete(`/posts/${id}/bookmark`);

export const getBookmarkedPosts = (page = 1, limit = 20) =>
  client.get<{ posts: Post[]; total: number }>('/posts/bookmarks', { params: { page, limit } });

export const getComments = (postId: number) =>
  client.get<Comment[]>(`/posts/${postId}/comments`);

export const createComment = (postId: number, content: string, parentId?: number) =>
  client.post<Comment>(`/posts/${postId}/comments`, { content, parent_id: parentId || undefined });

export const getReplies = (postId: number, commentId: number) =>
  client.get<Comment[]>(`/posts/${postId}/comments/${commentId}/replies`);

export const deleteComment = (id: number) =>
  client.delete(`/comments/${id}`);

export const searchPosts = (query: string, limit = 20, offset = 0) =>
  client.get<Post[]>('/search/posts', { params: { q: query, limit, offset } });

export const getDrafts = () =>
  client.get<Post[]>('/posts/drafts');

export const publishPost = (id: number) =>
  client.put<Post>(`/posts/${id}/publish`);

export const getReactions = (postId: number) =>
  client.get<ReactionResponse>(`/posts/${postId}/reactions`);

export const addReaction = (postId: number, emoji: string) =>
  client.post(`/posts/${postId}/reactions`, { emoji });

export const removeReaction = (postId: number, emoji: string) =>
  client.delete(`/posts/${postId}/reactions`, { data: { emoji } });
