import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDebounce } from '../useDebounce';

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('初期値をそのまま返す', () => {
    const { result } = renderHook(() => useDebounce('hello'));
    expect(result.current).toBe('hello');
  });

  it('デフォルト遅延（300ms）後に値が更新される', () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value),
      { initialProps: { value: 'initial' } },
    );

    rerender({ value: 'updated' });
    expect(result.current).toBe('initial');

    act(() => {
      vi.advanceTimersByTime(299);
    });
    expect(result.current).toBe('initial');

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(result.current).toBe('updated');
  });

  it('カスタム遅延を使用できる', () => {
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      { initialProps: { value: 'a', delay: 500 } },
    );

    rerender({ value: 'b', delay: 500 });

    act(() => {
      vi.advanceTimersByTime(499);
    });
    expect(result.current).toBe('a');

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(result.current).toBe('b');
  });

  it('値が連続変更された場合、最後の値のみが反映される', () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 200),
      { initialProps: { value: 'first' } },
    );

    rerender({ value: 'second' });
    act(() => {
      vi.advanceTimersByTime(100);
    });

    rerender({ value: 'third' });
    act(() => {
      vi.advanceTimersByTime(100);
    });

    // まだ最初の値のまま（thirdからの200ms未到達）
    expect(result.current).toBe('first');

    act(() => {
      vi.advanceTimersByTime(100);
    });
    // thirdからちょうど200ms
    expect(result.current).toBe('third');
  });

  it('数値型でも動作する', () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 100),
      { initialProps: { value: 0 } },
    );

    rerender({ value: 42 });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(result.current).toBe(42);
  });

  it('オブジェクト型でも動作する', () => {
    const initial = { name: 'Alice' };
    const updated = { name: 'Bob' };

    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 100),
      { initialProps: { value: initial } },
    );

    rerender({ value: updated });

    act(() => {
      vi.advanceTimersByTime(100);
    });
    expect(result.current).toEqual({ name: 'Bob' });
  });

  it('値が変更されなければタイマーは発火しない', () => {
    const { result } = renderHook(() => useDebounce('stable', 100));

    act(() => {
      vi.advanceTimersByTime(1000);
    });
    expect(result.current).toBe('stable');
  });

  it('アンマウント時にタイマーがクリアされる', () => {
    const { result, rerender, unmount } = renderHook(
      ({ value }) => useDebounce(value, 200),
      { initialProps: { value: 'before' } },
    );

    rerender({ value: 'after' });
    unmount();

    // アンマウント後にタイマーが進んでもエラーにならない
    act(() => {
      vi.advanceTimersByTime(200);
    });
    // アンマウント済みなのでresult.currentは最後のレンダリング時の値
    expect(result.current).toBe('before');
  });
});
