import client from './client';
import type { Post, Comment } from '../types/post';

export const getPosts = (page = 1, limit = 20) =>
  client.get<Post[]>('/posts', { params: { page, limit } });

export const getTimeline = (page = 1, limit = 20) =>
  client.get<Post[]>('/posts/timeline', { params: { page, limit } });

export const getPost = (id: number) =>
  client.get<Post>(`/posts/${id}`);

export const getUserPosts = (userId: number) =>
  client.get<Post[]>(`/users/${userId}/posts`);

export const createPost = (data: { title: string; content: string; image_urls?: string }) =>
  client.post<Post>('/posts', data);

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

export const updatePost = (id: number, data: { title: string; content: string }) =>
  client.put<Post>(`/posts/${id}`, data);

export const deletePost = (id: number) =>
  client.delete(`/posts/${id}`);

export const likePost = (id: number) =>
  client.post(`/posts/${id}/likes`);

export const unlikePost = (id: number) =>
  client.delete(`/posts/${id}/likes`);

export const getComments = (postId: number) =>
  client.get<Comment[]>(`/posts/${postId}/comments`);

export const createComment = (postId: number, content: string) =>
  client.post<Comment>(`/posts/${postId}/comments`, { content });

export const deleteComment = (id: number) =>
  client.delete(`/comments/${id}`);
