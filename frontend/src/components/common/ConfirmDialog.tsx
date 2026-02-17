import { useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { AlertTriangle, AlertCircle, Info } from 'lucide-react';

export interface ConfirmDialogProps {
  isOpen: boolean;
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  variant?: 'danger' | 'warning' | 'info';
  confirmText?: string;
  cancelText?: string;
}

const CONFIRM_STYLES = {
  danger: 'bg-red-600 hover:bg-red-700 focus:ring-red-500 text-white',
  warning: 'bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500 text-white',
  info: 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500 text-white',
};

const ICON_STYLES = {
  danger: 'text-red-400',
  warning: 'text-yellow-400',
  info: 'text-blue-400',
};

const ICON_COMPONENTS = {
  danger: AlertTriangle,
  warning: AlertCircle,
  info: Info,
};

export default function ConfirmDialog({
  isOpen,
  title,
  message,
  onConfirm,
  onCancel,
  variant = 'danger',
  confirmText,
  cancelText,
}: ConfirmDialogProps) {
  const { t } = useTranslation();

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onCancel();
      }
    },
    [onCancel]
  );

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      return () => document.removeEventListener('keydown', handleKeyDown);
    }
  }, [isOpen, handleKeyDown]);

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      data-testid="confirm-dialog-overlay"
      onClick={onCancel}
    >
      <div className="fixed inset-0 bg-black/50 transition-opacity" />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
        aria-describedby="confirm-dialog-message"
        className="relative z-10 w-full max-w-md mx-4 bg-gray-800 rounded-xl shadow-2xl p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start gap-4">
          {(() => {
            const IconComponent = ICON_COMPONENTS[variant];
            return (
              <div className={`flex-shrink-0 mt-0.5 ${ICON_STYLES[variant]}`}>
                <IconComponent className="w-6 h-6" />
              </div>
            );
          })()}
          <div className="flex-1">
            <h3 id="confirm-dialog-title" className="text-lg font-semibold text-white">
              {title}
            </h3>
            <p id="confirm-dialog-message" className="mt-2 text-sm text-gray-300">
              {message}
            </p>
          </div>
        </div>
        <div className="mt-6 flex justify-end gap-3">
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 text-sm font-medium text-gray-300 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
          >
            {cancelText || t('common.cancel')}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 ${CONFIRM_STYLES[variant]}`}
          >
            {confirmText || t('common.confirm')}
          </button>
        </div>
      </div>
    </div>
  );
}
