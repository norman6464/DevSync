import client from './client';
import type { Notification } from '../types/notification';

interface NotificationsResponse {
  notifications: Notification[];
  total: number;
}

export const getNotifications = (page = 1, limit = 20, type?: string) =>
  client.get<NotificationsResponse>('/notifications', {
    params: { page, limit, ...(type && { type }) },
  });

export const getUnreadCount = () =>
  client.get<{ count: number }>('/notifications/unread-count');

export const markAsRead = (id: number) =>
  client.put(`/notifications/${id}/read`);

export const markAllAsRead = () =>
  client.put('/notifications/read-all');

export const deleteNotification = (id: number) =>
  client.delete(`/notifications/${id}`);
