import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useBookReviews } from '../useBookReviews';
import { getBookReviews, createBookReview, deleteBookReview } from '../../api/bookReviews';
import type { BookReview } from '../../types/bookReview';

vi.mock('../../api/bookReviews', () => ({
  getBookReviews: vi.fn(),
  createBookReview: vi.fn(),
  updateBookReview: vi.fn(),
  deleteBookReview: vi.fn(),
  archiveBookReview: vi.fn(),
  unarchiveBookReview: vi.fn(),
  updateBookReviewStatus: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const makeReview = (overrides: Partial<BookReview> = {}): BookReview => ({
  id: 1,
  user_id: 1,
  title: 'テスト書籍',
  author: 'テスト著者',
  isbn: '',
  rating: 3,
  review: '',
  image_url: '',
  status: 'reading',
  is_archived: false,
  created_at: '2026-02-10T00:00:00Z',
  updated_at: '2026-02-10T00:00:00Z',
  ...overrides,
} as BookReview);

const mockReviews: BookReview[] = [
  makeReview({ id: 1, rating: 5, status: 'completed', created_at: '2026-02-01T00:00:00Z' }),
  makeReview({ id: 2, rating: 2, status: 'reading', created_at: '2026-02-03T00:00:00Z' }),
  makeReview({ id: 3, rating: 4, status: 'not_started', created_at: '2026-02-02T00:00:00Z' }),
  makeReview({ id: 4, rating: 3, is_archived: true, created_at: '2026-02-04T00:00:00Z' }),
];

describe('useBookReviews', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getBookReviews).mockResolvedValue({ reviews: mockReviews, total: 4 });
  });

  it('初期状態ではloadingがtrueであること', () => {
    const { result } = renderHook(() => useBookReviews());
    expect(result.current.loading).toBe(true);
  });

  it('レビューが取得されアーカイブ済みが除外されること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.reviews).toHaveLength(3);
    expect(result.current.reviews.every(r => !r.is_archived)).toBe(true);
  });

  it('ステータスフィルターが正しく動作すること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setStatusFilter('completed');
    });

    expect(result.current.reviews).toHaveLength(1);
    expect(result.current.reviews[0].status).toBe('completed');
  });

  it('レーティングフィルターが正しく動作すること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setRatingFilter(4);
    });

    expect(result.current.reviews.every(r => r.rating >= 4)).toBe(true);
  });

  it('アーカイブ表示切り替えが正しく動作すること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setShowArchived(true);
    });

    expect(result.current.reviews).toHaveLength(1);
    expect(result.current.reviews[0].is_archived).toBe(true);
  });

  it('ソート（oldest）が正しく動作すること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSortBy('oldest');
    });

    const ids = result.current.reviews.map(r => r.id);
    expect(ids).toEqual([1, 3, 2]);
  });

  it('ソート（ratingDesc）が正しく動作すること', async () => {
    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSortBy('ratingDesc');
    });

    const ratings = result.current.reviews.map(r => r.rating);
    expect(ratings).toEqual([5, 4, 2]);
  });

  it('レビュー作成が成功すること', async () => {
    const newReview = makeReview({ id: 10, title: '新しい書籍', rating: 5 });
    vi.mocked(createBookReview).mockResolvedValue(newReview);

    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: BookReview | null = null;
    await act(async () => {
      created = await result.current.createReview({ title: '新しい書籍', rating: 5 });
    });

    expect(created).toEqual(newReview);
    expect(result.current.reviews.some(r => r.id === 10)).toBe(true);
  });

  it('レビュー削除が成功すること', async () => {
    vi.mocked(deleteBookReview).mockResolvedValue(undefined);

    const { result } = renderHook(() => useBookReviews());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const target = result.current.reviews[0];
    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteReview(target);
    });

    expect(success).toBe(true);
    expect(result.current.reviews.find(r => r.id === target.id)).toBeUndefined();
  });
});
