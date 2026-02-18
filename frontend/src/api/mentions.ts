import client from './client';
import type { Mention } from '../types/post';

export const getMyMentions = (page = 1, limit = 20) =>
  client.get<Mention[]>('/mentions', { params: { page, limit } });

export const getPostMentions = (postId: number) =>
  client.get<Mention[]>(`/mentions/posts/${postId}`);
