import { Clock, Calendar, FileText, Flame } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { StreakInfo } from '../../types/learningLog';
import { panelClass } from '../../constants/styles';

interface WeeklySummaryCardProps {
  weeklyDuration: number;
  streakInfo: StreakInfo | null;
  logCount: number;
}

export default function WeeklySummaryCard({
  weeklyDuration,
  streakInfo,
  logCount,
}: WeeklySummaryCardProps) {
  const { t } = useTranslation();

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
      <div className={`${panelClass} flex items-center gap-3`}>
        <Clock className="w-8 h-8 text-blue-400" />
        <div>
          <p className="text-xs text-gray-400">{t('learningLogs.weeklyDuration')}</p>
          <p className="text-lg font-bold text-white">
            {weeklyDuration >= 60
              ? t('learningLogs.hoursMinutes', { hours: Math.floor(weeklyDuration / 60), minutes: weeklyDuration % 60 })
              : t('learningLogs.durationMinutes', { minutes: weeklyDuration })}
          </p>
        </div>
      </div>
      <div className={`${panelClass} flex items-center gap-3`}>
        <Calendar className="w-8 h-8 text-orange-400" />
        <div>
          <p className="text-xs text-gray-400">{t('learningLogs.currentStreak')}</p>
          <p className="text-lg font-bold text-white">
            {streakInfo?.current_streak ?? 0}{t('learningLogs.days')}
            {(streakInfo?.current_streak ?? 0) >= 7 && (
              <span className="inline-flex items-center gap-0.5 ml-2 px-1.5 py-0.5 text-xs rounded bg-orange-400/10 text-orange-400 font-medium align-middle">
                <Flame className="w-3 h-3" />
                {t('learningLogs.streakAchieved')}
              </span>
            )}
          </p>
        </div>
      </div>
      <div className={`${panelClass} flex items-center gap-3 col-span-2 md:col-span-1`}>
        <FileText className="w-8 h-8 text-green-400" />
        <div>
          <p className="text-xs text-gray-400">{t('learningLogs.logCount')}</p>
          <p className="text-lg font-bold text-white">{logCount}</p>
        </div>
      </div>
    </div>
  );
}
