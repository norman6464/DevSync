import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useNotifications } from '../useNotifications';
import {
  getNotifications, getUnreadCount, markAsRead, markAllAsRead,
  deleteNotification as deleteNotificationApi,
} from '../../api/notifications';

vi.mock('../../api/notifications', () => ({
  getNotifications: vi.fn(),
  getUnreadCount: vi.fn(),
  markAsRead: vi.fn(),
  markAllAsRead: vi.fn(),
  deleteNotification: vi.fn(),
}));

const mockNotifications = [
  { id: 1, type: 'like', read: false, created_at: '2026-02-19' },
  { id: 2, type: 'comment', read: true, created_at: '2026-02-18' },
];

describe('useNotifications', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    vi.mocked(getUnreadCount).mockResolvedValue({ data: { count: 3 } });
    vi.mocked(getNotifications).mockResolvedValue({
      data: { notifications: mockNotifications, total: 2 },
    });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('初期状態でfetchNotificationsとgetUnreadCountが呼ばれること', async () => {
    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getUnreadCount).toHaveBeenCalled();
    expect(getNotifications).toHaveBeenCalledWith(1, 20, undefined);
    expect(result.current.notifications).toEqual(mockNotifications);
    expect(result.current.total).toBe(2);
    expect(result.current.unreadCount).toBe(3);
  });

  it('markAsReadで既読化とunreadCount減少', async () => {
    vi.mocked(markAsRead).mockResolvedValue({});

    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.markAsRead(1);
    });

    expect(markAsRead).toHaveBeenCalledWith(1);
    expect(result.current.notifications[0].read).toBe(true);
    expect(result.current.unreadCount).toBe(2);
  });

  it('markAllAsReadで全既読化とunreadCount=0', async () => {
    vi.mocked(markAllAsRead).mockResolvedValue({});

    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.markAllAsRead();
    });

    expect(markAllAsRead).toHaveBeenCalled();
    expect(result.current.notifications.every(n => n.read)).toBe(true);
    expect(result.current.unreadCount).toBe(0);
  });

  it('未読通知の削除でunreadCountも減少すること', async () => {
    vi.mocked(deleteNotificationApi).mockResolvedValue({});

    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // id=1は未読（read: false）
    await act(async () => {
      await result.current.deleteNotification(1);
    });

    expect(deleteNotificationApi).toHaveBeenCalledWith(1);
    expect(result.current.notifications).toHaveLength(1);
    expect(result.current.total).toBe(1);
    expect(result.current.unreadCount).toBe(2);
  });

  it('既読通知の削除ではunreadCountは変わらないこと', async () => {
    vi.mocked(deleteNotificationApi).mockResolvedValue({});

    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // id=2は既読（read: true）
    await act(async () => {
      await result.current.deleteNotification(2);
    });

    expect(result.current.notifications).toHaveLength(1);
    expect(result.current.total).toBe(1);
    expect(result.current.unreadCount).toBe(3);
  });

  it('setFilterTypeでフィルタ変更時にpageが1にリセットされること', async () => {
    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setFilterType('like');
    });

    expect(result.current.filterType).toBe('like');
    expect(result.current.page).toBe(1);
  });

  it('nullデータの場合はデフォルト値が使われること', async () => {
    vi.mocked(getNotifications).mockResolvedValue({
      data: { notifications: null, total: null },
    });

    const { result } = renderHook(() => useNotifications());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.notifications).toEqual([]);
    expect(result.current.total).toBe(0);
  });
});
