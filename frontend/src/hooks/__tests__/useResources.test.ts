import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useResources } from '../useResources';
import { getPublicResources, createResource, deleteResource, likeResource, saveResource } from '../../api/resources';
import type { LearningResource } from '../../types/resource';
import toast from 'react-hot-toast';

vi.mock('../../api/resources', () => ({
  getPublicResources: vi.fn(),
  getSavedResources: vi.fn(),
  searchResources: vi.fn(),
  createResource: vi.fn(),
  updateResource: vi.fn(),
  deleteResource: vi.fn(),
  likeResource: vi.fn(),
  unlikeResource: vi.fn(),
  saveResource: vi.fn(),
  unsaveResource: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

vi.mock('../../store/authStore', () => ({
  useAuthStore: (selector: (s: { user: { id: 1; name: 'Test' } }) => unknown) =>
    selector({ user: { id: 1, name: 'Test' } }),
}));

const mockResources = [
  { id: 1, title: 'Go入門', category: 'book', difficulty: 'beginner', like_count: 3, save_count: 1, is_public: true, created_at: '2026-01-01', updated_at: '2026-01-01' },
  { id: 2, title: 'React動画', category: 'video', difficulty: 'intermediate', like_count: 5, save_count: 2, is_public: true, created_at: '2026-01-02', updated_at: '2026-01-02' },
  { id: 3, title: 'Docker記事', category: 'article', difficulty: 'advanced', like_count: 1, save_count: 0, is_public: true, created_at: '2026-01-03', updated_at: '2026-01-03' },
] as LearningResource[];

describe('useResources', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getPublicResources).mockResolvedValue({ resources: mockResources, total: 3 });
    vi.stubGlobal('confirm', () => true);
  });

  it('初期状態でリソース一覧が取得されること', async () => {
    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.resources).toHaveLength(3);
    expect(result.current.total).toBe(3);
    expect(result.current.tab).toBe('explore');
  });

  it('タブ変更でページが0にリセットされること', async () => {
    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPage(2);
    });
    expect(result.current.page).toBe(2);

    act(() => {
      result.current.setTab('saved');
    });
    expect(result.current.page).toBe(0);
    expect(result.current.tab).toBe('saved');
  });

  it('リソース作成が成功すること', async () => {
    const newResource = { id: 10, title: '新リソース', category: 'book', difficulty: 'beginner', like_count: 0, save_count: 0, is_public: true, created_at: '2026-02-01', updated_at: '2026-02-01' };
    vi.mocked(createResource).mockResolvedValue(newResource as never);

    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createResource({ title: '新リソース', category: 'book' });
    });

    expect(created).toEqual(newResource);
    expect(toast.success).toHaveBeenCalled();
    expect(result.current.resources.some(r => r.id === 10)).toBe(true);
  });

  it('リソース作成失敗時にエラートーストが表示されnullが返ること', async () => {
    vi.mocked(createResource).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createResource({ title: 'テスト', category: 'book' });
    });

    expect(created).toBeNull();
    expect(toast.error).toHaveBeenCalled();
  });

  it('リソース削除が成功すること', async () => {
    vi.mocked(deleteResource).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteResource(mockResources[0] as never);
    });

    expect(success).toBe(true);
    expect(result.current.resources.find(r => r.id === 1)).toBeUndefined();
    expect(toast.success).toHaveBeenCalled();
  });

  it('カテゴリフィルター変更でページが0にリセットされること', async () => {
    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPage(3);
    });
    expect(result.current.page).toBe(3);

    act(() => {
      result.current.setCategoryFilter('book');
    });
    expect(result.current.page).toBe(0);
    expect(result.current.categoryFilter).toBe('book');
  });

  it('難易度フィルター変更でページが0にリセットされること', async () => {
    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPage(2);
    });
    expect(result.current.page).toBe(2);

    act(() => {
      result.current.setDifficultyFilter('advanced');
    });
    expect(result.current.page).toBe(0);
    expect(result.current.difficultyFilter).toBe('advanced');
  });

  it('いいね失敗時にエラートーストが表示されること', async () => {
    vi.mocked(likeResource).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.likeResource(1);
    });

    expect(toast.error).toHaveBeenCalled();
  });

  it('保存失敗時にエラートーストが表示されること', async () => {
    vi.mocked(saveResource).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useResources());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.saveResource(1);
    });

    expect(toast.error).toHaveBeenCalled();
  });
});
