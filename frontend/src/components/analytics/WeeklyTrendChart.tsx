import { useTranslation } from 'react-i18next';
import type { WeeklyTrend } from '../../types/analytics';

interface Props {
  data: WeeklyTrend[];
  loading: boolean;
}

export default function WeeklyTrendChart({ data, loading }: Props) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="h-40 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
        <h3 className="text-sm font-medium text-white mb-4">{t('analytics.trendTitle')}</h3>
        <p className="text-xs text-gray-500 text-center py-6">{t('analytics.noData')}</p>
      </div>
    );
  }

  const maxMinutes = Math.max(...data.map((d) => d.total_minutes), 1);

  // 前週比計算
  const latestWeek = data[data.length - 1];
  const prevWeek = data.length >= 2 ? data[data.length - 2] : null;
  let growthRate: number | null = null;
  if (prevWeek && prevWeek.total_minutes > 0) {
    growthRate = ((latestWeek.total_minutes - prevWeek.total_minutes) / prevWeek.total_minutes) * 100;
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-medium text-white">{t('analytics.trendTitle')}</h3>
        {growthRate !== null && (
          <span className={`text-xs font-medium ${growthRate >= 0 ? 'text-green-400' : 'text-red-400'}`}>
            {growthRate >= 0 ? '+' : ''}{growthRate.toFixed(0)}%
          </span>
        )}
      </div>

      {/* バーチャート */}
      <div className="flex items-end gap-1 h-32 mb-2">
        {data.map((week) => {
          const height = (week.total_minutes / maxMinutes) * 100;
          const hours = Math.floor(week.total_minutes / 60);
          const mins = week.total_minutes % 60;
          return (
            <div key={week.week_start} className="flex-1 flex flex-col items-center">
              <span className="text-[10px] text-gray-500 mb-1">
                {hours > 0 ? `${hours}h` : `${mins}m`}
              </span>
              <div
                className="w-full bg-blue-500/80 rounded-t hover:bg-blue-400/80 transition-colors"
                style={{ height: `${Math.max(height, 2)}%` }}
                title={`${week.week_start}: ${week.total_minutes}${t('analytics.minutes')} (${week.log_count}${t('analytics.logs')})`}
              />
            </div>
          );
        })}
      </div>

      {/* 日付ラベル */}
      <div className="flex gap-1">
        {data.map((week, i) => (
          <div key={week.week_start} className="flex-1 text-center">
            {(i === 0 || i === data.length - 1 || i === Math.floor(data.length / 2)) && (
              <span className="text-[10px] text-gray-600">
                {week.week_start.slice(5)} {/* MM-DD */}
              </span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
