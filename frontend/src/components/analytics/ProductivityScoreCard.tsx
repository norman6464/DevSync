import { useTranslation } from 'react-i18next';
import type { ProductivityScore } from '../../types/analytics';

interface Props {
  score: ProductivityScore | null;
  loading: boolean;
}

function getScoreColor(value: number): string {
  if (value >= 80) return 'text-green-400';
  if (value >= 60) return 'text-blue-400';
  if (value >= 40) return 'text-yellow-400';
  return 'text-red-400';
}

function getBarColor(value: number): string {
  if (value >= 80) return 'bg-green-500';
  if (value >= 60) return 'bg-blue-500';
  if (value >= 40) return 'bg-yellow-500';
  return 'bg-red-500';
}

export default function ProductivityScoreCard({ score, loading }: Props) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="h-24 bg-gray-800 rounded animate-pulse" />
      </div>
    );
  }

  if (!score) return null;

  const metrics = [
    { label: t('analytics.pomodoroRate'), value: score.pomodoro_rate, weight: '30%' },
    { label: t('analytics.goalRate'), value: score.goal_rate, weight: '40%' },
    { label: t('analytics.streakConsistency'), value: score.streak_consistency, weight: '30%' },
  ];

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <h3 className="text-sm font-medium text-white mb-4">{t('analytics.productivityTitle')}</h3>

      {/* 総合スコア */}
      <div className="text-center mb-5">
        <div className={`text-4xl font-bold ${getScoreColor(score.overall_score)}`}>
          {score.overall_score.toFixed(0)}
        </div>
        <p className="text-xs text-gray-500 mt-1">{t('analytics.overallScore')}</p>
      </div>

      {/* 各指標 */}
      <div className="space-y-3">
        {metrics.map((metric) => (
          <div key={metric.label}>
            <div className="flex items-center justify-between mb-1">
              <span className="text-xs text-gray-400">{metric.label}</span>
              <div className="flex items-center gap-2">
                <span className="text-[10px] text-gray-600">{metric.weight}</span>
                <span className={`text-xs font-medium ${getScoreColor(metric.value)}`}>
                  {metric.value.toFixed(0)}%
                </span>
              </div>
            </div>
            <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full transition-all ${getBarColor(metric.value)}`}
                style={{ width: `${Math.min(metric.value, 100)}%` }}
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
