import { useToastStore } from '../store/toastStore';
import type { ToastType } from '../store/toastStore';

export function useToast() {
  const addToast = useToastStore((state) => state.addToast);

  return {
    success: (message: string, duration?: number) => addToast('success', message, duration),
    error: (message: string, duration?: number) => addToast('error', message, duration),
    warning: (message: string, duration?: number) => addToast('warning', message, duration),
    info: (message: string, duration?: number) => addToast('info', message, duration),
    toast: (type: ToastType, message: string, duration?: number) => addToast(type, message, duration),
  };
}
