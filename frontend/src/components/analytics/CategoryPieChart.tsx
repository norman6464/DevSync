import { useTranslation } from 'react-i18next';
import type { CategoryBreakdown } from '../../types/analytics';

interface Props {
  data: CategoryBreakdown[];
  loading: boolean;
}

const CATEGORY_COLORS: Record<string, string> = {
  coding: 'bg-blue-500',
  reading: 'bg-green-500',
  course: 'bg-purple-500',
  meetup: 'bg-orange-500',
  other: 'bg-gray-500',
};

const CATEGORY_TEXT_COLORS: Record<string, string> = {
  coding: 'text-blue-400',
  reading: 'text-green-400',
  course: 'text-purple-400',
  meetup: 'text-orange-400',
  other: 'text-gray-400',
};

export default function CategoryPieChart({ data, loading }: Props) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="h-32 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  const totalMinutes = data.reduce((sum, c) => sum + c.total_minutes, 0);
  const totalHours = Math.floor(totalMinutes / 60);
  const remainMinutes = totalMinutes % 60;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <h3 className="text-sm font-medium text-white mb-4">{t('analytics.categoryTitle')}</h3>

      {data.length === 0 ? (
        <p className="text-xs text-gray-500 text-center py-6">{t('analytics.noData')}</p>
      ) : (
        <>
          {/* 合計時間 */}
          <div className="text-center mb-4">
            <span className="text-2xl font-bold text-white">{totalHours}</span>
            <span className="text-sm text-gray-400 ml-1">{t('analytics.hours')}</span>
            {remainMinutes > 0 && (
              <>
                <span className="text-lg font-bold text-white ml-2">{remainMinutes}</span>
                <span className="text-sm text-gray-400 ml-1">{t('analytics.min')}</span>
              </>
            )}
          </div>

          {/* 横バー（total_minutesベースで幅を計算し、丸め誤差を防止） */}
          <div className="h-3 flex rounded-full overflow-hidden mb-4">
            {data.map((cat) => (
              <div
                key={cat.category}
                className={`${CATEGORY_COLORS[cat.category] || 'bg-gray-600'} transition-all`}
                style={{ width: totalMinutes > 0 ? `${(cat.total_minutes / totalMinutes) * 100}%` : '0%' }}
              />
            ))}
          </div>

          {/* カテゴリリスト */}
          <div className="space-y-2">
            {data.map((cat) => (
              <div key={cat.category} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className={`w-2.5 h-2.5 rounded-full ${CATEGORY_COLORS[cat.category] || 'bg-gray-600'}`} />
                  <span className={`text-xs ${CATEGORY_TEXT_COLORS[cat.category] || 'text-gray-400'}`}>
                    {t(`analytics.category_${cat.category}`, { defaultValue: cat.category })}
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs text-gray-500">{cat.log_count}{t('analytics.logs')}</span>
                  <span className="text-xs text-white font-medium w-12 text-right">
                    {cat.percentage.toFixed(0)}%
                  </span>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
