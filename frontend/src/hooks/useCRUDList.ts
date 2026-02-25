import { useState, useCallback } from 'react';
import toast from 'react-hot-toast';
import { useAsyncData } from './useAsyncData';

// ── useLocalList ────────────────────────────────────────────────
// サーバーデータの上にローカル状態を重ねる汎用コンポーザブル。
// 楽観的更新によりAPI完了を待たずに即座にUIを反映できる。

export function useLocalList<T>(serverData: T[]) {
  const [local, setLocal] = useState<T[] | null>(null);
  const items = local ?? serverData;

  const setItems = useCallback((updater: T[] | ((prev: T[]) => T[])) => {
    setLocal(prev => {
      const current = prev ?? serverData;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [serverData]);

  const resetLocal = useCallback(() => setLocal(null), []);

  return { items, setItems, resetLocal };
}

// ── useCRUDList ─────────────────────────────────────────────────
// useAsyncData + useLocalList + CRUD操作を統合した汎用フック。
// シンプルなリスト型CRUDの共通パターンを1箇所に集約する。

interface CRUDListOptions<T> {
  fetcher: () => Promise<T[]>;
  deps?: unknown[];
  enabled?: boolean;
  initialData?: T[];
}

interface MutateOptions {
  successMsg?: string;
  errorMsg?: string;
  trackSaving?: boolean;
}

export function useCRUDList<T extends { id: number }>(options: CRUDListOptions<T>) {
  const { data, loading, refetch } = useAsyncData(options.fetcher, {
    initialData: options.initialData ?? ([] as T[]),
    deps: options.deps,
    enabled: options.enabled,
  });

  const { items, setItems, resetLocal } = useLocalList(data);
  const [saving, setSaving] = useState(false);

  // リストの先頭に新アイテムを追加する
  const addItem = useCallback(async (
    apiCall: () => Promise<T>,
    opts?: MutateOptions,
  ): Promise<T | null> => {
    if (opts?.trackSaving !== false) setSaving(true);
    try {
      const newItem = await apiCall();
      setItems(prev => [newItem, ...prev]);
      if (opts?.successMsg) toast.success(opts.successMsg);
      return newItem;
    } catch {
      if (opts?.errorMsg) toast.error(opts.errorMsg);
      return null;
    } finally {
      if (opts?.trackSaving !== false) setSaving(false);
    }
  }, [setItems]);

  // IDマッチングでアイテムを差し替え更新する
  const updateItem = useCallback(async (
    apiCall: () => Promise<T>,
    opts?: MutateOptions,
  ): Promise<T | null> => {
    if (opts?.trackSaving) setSaving(true);
    try {
      const updated = await apiCall();
      setItems(prev => prev.map(i => i.id === updated.id ? updated : i));
      if (opts?.successMsg) toast.success(opts.successMsg);
      return updated;
    } catch {
      if (opts?.errorMsg) toast.error(opts.errorMsg);
      return null;
    } finally {
      if (opts?.trackSaving) setSaving(false);
    }
  }, [setItems]);

  // IDマッチングでアイテムを除去する
  const removeItem = useCallback(async (
    id: number,
    apiCall: () => Promise<unknown>,
    opts?: MutateOptions,
  ): Promise<boolean> => {
    try {
      await apiCall();
      setItems(prev => prev.filter(i => i.id !== id));
      if (opts?.successMsg) toast.success(opts.successMsg);
      return true;
    } catch {
      if (opts?.errorMsg) toast.error(opts.errorMsg);
      return false;
    }
  }, [setItems]);

  return {
    items,
    data,
    loading,
    saving,
    setItems,
    resetLocal,
    refetch,
    addItem,
    updateItem,
    removeItem,
  };
}
