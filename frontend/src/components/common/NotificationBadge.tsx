import type { ReactNode } from 'react';

const colorMap = {
  red: 'bg-red-500',
  blue: 'bg-blue-500',
  green: 'bg-green-500',
  yellow: 'bg-yellow-500',
};

interface NotificationBadgeProps {
  children: ReactNode;
  count: number;
  max?: number;
  dot?: boolean;
  color?: keyof typeof colorMap;
  pulse?: boolean;
  showZero?: boolean;
  className?: string;
}

export default function NotificationBadge({
  children,
  count,
  max = 99,
  dot = false,
  color = 'red',
  pulse = false,
  showZero = false,
  className = '',
}: NotificationBadgeProps) {
  const showBadge = count > 0 || showZero;
  const displayCount = count > max ? `${max}+` : count;

  return (
    <div className={`relative inline-flex ${className}`.trim()}>
      {children}
      {showBadge && (
        dot ? (
          <span
            className={`absolute -top-1 -right-1 w-2 h-2 rounded-full ${colorMap[color]} ${pulse ? 'animate-pulse' : ''}`}
          />
        ) : (
          <span
            className={`absolute -top-1 -right-1 min-w-[1.25rem] h-5 flex items-center justify-center px-1 text-xs font-bold text-white rounded-full ${colorMap[color]} ${pulse ? 'animate-pulse' : ''}`}
          >
            {displayCount}
          </span>
        )
      )}
    </div>
  );
}
