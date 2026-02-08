import { useTranslation } from 'react-i18next';
import { Flame, Trophy, Calendar } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useStreak } from '../../hooks';

export default function StreakWidget() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const { streakInfo, loading } = useStreak(user?.id);

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-3" />
        <div className="h-12 bg-gray-800 rounded animate-pulse mb-3" />
        <div className="grid grid-cols-2 gap-2">
          <div className="h-10 bg-gray-800 rounded animate-pulse" />
          <div className="h-10 bg-gray-800 rounded animate-pulse" />
        </div>
      </div>
    );
  }

  const currentStreak = streakInfo?.current_streak ?? 0;
  const longestStreak = streakInfo?.longest_streak ?? 0;
  const totalDays = streakInfo?.total_days ?? 0;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
        <Flame className="w-4 h-4 text-orange-400" />
        {t('streak.title')}
      </h3>

      {/* Main streak display */}
      <div className="text-center py-2 mb-3">
        <div className="flex items-center justify-center gap-2">
          <span className="text-4xl font-bold text-orange-400">
            {currentStreak}
          </span>
          <span className="text-sm text-gray-400">
            {t('streak.daysConsecutive')}
          </span>
        </div>
        {currentStreak > 0 ? (
          <p className="text-xs text-orange-300/70 mt-1">{t('streak.keepItUp')}</p>
        ) : (
          <p className="text-xs text-gray-500 mt-1">{t('streak.noStreak')}</p>
        )}
      </div>

      {/* Sub-stats */}
      <div className="grid grid-cols-2 gap-2">
        <div className="bg-gray-800/50 rounded-lg p-2 text-center">
          <Trophy className="w-3.5 h-3.5 text-yellow-400 mx-auto mb-1" />
          <div className="text-sm font-bold text-white">{longestStreak}</div>
          <div className="text-[10px] text-gray-500">{t('streak.longestStreak')}</div>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-2 text-center">
          <Calendar className="w-3.5 h-3.5 text-purple-400 mx-auto mb-1" />
          <div className="text-sm font-bold text-white">{totalDays}</div>
          <div className="text-[10px] text-gray-500">{t('streak.totalDays')}</div>
        </div>
      </div>
    </div>
  );
}
