import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useMyLevel, useLevel, useLevelBreakdown } from '../useLevel';
import { getMyLevelInfo, getLevelInfo, getXPBreakdown } from '../../api/level';

vi.mock('../../api/level', () => ({
  getMyLevelInfo: vi.fn(),
  getLevelInfo: vi.fn(),
  getXPBreakdown: vi.fn(),
}));

const mockLevelInfo = {
  level: 5,
  total_xp: 2500,
  current_level_xp: 2000,
  next_level_xp: 3000,
  progress_xp: 500,
  progress_percent: 50,
};

const mockBreakdown = {
  learning_logs: 500,
  posts: 800,
  github: 400,
  goals: 300,
  comments: 200,
  likes: 150,
  streak_bonus: 150,
  total: 2500,
};

describe('useMyLevel', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getMyLevelInfo).mockResolvedValue({ data: mockLevelInfo } as never);
  });

  it('自分のレベル情報が取得されること', async () => {
    const { result } = renderHook(() => useMyLevel());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.levelInfo).toEqual(mockLevelInfo);
    expect(getMyLevelInfo).toHaveBeenCalled();
  });

  it('API失敗時にエラーが設定されること', async () => {
    vi.mocked(getMyLevelInfo).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useMyLevel());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBeTruthy();
  });
});

describe('useLevel', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getLevelInfo).mockResolvedValue({ data: mockLevelInfo } as never);
  });

  it('指定ユーザーのレベル情報が取得されること', async () => {
    const { result } = renderHook(() => useLevel(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.levelInfo).toEqual(mockLevelInfo);
    expect(getLevelInfo).toHaveBeenCalledWith(1);
  });

  it('userIdがundefinedの場合APIを呼ばないこと', async () => {
    const { result } = renderHook(() => useLevel(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.levelInfo).toBeNull();
    expect(getLevelInfo).not.toHaveBeenCalled();
  });
});

describe('useLevelBreakdown', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getXPBreakdown).mockResolvedValue({ data: mockBreakdown } as never);
  });

  it('XP内訳が取得されること', async () => {
    const { result } = renderHook(() => useLevelBreakdown(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.breakdown).toEqual(mockBreakdown);
    expect(getXPBreakdown).toHaveBeenCalledWith(1);
  });

  it('userIdがundefinedの場合APIを呼ばないこと', async () => {
    const { result } = renderHook(() => useLevelBreakdown(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.breakdown).toBeNull();
    expect(getXPBreakdown).not.toHaveBeenCalled();
  });
});
