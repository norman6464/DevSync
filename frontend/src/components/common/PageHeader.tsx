import { Plus } from 'lucide-react';
import { buttonSecondaryClass } from '../../constants/styles';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  actionLabel?: string;
  onAction?: () => void;
}

export default function PageHeader({ title, subtitle, actionLabel, onAction }: PageHeaderProps) {
  return (
    <div className="flex items-center justify-between mb-6">
      <div>
        <h1 className="text-2xl font-bold text-white">{title}</h1>
        {subtitle && <p className="text-gray-400 text-sm mt-1">{subtitle}</p>}
      </div>
      {actionLabel && onAction && (
        <button
          onClick={onAction}
          className={`${buttonSecondaryClass} flex items-center gap-2`}
        >
          <Plus className="w-5 h-5" />
          {actionLabel}
        </button>
      )}
    </div>
  );
}
