import { TrendingUp, TrendingDown, BarChart3, Clock, Users, FileText } from 'lucide-react';

const iconMap = {
  'trending-up': TrendingUp,
  'trending-down': TrendingDown,
  'bar-chart': BarChart3,
  clock: Clock,
  users: Users,
  file: FileText,
};

interface StatCardProps {
  value: number;
  label: string;
  icon?: keyof typeof iconMap;
  change?: number;
  formatted?: string;
  suffix?: string;
  className?: string;
}

export default function StatCard({
  value,
  label,
  icon,
  change,
  formatted,
  suffix,
  className = '',
}: StatCardProps) {
  const IconComponent = icon ? iconMap[icon] : null;

  const changeColor =
    change === undefined
      ? ''
      : change > 0
        ? 'text-green-400'
        : change < 0
          ? 'text-red-400'
          : 'text-gray-400';

  return (
    <div className={`bg-gray-900 border border-gray-800 rounded-lg p-6 ${className}`.trim()}>
      <div className="flex items-center justify-between mb-4">
        <span className="text-sm text-gray-400">{label}</span>
        {IconComponent && (
          <IconComponent className="w-5 h-5 text-gray-500" />
        )}
      </div>
      <div className="flex items-baseline gap-2">
        <span className="text-3xl font-bold text-white">
          {formatted ?? value}
        </span>
        {suffix && <span className="text-sm text-gray-400">{suffix}</span>}
      </div>
      {change !== undefined && (
        <div className={`mt-2 text-sm ${changeColor}`}>
          {change > 0 ? '+' : ''}{change}%
        </div>
      )}
    </div>
  );
}
