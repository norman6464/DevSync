import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import type { AxiosResponse } from 'axios';
import type { ViewCount } from '../../types/post';
import { useViewCount, useMostViewed, useRecordView } from '../usePostViews';
import { recordView, getViewCount, getMostViewed } from '../../api/postViews';
import toast from 'react-hot-toast';

vi.mock('../../api/postViews', () => ({
  recordView: vi.fn(),
  getViewCount: vi.fn(),
  getMostViewed: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn() },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

describe('useViewCount', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('postIdが指定されている場合に閲覧数を取得する', async () => {
    vi.mocked(getViewCount).mockResolvedValue({ data: { post_id: 1, view_count: 42 } } as AxiosResponse<{ post_id: number; view_count: number }>);

    const { result } = renderHook(() => useViewCount(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getViewCount).toHaveBeenCalledWith(1);
    expect(result.current.viewCount).toBe(42);
  });

  it('postIdが未指定の場合は0を返しAPIを呼ばない', async () => {
    const { result } = renderHook(() => useViewCount(undefined));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getViewCount).not.toHaveBeenCalled();
    expect(result.current.viewCount).toBe(0);
  });

  it('APIレスポンスのdataがnullの場合は0を返す', async () => {
    vi.mocked(getViewCount).mockResolvedValue({ data: null } as unknown as AxiosResponse<{ post_id: number; view_count: number }>);

    const { result } = renderHook(() => useViewCount(1));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.viewCount).toBe(0);
  });
});

describe('useMostViewed', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('人気投稿リストを取得する', async () => {
    const mockData = [
      { post_id: 1, view_count: 100 },
      { post_id: 2, view_count: 50 },
    ];
    vi.mocked(getMostViewed).mockResolvedValue({ data: mockData } as unknown as AxiosResponse<ViewCount[]>);

    const { result } = renderHook(() => useMostViewed());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getMostViewed).toHaveBeenCalled();
    expect(result.current.mostViewed).toEqual(mockData);
  });

  it('APIレスポンスのdataがnullの場合は空配列を返す', async () => {
    vi.mocked(getMostViewed).mockResolvedValue({ data: null } as unknown as AxiosResponse<ViewCount[]>);

    const { result } = renderHook(() => useMostViewed());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.mostViewed).toEqual([]);
  });
});

describe('useRecordView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('閲覧を正常に記録する', async () => {
    vi.mocked(recordView).mockResolvedValue({} as AxiosResponse);

    const { result } = renderHook(() => useRecordView());

    await act(async () => {
      await result.current.recordView(1);
    });

    expect(recordView).toHaveBeenCalledWith(1);
    expect(toast.error).not.toHaveBeenCalled();
  });

  it('APIエラー時にトーストでエラーを表示する', async () => {
    vi.mocked(recordView).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useRecordView());

    await act(async () => {
      await result.current.recordView(1);
    });

    expect(recordView).toHaveBeenCalledWith(1);
    expect(toast.error).toHaveBeenCalledWith('postViews.recordFailed');
  });
});
