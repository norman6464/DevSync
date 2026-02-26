import client from './client';
import type { PostTemplate, CreatePostTemplateRequest, UpdatePostTemplateRequest, PostTemplateListResponse } from '../types/postTemplate';

export const createPostTemplate = (data: CreatePostTemplateRequest) =>
  client.post<PostTemplate>('/post-templates', data);

export const getMyPostTemplates = (limit = 20, offset = 0) =>
  client.get<PostTemplateListResponse>(`/post-templates?limit=${limit}&offset=${offset}`);

export const getPostTemplateById = (id: number) =>
  client.get<PostTemplate>(`/post-templates/${id}`);

export const updatePostTemplate = (id: number, data: UpdatePostTemplateRequest) =>
  client.put<PostTemplate>(`/post-templates/${id}`, data);

export const deletePostTemplate = (id: number) =>
  client.delete(`/post-templates/${id}`);
