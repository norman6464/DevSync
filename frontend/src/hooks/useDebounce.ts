import { useState, useEffect } from 'react';

/**
 * useDebounce - 値の変更をデバウンスする
 * @param value - デバウンスする値
 * @param delay - デバウンス時間（ミリ秒）
 * @returns デバウンスされた値
 */
export function useDebounce<T>(value: T, delay: number = 300): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    // 値が変更されたらタイマーをセット
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    // クリーンアップ: タイマーをクリア
    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}
