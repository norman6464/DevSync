import { Info, CheckCircle, AlertTriangle, XCircle, X } from 'lucide-react';

const variantConfig = {
  info: {
    icon: Info,
    borderColor: 'border-blue-500',
    bgColor: 'bg-blue-500/10',
    textColor: 'text-blue-400',
  },
  success: {
    icon: CheckCircle,
    borderColor: 'border-green-500',
    bgColor: 'bg-green-500/10',
    textColor: 'text-green-400',
  },
  warning: {
    icon: AlertTriangle,
    borderColor: 'border-yellow-500',
    bgColor: 'bg-yellow-500/10',
    textColor: 'text-yellow-400',
  },
  error: {
    icon: XCircle,
    borderColor: 'border-red-500',
    bgColor: 'bg-red-500/10',
    textColor: 'text-red-400',
  },
};

interface AlertProps {
  variant: 'info' | 'success' | 'warning' | 'error';
  message: string;
  title?: string;
  onClose?: () => void;
  className?: string;
}

export default function Alert({
  variant,
  message,
  title,
  onClose,
  className = '',
}: AlertProps) {
  const config = variantConfig[variant];
  const IconComponent = config.icon;

  return (
    <div
      role="alert"
      className={`flex items-start gap-3 p-4 border-l-4 rounded-r-lg ${config.borderColor} ${config.bgColor} ${className}`.trim()}
    >
      <IconComponent className={`w-5 h-5 ${config.textColor} mt-0.5 flex-shrink-0`} />
      <div className="flex-1 min-w-0">
        {title && (
          <h4 className={`font-medium ${config.textColor} mb-1`}>{title}</h4>
        )}
        <p className="text-sm text-gray-300">{message}</p>
      </div>
      {onClose && (
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-200 flex-shrink-0"
        >
          <X className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}
