import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useReactions } from '../useReactions';
import * as postsApi from '../../api/posts';

vi.mock('../../api/posts', () => ({
  getReactions: vi.fn(),
  addReaction: vi.fn(),
  removeReaction: vi.fn(),
}));

const mockAxiosResponse = (data: unknown) => ({
  data,
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {} as never,
});

describe('useReactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('初期状態は空のリアクション・ユーザーリアクション', () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({ reactions: [], user_reactions: [] })
    );
    const { result } = renderHook(() => useReactions(1));
    expect(result.current.reactions).toEqual([]);
    expect(result.current.userReactions).toEqual([]);
  });

  it('マウント時にgetReactionsでデータを取得する', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({
        reactions: [{ emoji: '👍', count: 3 }],
        user_reactions: ['👍'],
      })
    );
    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => {
      expect(result.current.reactions).toEqual([{ emoji: '👍', count: 3 }]);
    });
    expect(result.current.userReactions).toEqual(['👍']);
    expect(postsApi.getReactions).toHaveBeenCalledWith(1);
  });

  it('getReactions失敗時はエラーにならず空のまま', async () => {
    vi.mocked(postsApi.getReactions).mockRejectedValue(new Error('network'));
    const { result } = renderHook(() => useReactions(1));
    // エラーが発生してもクラッシュしない
    await waitFor(() => {
      expect(postsApi.getReactions).toHaveBeenCalled();
    });
    expect(result.current.reactions).toEqual([]);
    expect(result.current.userReactions).toEqual([]);
  });

  it('新しいリアクションを追加するとcountが1で追加される', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({ reactions: [], user_reactions: [] })
    );
    vi.mocked(postsApi.addReaction).mockResolvedValue(mockAxiosResponse({}));

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => expect(postsApi.getReactions).toHaveBeenCalled());

    await act(async () => {
      await result.current.toggleReaction('🎉');
    });

    expect(result.current.reactions).toEqual([{ emoji: '🎉', count: 1 }]);
    expect(result.current.userReactions).toContain('🎉');
    expect(postsApi.addReaction).toHaveBeenCalledWith(1, '🎉');
  });

  it('既存リアクションに追加するとcountが増加する', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({
        reactions: [{ emoji: '👍', count: 2 }],
        user_reactions: [],
      })
    );
    vi.mocked(postsApi.addReaction).mockResolvedValue(mockAxiosResponse({}));

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => {
      expect(result.current.reactions).toEqual([{ emoji: '👍', count: 2 }]);
    });

    await act(async () => {
      await result.current.toggleReaction('👍');
    });

    expect(result.current.reactions).toEqual([{ emoji: '👍', count: 3 }]);
    expect(result.current.userReactions).toContain('👍');
  });

  it('リアクション削除でcountが減少する', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({
        reactions: [{ emoji: '❤️', count: 5 }],
        user_reactions: ['❤️'],
      })
    );
    vi.mocked(postsApi.removeReaction).mockResolvedValue(mockAxiosResponse({}));

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => {
      expect(result.current.userReactions).toContain('❤️');
    });

    await act(async () => {
      await result.current.toggleReaction('❤️');
    });

    expect(result.current.reactions).toEqual([{ emoji: '❤️', count: 4 }]);
    expect(result.current.userReactions).not.toContain('❤️');
    expect(postsApi.removeReaction).toHaveBeenCalledWith(1, '❤️');
  });

  it('リアクション削除でcount=0になるとリストから除去される', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({
        reactions: [{ emoji: '🔥', count: 1 }],
        user_reactions: ['🔥'],
      })
    );
    vi.mocked(postsApi.removeReaction).mockResolvedValue(mockAxiosResponse({}));

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => {
      expect(result.current.reactions).toHaveLength(1);
    });

    await act(async () => {
      await result.current.toggleReaction('🔥');
    });

    expect(result.current.reactions).toEqual([]);
    expect(result.current.userReactions).toEqual([]);
  });

  it('addReaction失敗時は状態が変更されない', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({ reactions: [], user_reactions: [] })
    );
    vi.mocked(postsApi.addReaction).mockRejectedValue(new Error('fail'));
    const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => expect(postsApi.getReactions).toHaveBeenCalled());

    await act(async () => {
      await result.current.toggleReaction('👍');
    });

    expect(result.current.reactions).toEqual([]);
    expect(result.current.userReactions).toEqual([]);
    warnSpy.mockRestore();
  });

  it('removeReaction失敗時は状態が変更されない', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({
        reactions: [{ emoji: '👀', count: 2 }],
        user_reactions: ['👀'],
      })
    );
    vi.mocked(postsApi.removeReaction).mockRejectedValue(new Error('fail'));
    const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

    const { result } = renderHook(() => useReactions(1));
    await waitFor(() => {
      expect(result.current.userReactions).toContain('👀');
    });

    await act(async () => {
      await result.current.toggleReaction('👀');
    });

    // エラー時は状態変更されない（楽観的更新ではないため）
    expect(result.current.reactions).toEqual([{ emoji: '👀', count: 2 }]);
    expect(result.current.userReactions).toContain('👀');
    warnSpy.mockRestore();
  });

  it('postIdが変更されるとリアクションを再取得する', async () => {
    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({ reactions: [{ emoji: '👍', count: 1 }], user_reactions: [] })
    );

    const { rerender } = renderHook(({ id }) => useReactions(id), {
      initialProps: { id: 1 },
    });

    await waitFor(() => expect(postsApi.getReactions).toHaveBeenCalledWith(1));

    vi.mocked(postsApi.getReactions).mockResolvedValue(
      mockAxiosResponse({ reactions: [{ emoji: '🎉', count: 3 }], user_reactions: ['🎉'] })
    );

    rerender({ id: 2 });

    await waitFor(() => expect(postsApi.getReactions).toHaveBeenCalledWith(2));
  });
});
