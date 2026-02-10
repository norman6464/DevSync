import { useTranslation } from 'react-i18next';

interface LevelBadgeProps {
  level: number;
  size?: 'sm' | 'md';
}

/** レベルに応じた色を返す */
function getLevelColor(level: number): string {
  if (level >= 41) return 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30';
  if (level >= 31) return 'bg-purple-500/20 text-purple-400 border-purple-500/30';
  if (level >= 21) return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
  if (level >= 11) return 'bg-green-500/20 text-green-400 border-green-500/30';
  return 'bg-gray-500/20 text-gray-400 border-gray-500/30';
}

export default function LevelBadge({ level, size = 'sm' }: LevelBadgeProps) {
  const { t } = useTranslation();
  const colorClass = getLevelColor(level);
  const sizeClass = size === 'md'
    ? 'px-2.5 py-1 text-sm'
    : 'px-1.5 py-0.5 text-xs';

  return (
    <span className={`inline-flex items-center gap-1 font-bold rounded-md border ${colorClass} ${sizeClass}`}>
      {t('level.level')} {level}
    </span>
  );
}
