import { ReactNode } from 'react';
import { X } from 'lucide-react';

interface BottomSheetProps {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  title?: string;
  maxHeight?: string;
  className?: string;
}

export default function BottomSheet({
  open,
  onClose,
  children,
  title,
  maxHeight = '50vh',
  className = '',
}: BottomSheetProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50">
      <div
        data-testid="bottomsheet-overlay"
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />
      <div
        className={`absolute bottom-0 left-0 right-0 bg-gray-900 rounded-t-2xl shadow-xl flex flex-col ${className}`.trim()}
        style={{ maxHeight }}
      >
        <div className="flex justify-center pt-3 pb-1">
          <div
            data-testid="drag-handle"
            className="w-10 h-1 bg-gray-600 rounded-full"
          />
        </div>
        {title && (
          <div className="flex items-center justify-between px-4 py-2 border-b border-gray-800">
            <h3 className="text-lg font-semibold text-gray-200">{title}</h3>
            <button
              type="button"
              aria-label="閉じる"
              onClick={onClose}
              className="text-gray-400 hover:text-white"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        )}
        <div className="flex-1 overflow-y-auto p-4">{children}</div>
      </div>
    </div>
  );
}
