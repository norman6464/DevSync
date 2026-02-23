import { ReactNode } from 'react';
import { X } from 'lucide-react';

interface DrawerProps {
  open: boolean;
  onClose: () => void;
  children: ReactNode;
  title?: string;
  position?: 'left' | 'right' | 'bottom';
  width?: string;
  className?: string;
}

const positionClasses = {
  right: 'right-0 top-0 h-full',
  left: 'left-0 top-0 h-full',
  bottom: 'bottom-0 left-0 w-full',
};

export default function Drawer({
  open,
  onClose,
  children,
  title,
  position = 'right',
  width = '320px',
  className = '',
}: DrawerProps) {
  if (!open) return null;

  const isHorizontal = position !== 'bottom';
  const style = isHorizontal ? { width } : undefined;

  return (
    <div className="fixed inset-0 z-50">
      <div
        data-testid="drawer-overlay"
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />
      <div
        className={`absolute ${positionClasses[position]} bg-gray-900 border-gray-800 shadow-xl flex flex-col ${className}`.trim()}
        style={style}
      >
        {title && (
          <div className="flex items-center justify-between px-4 py-3 border-b border-gray-800">
            <h3 className="text-lg font-semibold text-gray-200">{title}</h3>
            <button type="button" onClick={onClose} className="text-gray-400 hover:text-white">
              <X className="w-5 h-5" />
            </button>
          </div>
        )}
        {!title && (
          <div className="flex justify-end px-4 py-2">
            <button type="button" onClick={onClose} className="text-gray-400 hover:text-white">
              <X className="w-5 h-5" />
            </button>
          </div>
        )}
        <div className="flex-1 overflow-y-auto p-4">{children}</div>
      </div>
    </div>
  );
}
