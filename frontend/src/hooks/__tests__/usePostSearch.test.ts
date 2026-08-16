import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import type { AxiosResponse } from 'axios';
import { usePostSearch } from '../usePostSearch';
import { searchPosts } from '../../api/posts';
import type { PostSearchResponse } from '../../api/posts';

vi.mock('../../api/posts', () => ({
  searchPosts: vi.fn(),
}));

const mockResponse: PostSearchResponse = {
  posts: [
    { id: 1, title: 'Test Post', content: 'content' } as never,
    { id: 2, title: 'Another Post', content: 'content2' } as never,
  ],
  total: 2,
  limit: 20,
  offset: 0,
};

describe('usePostSearch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('初期状態が正しいこと', () => {
    const { result } = renderHook(() => usePostSearch());

    expect(result.current.query).toBe('');
    expect(result.current.results).toEqual([]);
    expect(result.current.total).toBe(0);
    expect(result.current.loading).toBe(false);
    expect(result.current.searched).toBe(false);
    expect(result.current.filters).toEqual({ sortBy: 'latest' });
  });

  it('setQueryでクエリが更新されること', () => {
    const { result } = renderHook(() => usePostSearch());

    act(() => {
      result.current.setQuery('React');
    });

    expect(result.current.query).toBe('React');
  });

  it('空文字でsetQueryするとresultsとtotalがリセットされること', async () => {
    vi.mocked(searchPosts).mockResolvedValue({ data: mockResponse } as AxiosResponse<PostSearchResponse>);

    const { result } = renderHook(() => usePostSearch());

    act(() => {
      result.current.setQuery('React');
    });

    await act(async () => {
      await result.current.handleSearch('React');
    });

    expect(result.current.results.length).toBe(2);

    act(() => {
      result.current.setQuery('');
    });

    expect(result.current.results).toEqual([]);
    expect(result.current.total).toBe(0);
    expect(result.current.searched).toBe(false);
  });

  it('handleSearch成功時にresultsとtotalが設定されること', async () => {
    vi.mocked(searchPosts).mockResolvedValue({ data: mockResponse } as AxiosResponse<PostSearchResponse>);

    const { result } = renderHook(() => usePostSearch());

    await act(async () => {
      await result.current.handleSearch('React');
    });

    expect(searchPosts).toHaveBeenCalledWith('React', 20, 0, { sortBy: 'latest' });
    expect(result.current.results).toEqual(mockResponse.posts);
    expect(result.current.total).toBe(2);
    expect(result.current.searched).toBe(true);
    expect(result.current.loading).toBe(false);
  });

  it('空文字での検索はAPIを呼ばずリセットすること', async () => {
    const { result } = renderHook(() => usePostSearch());

    await act(async () => {
      await result.current.handleSearch('  ');
    });

    expect(searchPosts).not.toHaveBeenCalled();
    expect(result.current.results).toEqual([]);
    expect(result.current.searched).toBe(false);
  });

  it('APIエラー時にresultsが空になりsearchedがtrueになること', async () => {
    vi.mocked(searchPosts).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => usePostSearch());

    await act(async () => {
      await result.current.handleSearch('React');
    });

    expect(result.current.results).toEqual([]);
    expect(result.current.total).toBe(0);
    expect(result.current.searched).toBe(true);
    expect(result.current.loading).toBe(false);
  });

  it('setFiltersでフィルター状態が更新されること', () => {
    const { result } = renderHook(() => usePostSearch());

    act(() => {
      result.current.setFilters({ sortBy: 'popular', tags: ['typescript'] });
    });

    expect(result.current.filters).toEqual({ sortBy: 'popular', tags: ['typescript'] });
  });

  it('nullデータの場合はデフォルト値が使われること', async () => {
    vi.mocked(searchPosts).mockResolvedValue({ data: null } as unknown as AxiosResponse<PostSearchResponse>);

    const { result } = renderHook(() => usePostSearch());

    await act(async () => {
      await result.current.handleSearch('React');
    });

    expect(result.current.results).toEqual([]);
    expect(result.current.total).toBe(0);
    expect(result.current.searched).toBe(true);
  });
});
