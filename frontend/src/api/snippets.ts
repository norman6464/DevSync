import client from './client';
import type { CodeSnippet, SnippetComment } from '../types/post';

export const getSnippetsByPostId = (postId: number) =>
  client.get<CodeSnippet[]>(`/posts/${postId}/snippets`);

export const createSnippet = (postId: number, data: { language: string; file_name?: string; code: string }) =>
  client.post<CodeSnippet>(`/posts/${postId}/snippets`, data);

export const updateSnippet = (id: number, data: { language?: string; file_name?: string; code?: string }) =>
  client.put<CodeSnippet>(`/snippets/${id}`, data);

export const deleteSnippet = (id: number) =>
  client.delete(`/snippets/${id}`);

export const getSnippetComments = (snippetId: number) =>
  client.get<SnippetComment[]>(`/snippets/${snippetId}/comments`);

export const createSnippetComment = (snippetId: number, data: { line_number: number; content: string }) =>
  client.post<SnippetComment>(`/snippets/${snippetId}/comments`, data);

export const deleteSnippetComment = (snippetId: number, commentId: number) =>
  client.delete(`/snippets/${snippetId}/comments/${commentId}`);

export const forkSnippet = (snippetId: number, targetPostId: number) =>
  client.post<CodeSnippet>(`/snippets/${snippetId}/fork`, { target_post_id: targetPostId });

// Favorites
export const favoriteSnippet = (id: number) =>
  client.post(`/snippets/${id}/favorite`);

export const unfavoriteSnippet = (id: number) =>
  client.delete(`/snippets/${id}/favorite`);

export const getFavoritedSnippets = (limit = 20, offset = 0) =>
  client.get<{ snippets: CodeSnippet[]; total: number }>('/snippets/favorites', { params: { limit, offset } });

export const getSnippetsByLanguage = (language: string, limit = 20, offset = 0) =>
  client.get<CodeSnippet[]>(`/snippets/language/${language}`, { params: { limit, offset } });
