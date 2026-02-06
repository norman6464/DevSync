import client from './client';
import type { Notification } from '../types/notification';

export const getNotifications = (page = 1, limit = 20) =>
  client.get<Notification[]>('/notifications', { params: { page, limit } });

export const getUnreadCount = () =>
  client.get<{ count: number }>('/notifications/unread-count');

export const markAsRead = (id: number) =>
  client.put(`/notifications/${id}/read`);

export const markAllAsRead = () =>
  client.put('/notifications/read-all');

export const deleteNotification = (id: number) =>
  client.delete(`/notifications/${id}`);
