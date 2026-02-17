import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useConfirm } from '../useConfirm';

describe('useConfirm', () => {
  it('初期状態ではダイアログが閉じている', () => {
    const { result } = renderHook(() => useConfirm());
    expect(result.current.dialogProps.isOpen).toBe(false);
  });

  it('confirmを呼ぶとダイアログが開く', () => {
    const { result } = renderHook(() => useConfirm());
    act(() => {
      result.current.confirm({
        title: 'テスト',
        message: 'メッセージ',
      });
    });
    expect(result.current.dialogProps.isOpen).toBe(true);
    expect(result.current.dialogProps.title).toBe('テスト');
    expect(result.current.dialogProps.message).toBe('メッセージ');
  });

  it('onConfirmでPromiseがtrueに解決される', async () => {
    const { result } = renderHook(() => useConfirm());
    let promise: Promise<boolean>;
    act(() => {
      promise = result.current.confirm({
        title: 'テスト',
        message: 'メッセージ',
      });
    });
    act(() => {
      result.current.dialogProps.onConfirm();
    });
    const resolved = await promise!;
    expect(resolved).toBe(true);
    expect(result.current.dialogProps.isOpen).toBe(false);
  });

  it('onCancelでPromiseがfalseに解決される', async () => {
    const { result } = renderHook(() => useConfirm());
    let promise: Promise<boolean>;
    act(() => {
      promise = result.current.confirm({
        title: 'テスト',
        message: 'メッセージ',
      });
    });
    act(() => {
      result.current.dialogProps.onCancel();
    });
    const resolved = await promise!;
    expect(resolved).toBe(false);
    expect(result.current.dialogProps.isOpen).toBe(false);
  });

  it('variantが正しく渡される', () => {
    const { result } = renderHook(() => useConfirm());
    act(() => {
      result.current.confirm({
        title: 'テスト',
        message: 'メッセージ',
        variant: 'warning',
      });
    });
    expect(result.current.dialogProps.variant).toBe('warning');
  });

  it('confirmText/cancelTextが正しく渡される', () => {
    const { result } = renderHook(() => useConfirm());
    act(() => {
      result.current.confirm({
        title: 'テスト',
        message: 'メッセージ',
        confirmText: '削除',
        cancelText: 'やめる',
      });
    });
    expect(result.current.dialogProps.confirmText).toBe('削除');
    expect(result.current.dialogProps.cancelText).toBe('やめる');
  });
});
