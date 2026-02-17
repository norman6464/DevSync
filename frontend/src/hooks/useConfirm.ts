import { useState, useCallback, useRef } from 'react';
import type { ConfirmDialogProps } from '../components/common/ConfirmDialog';

interface ConfirmOptions {
  title: string;
  message: string;
  variant?: 'danger' | 'warning' | 'info';
  confirmText?: string;
  cancelText?: string;
}

type DialogProps = Omit<ConfirmDialogProps, 'children'>;

export function useConfirm() {
  const [dialogProps, setDialogProps] = useState<DialogProps>({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: () => {},
    onCancel: () => {},
  });

  const resolveRef = useRef<((value: boolean) => void) | null>(null);

  const confirm = useCallback((options: ConfirmOptions): Promise<boolean> => {
    return new Promise<boolean>((resolve) => {
      resolveRef.current = resolve;
      setDialogProps({
        isOpen: true,
        title: options.title,
        message: options.message,
        variant: options.variant,
        confirmText: options.confirmText,
        cancelText: options.cancelText,
        onConfirm: () => {
          resolveRef.current?.(true);
          resolveRef.current = null;
          setDialogProps((prev) => ({ ...prev, isOpen: false }));
        },
        onCancel: () => {
          resolveRef.current?.(false);
          resolveRef.current = null;
          setDialogProps((prev) => ({ ...prev, isOpen: false }));
        },
      });
    });
  }, []);

  return { confirm, dialogProps };
}
