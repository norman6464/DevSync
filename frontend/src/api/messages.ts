import client from './client';
import type { Message, Conversation } from '../types/message';

export const getConversations = () =>
  client.get<Conversation[]>('/messages/conversations');

export const getMessages = (userId: number, page = 1, limit = 50) =>
  client.get<Message[]>(`/messages/${userId}`, { params: { page, limit } });

export const sendMessage = (userId: number, content: string) =>
  client.post<Message>(`/messages/${userId}`, { content });
