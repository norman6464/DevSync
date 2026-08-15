import type { ReactNode } from 'react';
import { Loader2 } from 'lucide-react';

interface SpinnerButtonProps {
  children: ReactNode;
  onClick: () => void;
  loading?: boolean;
  loadingText?: string;
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
}

const variantClasses = {
  primary: 'bg-blue-600 hover:bg-blue-500 text-white',
  secondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200',
  danger: 'bg-red-600 hover:bg-red-500 text-white',
};

const sizeClasses = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-sm',
  lg: 'px-6 py-3 text-base',
};

export default function SpinnerButton({
  children,
  onClick,
  loading = false,
  loadingText,
  variant = 'primary',
  size = 'md',
  disabled = false,
  className = '',
}: SpinnerButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled || loading}
      className={`inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors disabled:opacity-50 ${variantClasses[variant]} ${sizeClasses[size]} ${className}`.trim()}
    >
      {loading && (
        <Loader2 data-testid="spinner" className="w-4 h-4 animate-spin" />
      )}
      {loading && loadingText ? loadingText : children}
    </button>
  );
}
