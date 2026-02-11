import { useTranslation } from 'react-i18next';
import { Lightbulb, Clock, TrendingUp, Flame, LayoutGrid } from 'lucide-react';
import type { AIInsight } from '../../types/analytics';

interface Props {
  insights: AIInsight[];
  loading: boolean;
}

const ICON_MAP: Record<string, typeof Lightbulb> = {
  peak_time: Clock,
  category_focus: LayoutGrid,
  weekly_growth: TrendingUp,
  streak_trend: Flame,
};

const COLOR_MAP: Record<string, string> = {
  peak_time: 'text-blue-400 bg-blue-400/10',
  category_focus: 'text-purple-400 bg-purple-400/10',
  weekly_growth: 'text-green-400 bg-green-400/10',
  streak_trend: 'text-orange-400 bg-orange-400/10',
};

export default function AIInsightCard({ insights, loading }: Props) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/3 mb-4" />
        <div className="space-y-3">
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
          <div className="h-16 bg-gray-800 rounded animate-pulse" />
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
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
                <p className="text-xs text-gray-300 leading-relaxed">{insight.message}</p>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
