import { Info, CheckCircle, AlertTriangle, XCircle, X } from 'lucide-react';

interface BannerProps {
  message: string;
  variant?: 'info' | 'success' | 'warning' | 'error';
  title?: string;
  onClose?: () => void;
  actionLabel?: string;
  onAction?: () => void;
  className?: string;
}

const variantStyles = {
  info: { bg: 'bg-blue-900/30 border-blue-800', icon: Info, iconColor: 'text-blue-400' },
  success: { bg: 'bg-green-900/30 border-green-800', icon: CheckCircle, iconColor: 'text-green-400' },
  warning: { bg: 'bg-yellow-900/30 border-yellow-800', icon: AlertTriangle, iconColor: 'text-yellow-400' },
  error: { bg: 'bg-red-900/30 border-red-800', icon: XCircle, iconColor: 'text-red-400' },
};

export default function Banner({
  message,
  variant = 'info',
  title,
  onClose,
  actionLabel,
  onAction,
  className = '',
}: BannerProps) {
  const style = variantStyles[variant];
  const Icon = style.icon;

  return (
    <div className={`flex items-start gap-3 px-4 py-3 border rounded-lg ${style.bg} ${className}`.trim()}>
      <Icon className={`w-5 h-5 flex-shrink-0 mt-0.5 ${style.iconColor}`} />
      <div className="flex-1 min-w-0">
        {title && <p className="text-sm font-semibold text-gray-200">{title}</p>}
        <p className="text-sm text-gray-300">{message}</p>
        {actionLabel && onAction && (
          <button
            type="button"
            onClick={onAction}
            className="mt-1 text-sm font-medium text-blue-400 hover:text-blue-300"
          >
            {actionLabel}
          </button>
        )}
      </div>
      {onClose && (
        <button
          type="button"
          aria-label="閉じる"
          onClick={onClose}
          className="text-gray-400 hover:text-white flex-shrink-0"
        >
          <X className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}
