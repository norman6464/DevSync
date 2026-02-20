import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useCircleSearch } from '../useCircleSearch';
import { searchCircles } from '../../api/studyCircles';
import type { StudyCircle } from '../../types/studyCircle';

vi.mock('../../api/studyCircles', () => ({
  searchCircles: vi.fn(),
}));

const mockCircles: StudyCircle[] = [
  { id: 1, name: 'React勉強会', description: 'Reactを学ぶ', owner_id: 1, max_members: 10, is_public: true, created_at: '', updated_at: '' },
  { id: 2, name: 'Go勉強会', description: 'Goを学ぶ', owner_id: 2, max_members: 5, is_public: true, created_at: '', updated_at: '' },
];

describe('useCircleSearch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('初期状態が正しい', () => {
    const { result } = renderHook(() => useCircleSearch());
    expect(result.current.query).toBe('');
    expect(result.current.results).toEqual([]);
    expect(result.current.loading).toBe(false);
    expect(result.current.searched).toBe(false);
  });

  it('setQueryでクエリが更新される', () => {
    const { result } = renderHook(() => useCircleSearch());
    act(() => {
      result.current.setQuery('React');
    });
    expect(result.current.query).toBe('React');
  });

  it('空文字をsetQueryするとresultsとsearchedがリセットされる', async () => {
    vi.mocked(searchCircles).mockResolvedValue({ data: mockCircles });
    const { result } = renderHook(() => useCircleSearch());

    await act(async () => {
      await result.current.handleSearch('React');
    });
    expect(result.current.searched).toBe(true);

    act(() => {
      result.current.setQuery('');
    });
    expect(result.current.results).toEqual([]);
    expect(result.current.searched).toBe(false);
  });

  it('handleSearchで検索結果が取得される', async () => {
    vi.mocked(searchCircles).mockResolvedValue({ data: mockCircles });
    const { result } = renderHook(() => useCircleSearch());

    await act(async () => {
      await result.current.handleSearch('勉強会');
    });

    expect(searchCircles).toHaveBeenCalledWith('勉強会');
    expect(result.current.results).toEqual(mockCircles);
    expect(result.current.searched).toBe(true);
    expect(result.current.loading).toBe(false);
  });

  it('handleSearchで引数なしの場合はquery状態を使用する', async () => {
    vi.mocked(searchCircles).mockResolvedValue({ data: mockCircles });
    const { result } = renderHook(() => useCircleSearch());

    act(() => {
      result.current.setQuery('Go');
    });

    await act(async () => {
      await result.current.handleSearch();
    });

    expect(searchCircles).toHaveBeenCalledWith('Go');
    expect(result.current.results).toEqual(mockCircles);
  });

  it('空白のみのクエリでは検索しない', async () => {
    const { result } = renderHook(() => useCircleSearch());

    await act(async () => {
      await result.current.handleSearch('   ');
    });

    expect(searchCircles).not.toHaveBeenCalled();
    expect(result.current.results).toEqual([]);
    expect(result.current.searched).toBe(false);
  });

  it('APIエラー時はresultsが空になりsearchedがtrueになる', async () => {
    vi.mocked(searchCircles).mockRejectedValue(new Error('Network error'));
    const { result } = renderHook(() => useCircleSearch());

    await act(async () => {
      await result.current.handleSearch('test');
    });

    expect(result.current.results).toEqual([]);
    expect(result.current.searched).toBe(true);
    expect(result.current.loading).toBe(false);
  });

  it('APIレスポンスのdataがnullの場合は空配列になる', async () => {
    vi.mocked(searchCircles).mockResolvedValue({ data: null });
    const { result } = renderHook(() => useCircleSearch());

    await act(async () => {
      await result.current.handleSearch('test');
    });

    expect(result.current.results).toEqual([]);
    expect(result.current.searched).toBe(true);
  });

  it('検索中はloadingがtrueになる', async () => {
    let resolvePromise: (value: { data: StudyCircle[] }) => void;
    vi.mocked(searchCircles).mockImplementation(
      () => new Promise((resolve) => { resolvePromise = resolve; })
    );

    const { result } = renderHook(() => useCircleSearch());

    let searchPromise: Promise<void>;
    act(() => {
      searchPromise = result.current.handleSearch('React');
    });

    expect(result.current.loading).toBe(true);

    await act(async () => {
      resolvePromise!({ data: mockCircles });
      await searchPromise;
    });

    expect(result.current.loading).toBe(false);
  });
});
