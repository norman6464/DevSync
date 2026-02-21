import { useTranslation } from 'react-i18next';
import { Trophy, ThumbsUp, MinusCircle, TrendingDown } from 'lucide-react';
import type { ProductivityScore } from '../../types/analytics';
import { panelClass } from '../../constants/styles';

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
      <div className={panelClass}>
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
    <div className={panelClass}>
      <h3 className="text-sm font-medium text-white mb-4">{t('analytics.productivityTitle')}</h3>

      {/* 総合スコア */}
      <div className="text-center mb-5">
        <div className={`text-4xl font-bold ${getScoreColor(score.overall_score)}`}>
          {score.overall_score.toFixed(0)}
        </div>
        <p className="text-xs text-gray-500 mt-1">{t('analytics.overallScore')}</p>
        {score.overall_score >= 80 ? (
          <span className="inline-flex items-center gap-1 mt-2 px-2 py-0.5 text-xs rounded bg-green-400/10 text-green-400 font-medium">
            <Trophy className="w-3 h-3" />
            {t('analytics.rankExcellent')}
          </span>
        ) : score.overall_score >= 60 ? (
          <span className="inline-flex items-center gap-1 mt-2 px-2 py-0.5 text-xs rounded bg-blue-400/10 text-blue-400 font-medium">
            <ThumbsUp className="w-3 h-3" />
            {t('analytics.rankGood')}
          </span>
        ) : score.overall_score >= 40 ? (
          <span className="inline-flex items-center gap-1 mt-2 px-2 py-0.5 text-xs rounded bg-yellow-400/10 text-yellow-400 font-medium">
            <MinusCircle className="w-3 h-3" />
            {t('analytics.rankFair')}
          </span>
        ) : (
          <span className="inline-flex items-center gap-1 mt-2 px-2 py-0.5 text-xs rounded bg-red-400/10 text-red-400 font-medium">
            <TrendingDown className="w-3 h-3" />
            {t('analytics.rankNeedsWork')}
          </span>
        )}
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
