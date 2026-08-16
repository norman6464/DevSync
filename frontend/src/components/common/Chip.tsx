import type { ReactNode } from 'react';
import { X } from 'lucide-react';

const sizeClasses = {
  sm: 'text-xs px-2 py-0.5',
  md: 'text-sm px-3 py-1',
  lg: 'text-base px-4 py-1.5',
};

interface ChipProps {
  label: string;
  selected?: boolean;
  onClick?: () => void;
  onDelete?: () => void;
  icon?: ReactNode;
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
}

export default function Chip({
  label,
  selected = false,
  onClick,
  onDelete,
  icon,
  size = 'md',
  disabled = false,
  className = '',
}: ChipProps) {
  return (
    <span
      onClick={disabled ? undefined : onClick}
      className={`inline-flex items-center gap-1.5 rounded-full border transition-colors ${sizeClasses[size]} ${
        selected ? 'bg-blue-600 border-blue-500 text-white' : 'bg-gray-800 border-gray-700 text-gray-300'
      } ${onClick && !disabled ? 'cursor-pointer hover:border-gray-600' : ''} ${disabled ? 'opacity-50 cursor-not-allowed' : ''} ${className}`.trim()}
    >
      {icon}
      {label}
      {onDelete && (
        <button
          type="button"
          onClick={(e) => { e.stopPropagation(); onDelete(); }}
          className="ml-0.5 text-gray-400 hover:text-white"
        >
          <X className="w-3 h-3" />
        </button>
      )}
    </span>
  );
}
