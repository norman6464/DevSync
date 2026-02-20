import { Clock, Calendar, FileText } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { StreakInfo } from '../../types/learningLog';

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
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3">
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
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3">
        <Calendar className="w-8 h-8 text-orange-400" />
        <div>
          <p className="text-xs text-gray-400">{t('learningLogs.currentStreak')}</p>
          <p className="text-lg font-bold text-white">{streakInfo?.current_streak ?? 0}{t('learningLogs.days')}</p>
        </div>
      </div>
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3 col-span-2 md:col-span-1">
        <FileText className="w-8 h-8 text-green-400" />
        <div>
          <p className="text-xs text-gray-400">{t('learningLogs.logCount')}</p>
          <p className="text-lg font-bold text-white">{logCount}</p>
        </div>
      </div>
    </div>
  );
}
