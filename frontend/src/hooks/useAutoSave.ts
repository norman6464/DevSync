import { useState, useEffect, useCallback, useRef } from 'react';

export type SaveStatus = 'idle' | 'saving' | 'saved';

interface UseAutoSaveOptions<T> {
  key: string;
  data: T | null;
  delay?: number;
}

interface UseAutoSaveReturn<T> {
  saveStatus: SaveStatus;
  lastSaved: Date | null;
  clearSaved: () => void;
  getSaved: () => T | null;
}

/**
 * 自動保存フック
 * データをLocalStorageに自動的に保存します
 *
 * @param key - LocalStorageのキー
 * @param data - 保存するデータ
 * @param delay - 保存までの遅延時間（ミリ秒）デフォルト: 3000ms
 * @returns 保存状態と操作関数
 */
export function useAutoSave<T>({
  key,
  data,
  delay = 3000,
}: UseAutoSaveOptions<T>): UseAutoSaveReturn<T> {
  const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle');
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const prevDataRef = useRef<T | null>(null);

  // 保存処理
  const save = useCallback(() => {
    if (!data) return;

    try {
      setSaveStatus('saving');
      localStorage.setItem(key, JSON.stringify(data));
      setLastSaved(new Date());
      setSaveStatus('saved');
    } catch (error) {
      console.error('Failed to save to localStorage:', error);
      setSaveStatus('idle');
    }
  }, [key, data]);

  // 保存データのクリア
  const clearSaved = useCallback(() => {
    localStorage.removeItem(key);
    setLastSaved(null);
    setSaveStatus('idle');
  }, [key]);

  // 保存データの取得
  const getSaved = useCallback((): T | null => {
    try {
      const saved = localStorage.getItem(key);
      return saved ? JSON.parse(saved) : null;
    } catch (error) {
      console.error('Failed to get saved data:', error);
      return null;
    }
  }, [key]);

  // データが変更されたら自動保存
  useEffect(() => {
    // dataがnullの場合は保存しない
    if (!data) return;

    // 前回のデータと同じ場合はスキップ
    if (JSON.stringify(data) === JSON.stringify(prevDataRef.current)) {
      return;
    }

    // タイマーをクリア
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // 新しいタイマーをセット
    timeoutRef.current = setTimeout(() => {
      save();
      prevDataRef.current = data;
    }, delay);

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [data, delay, save]);

  return {
    saveStatus,
    lastSaved,
    clearSaved,
    getSaved,
  };
}
