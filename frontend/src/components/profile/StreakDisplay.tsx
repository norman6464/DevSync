import { useTranslation } from 'react-i18next';
import { Flame, Trophy, Calendar } from 'lucide-react';
import type { StreakInfo } from '../../types/learningLog';

interface StreakDisplayProps {
  streakInfo: StreakInfo | null;
}

export default function StreakDisplay({ streakInfo }: StreakDisplayProps) {
  const { t } = useTranslation();

  if (!streakInfo) return null;

  const { current_streak, longest_streak, total_days, last_log_date } = streakInfo;

  // Don't show if user has never logged
  if (total_days === 0) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
      <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide flex items-center gap-2 mb-4">
        <Flame className="w-5 h-5 text-orange-400" />
        {t('streak.title')}
      </h2>

      <div className="grid grid-cols-3 gap-4">
        <div className="text-center">
          <div className="flex items-center justify-center gap-1.5 mb-1">
            <Flame className="w-4 h-4 text-orange-400" />
            <span className="text-2xl font-bold text-orange-400">{current_streak}</span>
          </div>
          <div className="text-xs text-gray-500">{t('streak.currentStreak')}</div>
        </div>

        <div className="text-center">
          <div className="flex items-center justify-center gap-1.5 mb-1">
            <Trophy className="w-4 h-4 text-yellow-400" />
            <span className="text-2xl font-bold text-yellow-400">{longest_streak}</span>
          </div>
          <div className="text-xs text-gray-500">{t('streak.longestStreak')}</div>
        </div>

        <div className="text-center">
          <div className="flex items-center justify-center gap-1.5 mb-1">
            <Calendar className="w-4 h-4 text-purple-400" />
            <span className="text-2xl font-bold text-purple-400">{total_days}</span>
          </div>
          <div className="text-xs text-gray-500">{t('streak.totalDays')}</div>
        </div>
      </div>

      {last_log_date && (
        <p className="text-xs text-gray-500 text-center mt-3">
          {t('streak.lastLogDate')}: {last_log_date}
        </p>
      )}
    </div>
  );
}
