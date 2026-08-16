import { createContext, useContext } from 'react';

export interface ToastContextValue {
  showToast: (message: string, options?: { icon?: string; color?: string }) => void;
}

// Provider（コンポーネント）と分離した非コンポーネントモジュール。
// コンポーネント以外の export が Provider のファイルにあると Fast Refresh が効かなくなるため。
export const ToastContext = createContext<ToastContextValue | null>(null);

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error('useToast must be used within ToastProvider');
  return ctx;
}
