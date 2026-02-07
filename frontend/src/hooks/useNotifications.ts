import { useState, useEffect, useCallback } from 'react';
import {
  getNotifications, getUnreadCount, markAsRead, markAllAsRead,
  deleteNotification as deleteNotificationApi,
} from '../api/notifications';
import type { Notification, NotificationType } from '../types/notification';

export function useNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [filterType, setFilterType] = useState<NotificationType | ''>('');
  const limit = 20;

  useEffect(() => {
    const fetchUnreadCount = async () => {
      try {
        const response = await getUnreadCount();
        setUnreadCount(response.data.count);
      } catch {
        // silently fail
      }
    };
    fetchUnreadCount();
    const interval = setInterval(fetchUnreadCount, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchNotifications = useCallback(async (p?: number, type?: NotificationType | '') => {
    setLoading(true);
    try {
      const currentPage = p ?? page;
      const currentType = type ?? filterType;
      const response = await getNotifications(currentPage, limit, currentType || undefined);
      setNotifications(response.data.notifications ?? []);
      setTotal(response.data.total ?? 0);
    } catch {
      // silently fail
    } finally {
      setLoading(false);
    }
  }, [page, filterType]);

  const handleMarkAsRead = useCallback(async (id: number) => {
    try {
      await markAsRead(id);
      setNotifications(prev => prev.map(n => n.id === id ? { ...n, read: true } : n));
      setUnreadCount(prev => Math.max(0, prev - 1));
    } catch {
      // silently fail
    }
  }, []);

  const handleMarkAllAsRead = useCallback(async () => {
    try {
      await markAllAsRead();
      setNotifications(prev => prev.map(n => ({ ...n, read: true })));
      setUnreadCount(0);
    } catch {
      // silently fail
    }
  }, []);

  const handleDelete = useCallback(async (id: number) => {
    try {
      const notification = notifications.find(n => n.id === id);
      await deleteNotificationApi(id);
      setNotifications(prev => prev.filter(n => n.id !== id));
      setTotal(prev => prev - 1);
      if (notification && !notification.read) {
        setUnreadCount(prev => Math.max(0, prev - 1));
      }
    } catch {
      // silently fail
    }
  }, [notifications]);

  const handleFilterChange = useCallback((type: NotificationType | '') => {
    setFilterType(type);
    setPage(1);
  }, []);

  useEffect(() => {
    fetchNotifications(page, filterType);
  }, [page, filterType]);

  return {
    notifications,
    unreadCount,
    total,
    loading,
    page,
    setPage,
    limit,
    filterType,
    setFilterType: handleFilterChange,
    fetchNotifications,
    markAsRead: handleMarkAsRead,
    markAllAsRead: handleMarkAllAsRead,
    deleteNotification: handleDelete,
  };
}
