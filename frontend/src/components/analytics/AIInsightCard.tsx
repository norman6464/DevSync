import { useTranslation } from 'react-i18next';
import { Lightbulb, Clock, TrendingUp, TrendingDown, Flame, LayoutGrid } from 'lucide-react';
import type { AIInsight } from '../../types/analytics';
import { panelClass } from '../../constants/styles';

interface Props {
  insights: AIInsight[];
  loading: boolean;
}

const ICON_MAP: Record<string, typeof Lightbulb> = {
  peak_time: Clock,
  category_focus: LayoutGrid,
  weekly_growth_up: TrendingUp,
  weekly_growth_down: TrendingDown,
  streak_active: Flame,
  streak_record: Flame,
};

const COLOR_MAP: Record<string, string> = {
  peak_time: 'text-blue-400 bg-blue-400/10',
  category_focus: 'text-purple-400 bg-purple-400/10',
  weekly_growth_up: 'text-green-400 bg-green-400/10',
  weekly_growth_down: 'text-red-400 bg-red-400/10',
  streak_active: 'text-orange-400 bg-orange-400/10',
  streak_record: 'text-orange-400 bg-orange-400/10',
};

export default function AIInsightCard({ insights, loading }: Props) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className={panelClass}>
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="space-y-3">
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
        </div>
      </div>
    );
  }

  const formatInsight = (insight: AIInsight): string => {
    const params = { ...insight.params };
    // 曜日番号をi18nで翻訳
    if ('day_of_week' in params) {
      const dayLabels = t('analytics.dayLabels', { returnObjects: true }) as string[];
      params.day = dayLabels[params.day_of_week as number] || '';
    }
    return t(`analytics.insight_${insight.type}`, params as Record<string, string>);
  };

  return (
    <div className={panelClass}>
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-4">
        <Lightbulb className="w-4 h-4 text-yellow-400" />
        {t('analytics.insightsTitle')}
      </h3>

      {insights.length === 0 ? (
        <p className="text-xs text-gray-500 text-center py-4">{t('analytics.noInsights')}</p>
      ) : (
        <div className="space-y-3">
          {insights.map((insight, idx) => {
            const Icon = ICON_MAP[insight.type] || Lightbulb;
            const colorClass = COLOR_MAP[insight.type] || 'text-gray-400 bg-gray-400/10';
            return (
              <div key={idx} className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
                <div className={`p-1.5 rounded-lg shrink-0 ${colorClass}`}>
                  <Icon className="w-4 h-4" />
                </div>
                <p className="text-xs text-gray-300 leading-relaxed">{formatInsight(insight)}</p>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
