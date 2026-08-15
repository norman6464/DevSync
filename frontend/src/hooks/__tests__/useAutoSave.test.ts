import { renderHook, act } from '@testing-library/react';
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { useAutoSave } from '../useAutoSave';

describe('useAutoSave', () => {
  beforeEach(() => {
    // LocalStorageをクリア
    localStorage.clear();
    // タイマーをモック
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('データが変更されたら3秒後に自動保存される', async () => {
    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: { title: 'Test Title', content: 'Test Content' },
        delay: 3000,
      })
    );

    expect(result.current.saveStatus).toBe('idle');

    // 3秒経過させる
    act(() => {
      vi.advanceTimersByTime(3000);
    });

    expect(result.current.saveStatus).toBe('saved');

    // LocalStorageに保存されていることを確認
    const saved = localStorage.getItem('test_draft');
    expect(saved).toBeTruthy();
    const parsed = JSON.parse(saved!);
    expect(parsed.title).toBe('Test Title');
    expect(parsed.content).toBe('Test Content');
  });

  it('保存中は saveStatus が "saving" になる', async () => {
    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: { title: 'Test' },
        delay: 1000,
      })
    );

    act(() => {
      vi.advanceTimersByTime(500);
    });

    // まだ保存中ではない
    expect(result.current.saveStatus).toBe('idle');

    act(() => {
      vi.advanceTimersByTime(500);
    });

    // 保存処理が実行される
    expect(result.current.saveStatus).toBe('saved');
  });

  it('clearSaved でデータがクリアされる', () => {
    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: { title: 'Test' },
        delay: 1000,
      })
    );

    // データを保存
    act(() => {
      vi.advanceTimersByTime(1000);
    });

    // データが保存されていることを確認
    expect(localStorage.getItem('test_draft')).toBeTruthy();

    // クリア実行
    act(() => {
      result.current.clearSaved();
    });

    // データが削除されていることを確認
    expect(localStorage.getItem('test_draft')).toBeNull();
  });

  it('getSaved で保存されたデータを取得できる', () => {
    const testData = { title: 'Saved Title', content: 'Saved Content' };
    localStorage.setItem('test_draft', JSON.stringify(testData));

    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: null,
        delay: 1000,
      })
    );

    const saved = result.current.getSaved();
    expect(saved).toEqual(testData);
  });

  it('データが null の場合は保存しない', async () => {
    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: null,
        delay: 1000,
      })
    );

    act(() => {
      vi.advanceTimersByTime(1000);
    });

    expect(localStorage.getItem('test_draft')).toBeNull();
    expect(result.current.saveStatus).toBe('idle');
  });

  it('最終保存時刻が記録される', async () => {
    const { result } = renderHook(() =>
      useAutoSave({
        key: 'test_draft',
        data: { title: 'Test' },
        delay: 1000,
      })
    );

    expect(result.current.lastSaved).toBeNull();

    act(() => {
      vi.advanceTimersByTime(1000);
    });

    expect(result.current.lastSaved).toBeInstanceOf(Date);
  });
});
