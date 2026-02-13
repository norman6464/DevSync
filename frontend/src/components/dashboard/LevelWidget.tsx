import { useTranslation } from 'react-i18next';
import { Star } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useMyLevel } from '../../hooks';

export default function LevelWidget() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const { levelInfo, loading } = useMyLevel();

  if (!user) return null;

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-3" />
        <div className="h-12 bg-gray-800 rounded animate-pulse mb-3" />
        <div className="h-3 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  const level = levelInfo?.level ?? 0;
  const totalXP = levelInfo?.total_xp ?? 0;
  const nextLevelXP = levelInfo?.next_level_xp ?? 100;
  const currentLevelXP = levelInfo?.current_level_xp ?? 0;
  const progressPercent = levelInfo?.progress_percent ?? 0;
  const xpToNext = nextLevelXP - totalXP;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
        <Star className="w-4 h-4 text-gray-400" />
        {t('level.title')}
      </h3>

      {/* Level display */}
      <div className="text-center py-2 mb-3">
        <div className="flex items-center justify-center gap-2">
          <span className="text-sm text-gray-400">{t('level.level')}</span>
          <span className="text-4xl font-bold text-white">{level}</span>
        </div>
        <p className="text-xs text-gray-400 mt-1">
          {t('level.xpProgress', { current: totalXP.toLocaleString(), next: nextLevelXP.toLocaleString() })}
        </p>
      </div>

      {/* Progress bar */}
      <div className="mb-2">
        <div className="h-2 bg-gray-800 rounded-full overflow-hidden">
          <div
            className="h-full bg-blue-600 rounded-full transition-all"
            style={{ width: `${Math.min(progressPercent, 100)}%` }}
          />
        </div>
        <div className="flex justify-between mt-1">
          <span className="text-[10px] text-gray-500">{currentLevelXP.toLocaleString()} XP</span>
          <span className="text-[10px] text-gray-500">{nextLevelXP.toLocaleString()} XP</span>
        </div>
      </div>

      {/* XP to next level */}
      <p className="text-xs text-gray-500 text-center">
        {t('level.xpNeeded', { xp: xpToNext.toLocaleString() })}
      </p>
    </div>
  );
}
