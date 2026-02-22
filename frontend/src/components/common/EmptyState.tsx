import { ReactNode } from 'react';
import { Inbox, FileText, Users, Search, FolderOpen } from 'lucide-react';

const iconMap = {
  inbox: Inbox,
  file: FileText,
  users: Users,
  search: Search,
  folder: FolderOpen,
};

const iconSizeMap = {
  sm: 'w-8 h-8',
  md: 'w-12 h-12',
  lg: 'w-16 h-16',
};

interface EmptyStateProps {
  title: string;
  description?: string;
  icon?: keyof typeof iconMap;
  iconSize?: 'sm' | 'md' | 'lg';
  actionLabel?: string;
  onAction?: () => void;
  className?: string;
  children?: ReactNode;
}

export default function EmptyState({
  title,
  description,
  icon,
  iconSize = 'md',
  actionLabel,
  onAction,
  className = '',
  children,
}: EmptyStateProps) {
  const IconComponent = icon ? iconMap[icon] : null;

  return (
    <div className={`flex flex-col items-center justify-center py-12 text-center ${className}`.trim()}>
      {IconComponent && (
        <IconComponent className={`${iconSizeMap[iconSize]} text-gray-600 mb-4`} />
      )}
      <h3 className="text-lg font-medium text-gray-300 mb-2">{title}</h3>
      {description && (
        <p className="text-sm text-gray-500 mb-4 max-w-md">{description}</p>
      )}
      {actionLabel && onAction && (
        <button
          onClick={onAction}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          {actionLabel}
        </button>
      )}
      {children}
    </div>
  );
}
