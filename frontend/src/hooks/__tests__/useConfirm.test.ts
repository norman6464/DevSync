import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useConfirm } from '../useConfirm';

describe('useConfirm', () => {
  it('初期状態でisOpenがfalseであること', () => {
    const { result } = renderHook(() => useConfirm());

    expect(result.current.dialogProps.isOpen).toBe(false);
    expect(result.current.dialogProps.title).toBe('');
    expect(result.current.dialogProps.message).toBe('');
  });

  it('confirm()呼び出しでダイアログが開きオプションが設定されること', () => {
    const { result } = renderHook(() => useConfirm());

    act(() => {
      result.current.confirm({
        title: '削除確認',
        message: '本当に削除しますか？',
        variant: 'danger',
        confirmText: '削除',
        cancelText: 'やめる',
      });
    });

    expect(result.current.dialogProps.isOpen).toBe(true);
    expect(result.current.dialogProps.title).toBe('削除確認');
    expect(result.current.dialogProps.message).toBe('本当に削除しますか？');
    expect(result.current.dialogProps.variant).toBe('danger');
    expect(result.current.dialogProps.confirmText).toBe('削除');
    expect(result.current.dialogProps.cancelText).toBe('やめる');
  });

  it('onConfirmでPromiseがtrueで解決しダイアログが閉じること', async () => {
    const { result } = renderHook(() => useConfirm());

    let resolved: boolean | undefined;

    act(() => {
      result.current.confirm({ title: 'テスト', message: 'メッセージ' }).then((v) => {
        resolved = v;
      });
    });

    expect(result.current.dialogProps.isOpen).toBe(true);

    act(() => {
      result.current.dialogProps.onConfirm();
    });

    expect(result.current.dialogProps.isOpen).toBe(false);
    // Promiseのmicrotask解決を待つ
    await vi.waitFor(() => {
      expect(resolved).toBe(true);
    });
  });

  it('onCancelでPromiseがfalseで解決しダイアログが閉じること', async () => {
    const { result } = renderHook(() => useConfirm());

    let resolved: boolean | undefined;

    act(() => {
      result.current.confirm({ title: 'テスト', message: 'メッセージ' }).then((v) => {
        resolved = v;
      });
    });

    expect(result.current.dialogProps.isOpen).toBe(true);

    act(() => {
      result.current.dialogProps.onCancel();
    });

    expect(result.current.dialogProps.isOpen).toBe(false);
    await vi.waitFor(() => {
      expect(resolved).toBe(false);
    });
  });

  it('連続してconfirm()を呼び出しても正しく動作すること', async () => {
    const { result } = renderHook(() => useConfirm());

    // 1回目: confirm
    let first: boolean | undefined;
    act(() => {
      result.current.confirm({ title: '1回目', message: 'msg1' }).then((v) => {
        first = v;
      });
    });
    act(() => {
      result.current.dialogProps.onConfirm();
    });
    await vi.waitFor(() => {
      expect(first).toBe(true);
    });

    // 2回目: cancel
    let second: boolean | undefined;
    act(() => {
      result.current.confirm({ title: '2回目', message: 'msg2' }).then((v) => {
        second = v;
      });
    });

    expect(result.current.dialogProps.title).toBe('2回目');

    act(() => {
      result.current.dialogProps.onCancel();
    });
    await vi.waitFor(() => {
      expect(second).toBe(false);
    });
  });

  it('デフォルトオプション省略時にvariant等がundefinedであること', () => {
    const { result } = renderHook(() => useConfirm());

    act(() => {
      result.current.confirm({ title: 'シンプル', message: '確認' });
    });

    expect(result.current.dialogProps.variant).toBeUndefined();
    expect(result.current.dialogProps.confirmText).toBeUndefined();
    expect(result.current.dialogProps.cancelText).toBeUndefined();
  });
});
