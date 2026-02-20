import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useReport } from '../useReport';
import {
  getMyWeeklyReport,
  getMyMonthlyReport,
  getComparison,
} from '../../api/reports';
import type { ActivityReport, ReportComparison } from '../../api/reports';

vi.mock('../../api/reports', () => ({
  getMyWeeklyReport: vi.fn(),
  getMyMonthlyReport: vi.fn(),
  getComparison: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn() },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

const mockReport: ActivityReport = {
  period: 'weekly',
  start_date: '2026-02-10',
  end_date: '2026-02-16',
  user_id: 1,
  total_contributions: 10,
  posts_created: 3,
  comments_created: 5,
  likes_received: 20,
  goals_completed: 1,
  goals_progress: 50,
  new_followers: 2,
  messages_exchanged: 8,
  daily_contributions: [],
  top_languages: [],
};

const mockComparison: ReportComparison = {
  contributions_diff: 5,
  posts_diff: 2,
  followers_diff: 1,
  goals_diff: 0,
  trend_percentage: 15,
};

describe('useReport', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('初期状態ではweeklyレポートを取得する', async () => {
    vi.mocked(getMyWeeklyReport).mockResolvedValue({ data: mockReport });
    vi.mocked(getComparison).mockResolvedValue({ data: mockComparison });

    const { result } = renderHook(() => useReport());

    expect(result.current.period).toBe('weekly');

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getMyWeeklyReport).toHaveBeenCalled();
    expect(getComparison).toHaveBeenCalledWith('weekly');
    expect(result.current.report).toEqual(mockReport);
    expect(result.current.comparison).toEqual(mockComparison);
  });

  it('periodをmonthlyに切り替えるとmonthlyレポートを取得する', async () => {
    vi.mocked(getMyWeeklyReport).mockResolvedValue({ data: mockReport });
    vi.mocked(getMyMonthlyReport).mockResolvedValue({ data: { ...mockReport, period: 'monthly' } });
    vi.mocked(getComparison).mockResolvedValue({ data: mockComparison });

    const { result } = renderHook(() => useReport());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPeriod('monthly');
    });

    await vi.waitFor(() => {
      expect(getMyMonthlyReport).toHaveBeenCalled();
    });

    expect(getComparison).toHaveBeenCalledWith('monthly');
  });

  it('データがnullの場合はデフォルト値を返す', async () => {
    vi.mocked(getMyWeeklyReport).mockResolvedValue({ data: null });
    vi.mocked(getComparison).mockResolvedValue({ data: null });

    const { result } = renderHook(() => useReport());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.report).toBeNull();
    expect(result.current.comparison).toBeNull();
  });

  it('setPeriodでperiod状態が更新される', async () => {
    vi.mocked(getMyWeeklyReport).mockResolvedValue({ data: mockReport });
    vi.mocked(getMyMonthlyReport).mockResolvedValue({ data: mockReport });
    vi.mocked(getComparison).mockResolvedValue({ data: mockComparison });

    const { result } = renderHook(() => useReport());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.period).toBe('weekly');

    act(() => {
      result.current.setPeriod('monthly');
    });

    expect(result.current.period).toBe('monthly');
  });
});
