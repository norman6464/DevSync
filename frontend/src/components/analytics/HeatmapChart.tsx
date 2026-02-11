import { useTranslation } from 'react-i18next';
import type { HeatmapEntry } from '../../types/analytics';

interface Props {
  data: HeatmapEntry[];
  loading: boolean;
}

const HOURS = Array.from({ length: 24 }, (_, i) => i);
const DAYS = [0, 1, 2, 3, 4, 5, 6]; // 日〜土

function getColor(minutes: number, max: number): string {
  if (minutes === 0 || max === 0) return 'bg-gray-800';
  const ratio = minutes / max;
  if (ratio > 0.75) return 'bg-green-500';
  if (ratio > 0.5) return 'bg-green-600';
  if (ratio > 0.25) return 'bg-green-700';
  return 'bg-green-900';
}

export default function HeatmapChart({ data, loading }: Props) {
  const { t } = useTranslation();
  const dayLabels = t('analytics.dayLabels', { returnObjects: true }) as string[];

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="h-40 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  // データをマップに変換
  const dataMap = new Map<string, number>();
  let maxMinutes = 0;
  for (const entry of data) {
    const key = `${entry.day_of_week}-${entry.hour}`;
    dataMap.set(key, entry.total_minutes);
    if (entry.total_minutes > maxMinutes) maxMinutes = entry.total_minutes;
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <h3 className="text-sm font-medium text-white mb-4">{t('analytics.heatmapTitle')}</h3>

      <div className="overflow-x-auto">
        <div className="min-w-[600px]">
          {/* 時間ラベル */}
          <div className="flex gap-[2px] mb-1 ml-10">
            {HOURS.filter((h) => h % 3 === 0).map((h) => (
              <span
                key={h}
                className="text-[10px] text-gray-500"
                style={{ width: `${(100 / 8)}%`, textAlign: 'center' }}
              >
                {h}
              </span>
            ))}
          </div>

          {/* ヒートマップ本体 */}
          {DAYS.map((day) => (
            <div key={day} className="flex items-center gap-[2px] mb-[2px]">
              <span className="text-[10px] text-gray-500 w-10 text-right pr-2 shrink-0">
                {dayLabels[day]}
              </span>
              {HOURS.map((hour) => {
                const minutes = dataMap.get(`${day}-${hour}`) || 0;
                return (
                  <div
                    key={hour}
                    className={`flex-1 h-4 rounded-sm ${getColor(minutes, maxMinutes)} transition-colors`}
                    title={`${dayLabels[day]} ${hour}:00 — ${minutes}${t('analytics.minutes')}`}
                  />
                );
              })}
            </div>
          ))}

          {/* 凡例 */}
          <div className="flex items-center justify-end gap-1 mt-3">
            <span className="text-[10px] text-gray-500">{t('analytics.less')}</span>
            <div className="w-3 h-3 rounded-sm bg-gray-800" />
            <div className="w-3 h-3 rounded-sm bg-green-900" />
            <div className="w-3 h-3 rounded-sm bg-green-700" />
            <div className="w-3 h-3 rounded-sm bg-green-600" />
            <div className="w-3 h-3 rounded-sm bg-green-500" />
            <span className="text-[10px] text-gray-500">{t('analytics.more')}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
