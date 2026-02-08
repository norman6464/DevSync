import client from './client';
import type { ChatRoom, ChatRoomMember, GroupMessage } from '../types/chat';

export const getChatRooms = () =>
  client.get<ChatRoom[]>('/chat-rooms');

export const createChatRoom = (data: { name: string; description?: string; member_ids?: number[] }) =>
  client.post<ChatRoom>('/chat-rooms', data);

export const getChatRoom = (id: number) =>
  client.get<ChatRoom>(`/chat-rooms/${id}`);

export const updateChatRoom = (id: number, data: { name?: string; description?: string }) =>
  client.put<ChatRoom>(`/chat-rooms/${id}`, data);

export const deleteChatRoom = (id: number) =>
  client.delete(`/chat-rooms/${id}`);

export const getChatRoomMembers = (id: number) =>
  client.get<ChatRoomMember[]>(`/chat-rooms/${id}/members`);

export const addChatRoomMember = (id: number, userId: number) =>
  client.post(`/chat-rooms/${id}/members`, { user_id: userId });

export const removeChatRoomMember = (id: number, userId: number) =>
  client.delete(`/chat-rooms/${id}/members/${userId}`);

export const getChatRoomMessages = (id: number, page = 1, limit = 50) =>
  client.get<GroupMessage[]>(`/chat-rooms/${id}/messages`, { params: { page, limit } });

export const sendGroupMessage = (id: number, content: string) =>
  client.post<GroupMessage>(`/chat-rooms/${id}/messages`, { content });
