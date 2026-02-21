import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useRankings } from '../useRankings';
import { getContributionRanking, getLanguageRanking, getLevelRanking, getAvailableLanguages } from '../../api/rankings';

vi.mock('../../api/rankings', () => ({
  getContributionRanking: vi.fn(),
  getLanguageRanking: vi.fn(),
  getLevelRanking: vi.fn(),
  getAvailableLanguages: vi.fn(),
}));

const mockRankings = [
  { user_id: 1, username: 'alice', name: 'Alice Tanaka', avatar_url: '', score: 100, rank: 1 },
  { user_id: 2, username: 'bob', name: 'Bob Suzuki', avatar_url: '', score: 80, rank: 2 },
  { user_id: 3, username: 'charlie', name: 'Charlie Sato', avatar_url: '', score: 60, rank: 3 },
];

describe('useRankings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getAvailableLanguages).mockResolvedValue({ data: ['Go', 'Python', 'Rust'] });
    vi.mocked(getContributionRanking).mockResolvedValue({ data: mockRankings });
    vi.mocked(getLanguageRanking).mockResolvedValue({ data: mockRankings });
    vi.mocked(getLevelRanking).mockResolvedValue({ data: mockRankings });
  });

  it('初期状態でcontributionsタブ・weekly期間でランキングが取得されること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.tab).toBe('contributions');
    expect(result.current.period).toBe('weekly');
    expect(result.current.rankings).toHaveLength(3);
    expect(getContributionRanking).toHaveBeenCalledWith('weekly');
  });

  it('名前で検索フィルタリングされること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSearchQuery('Alice');
    });

    expect(result.current.rankings).toHaveLength(1);
    expect(result.current.rankings[0].name).toBe('Alice Tanaka');
  });

  it('ユーザー名で検索フィルタリングされること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSearchQuery('bob');
    });

    expect(result.current.rankings).toHaveLength(1);
    expect(result.current.rankings[0].username).toBe('bob');
  });

  it('空の検索クエリで全件返却されること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSearchQuery('alice');
    });
    expect(result.current.rankings).toHaveLength(1);

    act(() => {
      result.current.setSearchQuery('');
    });
    expect(result.current.rankings).toHaveLength(3);
  });

  it('利用可能言語がAPIから取得されること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.languages).toEqual(['Go', 'Python', 'Rust']);
    });

    expect(getAvailableLanguages).toHaveBeenCalled();
  });

  it('利用可能言語APIが空の場合デフォルト言語が使われること', async () => {
    vi.mocked(getAvailableLanguages).mockResolvedValue({ data: [] });

    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.languages).toContain('JavaScript');
    expect(result.current.languages).toContain('TypeScript');
  });

  it('期間を変更するとsetPeriodが反映されること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPeriod('monthly');
    });

    expect(result.current.period).toBe('monthly');
  });

  it('タブを変更するとsetTabが反映されること', async () => {
    const { result } = renderHook(() => useRankings());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setTab('level');
    });

    expect(result.current.tab).toBe('level');
  });
});
