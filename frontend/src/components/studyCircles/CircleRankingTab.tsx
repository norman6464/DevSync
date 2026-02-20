import { useTranslation } from 'react-i18next';
import type { CircleMemberStreak } from '../../types/studyCircle';
import Avatar from '../common/Avatar';

interface CircleRankingTabProps {
  streaks: CircleMemberStreak[];
}

export default function CircleRankingTab({ streaks }: CircleRankingTabProps) {
  const { t } = useTranslation();

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <h3 className="text-sm font-medium text-white mb-3">{t('studyCircle.streak.title')}</h3>
      {streaks.length === 0 ? (
        <p className="text-xs text-gray-500 text-center py-4">{t('studyCircle.streak.noStreaks')}</p>
      ) : (
        <div className="space-y-2">
          {streaks.map((s, i) => (
            <div
              key={s.user_id}
              className="flex items-center gap-3 p-2.5 rounded-lg bg-gray-800/30"
            >
              <span className={`text-sm font-bold w-6 text-center ${
                i === 0 ? 'text-yellow-400' : i === 1 ? 'text-gray-300' : i === 2 ? 'text-amber-600' : 'text-gray-600'
              }`}>
                {i + 1}
              </span>
              <Avatar name={s.user_name} avatarUrl={s.avatar_url} size="sm" />
              <div className="flex-1 min-w-0">
                <p className="text-sm text-white truncate">{s.user_name}</p>
                <p className="text-[10px] text-gray-500">
                  {t('studyCircle.streak.totalCheckins')}: {s.total_checkins}
                </p>
              </div>
              <div className="text-right shrink-0">
                <div className="text-lg font-bold text-orange-400">{s.current_streak}</div>
                <div className="text-[10px] text-gray-500">{t('studyCircle.streak.currentStreak')}</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
