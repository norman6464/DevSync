import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useCommentLike } from '../useCommentLike';
import * as postsApi from '../../api/posts';

vi.mock('../../api/posts', () => ({
  likeComment: vi.fn(),
  unlikeComment: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn() },
}));

describe('useCommentLike', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('初期状態が正しく設定される（未いいね）', () => {
    const { result } = renderHook(() => useCommentLike(1, false, 5));
    expect(result.current.liked).toBe(false);
    expect(result.current.likeCount).toBe(5);
    expect(result.current.loading).toBe(false);
  });

  it('初期状態が正しく設定される（いいね済み）', () => {
    const { result } = renderHook(() => useCommentLike(2, true, 10));
    expect(result.current.liked).toBe(true);
    expect(result.current.likeCount).toBe(10);
  });

  it('いいね成功: liked=true・likeCount+1・likeComment呼び出し', async () => {
    vi.mocked(postsApi.likeComment).mockResolvedValue({ data: undefined, status: 200, statusText: 'OK', headers: {}, config: {} as never });

    const { result } = renderHook(() => useCommentLike(1, false, 3));

    await act(async () => {
      await result.current.toggleLike();
    });

    expect(result.current.liked).toBe(true);
    expect(result.current.likeCount).toBe(4);
    expect(postsApi.likeComment).toHaveBeenCalledWith(1);
    expect(postsApi.unlikeComment).not.toHaveBeenCalled();
  });

  it('いいね解除成功: liked=false・likeCount-1・unlikeComment呼び出し', async () => {
    vi.mocked(postsApi.unlikeComment).mockResolvedValue({ data: undefined, status: 200, statusText: 'OK', headers: {}, config: {} as never });

    const { result } = renderHook(() => useCommentLike(1, true, 7));

    await act(async () => {
      await result.current.toggleLike();
    });

    expect(result.current.liked).toBe(false);
    expect(result.current.likeCount).toBe(6);
    expect(postsApi.unlikeComment).toHaveBeenCalledWith(1);
    expect(postsApi.likeComment).not.toHaveBeenCalled();
  });

  it('いいね時APIエラー → 元の状態にロールバック', async () => {
    vi.mocked(postsApi.likeComment).mockRejectedValue(new Error('network error'));

    const { result } = renderHook(() => useCommentLike(1, false, 5));

    await act(async () => {
      await result.current.toggleLike();
    });

    // ロールバック後は元の状態に戻る
    expect(result.current.liked).toBe(false);
    expect(result.current.likeCount).toBe(5);
  });

  it('いいね解除時APIエラー → 元の状態にロールバック', async () => {
    vi.mocked(postsApi.unlikeComment).mockRejectedValue(new Error('network error'));

    const { result } = renderHook(() => useCommentLike(1, true, 8));

    await act(async () => {
      await result.current.toggleLike();
    });

    // ロールバック後は元の状態に戻る
    expect(result.current.liked).toBe(true);
    expect(result.current.likeCount).toBe(8);
  });

  it('loading中にtoggleLikeを呼んでも二重実行しない', async () => {
    let resolvePromise!: () => void;
    vi.mocked(postsApi.likeComment).mockImplementation(
      () => new Promise((resolve) => { resolvePromise = () => resolve({ data: undefined, status: 200, statusText: 'OK', headers: {}, config: {} as never }); })
    );

    const { result } = renderHook(() => useCommentLike(1, false, 0));

    // 最初のtoggleLikeを開始（非同期で待たない）
    act(() => { result.current.toggleLike(); });

    // loading=trueの間に再度toggleLikeを呼ぶ
    await act(async () => {
      await result.current.toggleLike();
    });

    // resolveして完了させる
    await act(async () => {
      resolvePromise();
    });

    // likeCommentは1回だけ呼ばれる
    expect(postsApi.likeComment).toHaveBeenCalledTimes(1);
  });
});
