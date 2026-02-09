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
