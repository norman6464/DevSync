import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useLearningLogs, useStreak, useWeeklyDuration, useLearningLogCalendar } from '../useLearningLogs';
import { getMyLogs, createLog, deleteLog, favoriteLog, unfavoriteLog, getStreakInfo, getWeeklyDuration, getCalendarData } from '../../api/learningLogs';
import toast from 'react-hot-toast';

vi.mock('../../api/learningLogs', () => ({
  getMyLogs: vi.fn(),
  createLog: vi.fn(),
  updateLog: vi.fn(),
  deleteLog: vi.fn(),
  favoriteLog: vi.fn(),
  unfavoriteLog: vi.fn(),
  getStreakInfo: vi.fn(),
  getWeeklyDuration: vi.fn(),
  getCalendarData: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const mockLogs = [
  { id: 1, title: 'Go学習', content: '...', category: 'language', duration: 60, is_favorite: false, created_at: '2026-01-01' },
  { id: 2, title: 'React学習', content: '...', category: 'framework', duration: 90, is_favorite: true, created_at: '2026-01-02' },
  { id: 3, title: 'Docker入門', content: '...', category: 'skill', duration: 45, is_favorite: false, created_at: '2026-01-03' },
];

describe('useLearningLogs', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getMyLogs).mockResolvedValue({ data: mockLogs } as never);
    vi.stubGlobal('confirm', () => true);
  });

  it('初期状態でログ一覧が取得されること', async () => {
    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.logs).toHaveLength(3);
  });

  it('ログ作成が成功すること', async () => {
    const newLog = { id: 10, title: '新ログ', content: 'test', category: 'language', duration: 30, is_favorite: false, created_at: '2026-02-01' };
    vi.mocked(createLog).mockResolvedValue({ data: newLog } as never);

    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createLog({ title: '新ログ', content: 'test' });
    });

    expect(created).toEqual(newLog);
    expect(toast.success).toHaveBeenCalled();
    expect(result.current.logs.some(l => l.id === 10)).toBe(true);
  });

  it('ログ作成失敗時にエラートーストが表示されnullが返ること', async () => {
    vi.mocked(createLog).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createLog({ title: 'テスト', content: 'test' });
    });

    expect(created).toBeNull();
    expect(toast.error).toHaveBeenCalled();
  });

  it('ログ削除が成功すること', async () => {
    vi.mocked(deleteLog).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteLog(1);
    });

    expect(success).toBe(true);
    expect(result.current.logs.find(l => l.id === 1)).toBeUndefined();
    expect(toast.success).toHaveBeenCalled();
  });

  it('お気に入りトグルが動作すること', async () => {
    const updated = { ...mockLogs[0], is_favorite: true };
    vi.mocked(favoriteLog).mockResolvedValue({ data: updated } as never);

    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.toggleFavorite(1);
    });

    expect(favoriteLog).toHaveBeenCalledWith(1);
    expect(result.current.logs.find(l => l.id === 1)?.is_favorite).toBe(true);
  });

  it('お気に入り解除が動作すること', async () => {
    const updated = { ...mockLogs[1], is_favorite: false };
    vi.mocked(unfavoriteLog).mockResolvedValue({ data: updated } as never);

    const { result } = renderHook(() => useLearningLogs());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.toggleFavorite(2);
    });

    expect(unfavoriteLog).toHaveBeenCalledWith(2);
    expect(result.current.logs.find(l => l.id === 2)?.is_favorite).toBe(false);
  });
});

describe('useStreak', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('ストリーク情報が取得されること', async () => {
    vi.mocked(getStreakInfo).mockResolvedValue({ data: { current_streak: 5, longest_streak: 10 } } as never);

    const { result } = renderHook(() => useStreak(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.streakInfo?.current_streak).toBe(5);
  });

  it('userIdがundefinedの場合nullが返ること', async () => {
    const { result } = renderHook(() => useStreak(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.streakInfo).toBeNull();
    expect(getStreakInfo).not.toHaveBeenCalled();
  });
});

describe('useWeeklyDuration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('週間学習時間が取得されること', async () => {
    vi.mocked(getWeeklyDuration).mockResolvedValue({ data: { duration: 300 } } as never);

    const { result } = renderHook(() => useWeeklyDuration(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.weeklyDuration).toBe(300);
  });

  it('userIdがundefinedの場合0が返ること', async () => {
    const { result } = renderHook(() => useWeeklyDuration(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.weeklyDuration).toBe(0);
    expect(getWeeklyDuration).not.toHaveBeenCalled();
  });
});

describe('useLearningLogCalendar', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('カレンダーデータが取得されること', async () => {
    const calendarEntries = [
      { date: '2026-01-01', count: 2 },
      { date: '2026-01-02', count: 1 },
    ];
    vi.mocked(getCalendarData).mockResolvedValue({ data: calendarEntries } as never);

    const { result } = renderHook(() => useLearningLogCalendar(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.calendarData).toHaveLength(2);
  });

  it('userIdがundefinedの場合空配列が返ること', async () => {
    const { result } = renderHook(() => useLearningLogCalendar(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.calendarData).toEqual([]);
    expect(getCalendarData).not.toHaveBeenCalled();
  });
});
