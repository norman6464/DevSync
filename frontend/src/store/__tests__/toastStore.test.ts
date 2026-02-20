import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { useToastStore } from '../toastStore';

describe('toastStore', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    // ストアをリセット
    useToastStore.setState({ toasts: [] });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe('addToast', () => {
    it('トーストを追加する', () => {
      useToastStore.getState().addToast('success', 'テスト成功');

      const { toasts } = useToastStore.getState();
      expect(toasts).toHaveLength(1);
      expect(toasts[0].type).toBe('success');
      expect(toasts[0].message).toBe('テスト成功');
    });

    it('一意のIDを自動生成する', () => {
      useToastStore.getState().addToast('info', 'メッセージ1');
      useToastStore.getState().addToast('info', 'メッセージ2');

      const { toasts } = useToastStore.getState();
      expect(toasts).toHaveLength(2);
      expect(toasts[0].id).not.toBe(toasts[1].id);
    });

    it('デフォルトduration（3000ms）で自動削除される', () => {
      useToastStore.getState().addToast('success', 'auto remove');

      expect(useToastStore.getState().toasts).toHaveLength(1);

      vi.advanceTimersByTime(3000);

      expect(useToastStore.getState().toasts).toHaveLength(0);
    });

    it('カスタムdurationで自動削除される', () => {
      useToastStore.getState().addToast('warning', 'custom', 5000);

      vi.advanceTimersByTime(4999);
      expect(useToastStore.getState().toasts).toHaveLength(1);

      vi.advanceTimersByTime(1);
      expect(useToastStore.getState().toasts).toHaveLength(0);
    });

    it('duration=0の場合は自動削除されない', () => {
      useToastStore.getState().addToast('error', 'persistent', 0);

      vi.advanceTimersByTime(10000);
      expect(useToastStore.getState().toasts).toHaveLength(1);
    });

    it('複数のトーストを並行管理できる', () => {
      useToastStore.getState().addToast('success', '1つ目');
      useToastStore.getState().addToast('error', '2つ目');
      useToastStore.getState().addToast('warning', '3つ目');

      expect(useToastStore.getState().toasts).toHaveLength(3);
    });

    it('各トーストタイプを正しく設定する', () => {
      const types = ['success', 'error', 'warning', 'info'] as const;
      types.forEach((type) => {
        useToastStore.getState().addToast(type, `${type} message`);
      });

      const { toasts } = useToastStore.getState();
      types.forEach((type, i) => {
        expect(toasts[i].type).toBe(type);
      });
    });
  });

  describe('removeToast', () => {
    it('指定IDのトーストを削除する', () => {
      useToastStore.getState().addToast('success', '残る');
      useToastStore.getState().addToast('error', '消える');

      const toasts = useToastStore.getState().toasts;
      const idToRemove = toasts[1].id;

      useToastStore.getState().removeToast(idToRemove);

      const remaining = useToastStore.getState().toasts;
      expect(remaining).toHaveLength(1);
      expect(remaining[0].message).toBe('残る');
    });

    it('存在しないIDを指定してもエラーにならない', () => {
      useToastStore.getState().addToast('info', 'test');
      expect(() => useToastStore.getState().removeToast('nonexistent')).not.toThrow();
      expect(useToastStore.getState().toasts).toHaveLength(1);
    });
  });

  describe('clearAll', () => {
    it('全てのトーストを削除する', () => {
      useToastStore.getState().addToast('success', '1');
      useToastStore.getState().addToast('error', '2');
      useToastStore.getState().addToast('warning', '3');

      expect(useToastStore.getState().toasts).toHaveLength(3);

      useToastStore.getState().clearAll();

      expect(useToastStore.getState().toasts).toHaveLength(0);
    });

    it('空の状態でもエラーにならない', () => {
      expect(() => useToastStore.getState().clearAll()).not.toThrow();
      expect(useToastStore.getState().toasts).toHaveLength(0);
    });
  });
});
