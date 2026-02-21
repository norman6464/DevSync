import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useBookmarks } from '../useBookmarks';
import { getBookmarkedPosts } from '../../api/posts';

vi.mock('../../api/posts', () => ({
  getBookmarkedPosts: vi.fn(),
}));

const mockPosts = [
  { id: 1, title: 'Go入門', content: '...', user_id: 1 },
  { id: 2, title: 'React Tips', content: '...', user_id: 2 },
  { id: 3, title: 'Docker活用', content: '...', user_id: 1 },
];

describe('useBookmarks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getBookmarkedPosts).mockResolvedValue({ data: { posts: mockPosts, total: 3 } } as never);
  });

  it('初期状態でブックマーク一覧が取得されること', async () => {
    const { result } = renderHook(() => useBookmarks());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.posts).toHaveLength(3);
    expect(result.current.total).toBe(3);
    expect(getBookmarkedPosts).toHaveBeenCalledWith(1, 20);
  });

  it('ページ番号を指定して取得できること', async () => {
    vi.mocked(getBookmarkedPosts).mockResolvedValue({ data: { posts: [mockPosts[2]], total: 3 } } as never);

    const { result } = renderHook(() => useBookmarks(2, 2));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.posts).toHaveLength(1);
    expect(getBookmarkedPosts).toHaveBeenCalledWith(2, 2);
  });

  it('ページ変更で再取得されること', async () => {
    const { result, rerender } = renderHook(
      ({ page }) => useBookmarks(page, 20),
      { initialProps: { page: 1 } }
    );

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    vi.mocked(getBookmarkedPosts).mockResolvedValue({ data: { posts: [mockPosts[2]], total: 3 } } as never);
    rerender({ page: 2 });

    await vi.waitFor(() => {
      expect(getBookmarkedPosts).toHaveBeenCalledWith(2, 20);
    });
  });

  it('API失敗時にエラーが設定されること', async () => {
    vi.mocked(getBookmarkedPosts).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useBookmarks());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBeTruthy();
    expect(result.current.posts).toEqual([]);
  });

  it('空の場合に初期値が返ること', async () => {
    vi.mocked(getBookmarkedPosts).mockResolvedValue({ data: { posts: [], total: 0 } } as never);

    const { result } = renderHook(() => useBookmarks());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.posts).toEqual([]);
    expect(result.current.total).toBe(0);
  });
});
