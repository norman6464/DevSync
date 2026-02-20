import { useCallback } from 'react';

/**
 * Ctrl+Enter / Cmd+Enter でフォーム送信するキーハンドラーを返すフック。
 * canSubmit が true の時のみ onSubmit を呼び出す。
 */
export function useSubmitShortcut(
  onSubmit: () => void,
  canSubmit: boolean,
) {
  return useCallback(
    (e: React.KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        if (canSubmit) onSubmit();
      }
    },
    [onSubmit, canSubmit],
  );
}
