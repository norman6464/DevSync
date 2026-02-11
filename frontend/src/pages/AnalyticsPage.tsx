import { useTranslation } from 'react-i18next';
import { BarChart3 } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import {
  useHeatmap,
  useCategoryBreakdown,
  useProductivityScore,
  useWeeklyTrends,
  useInsights,
} from '../hooks';
import HeatmapChart from '../components/analytics/HeatmapChart';
import CategoryPieChart from '../components/analytics/CategoryPieChart';
import ProductivityScoreCard from '../components/analytics/ProductivityScoreCard';
import WeeklyTrendChart from '../components/analytics/WeeklyTrendChart';
import AIInsightCard from '../components/analytics/AIInsightCard';

export default function AnalyticsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const { heatmap, loading: heatmapLoading } = useHeatmap(user?.id);
  const { categories, loading: categoriesLoading } = useCategoryBreakdown(user?.id);
  const { score, loading: scoreLoading } = useProductivityScore(user?.id);
  const { trends, loading: trendsLoading } = useWeeklyTrends(user?.id, 12);
  const { insights, loading: insightsLoading } = useInsights();

  return (
    <div className="space-y-6">
      {/* ヘッダー */}
      <div className="flex items-center gap-3">
        <BarChart3 className="w-6 h-6 text-blue-400" />
        <div>
          <h1 className="text-xl font-bold text-white">{t('analytics.pageTitle')}</h1>
          <p className="text-sm text-gray-400">{t('analytics.pageDescription')}</p>
        </div>
      </div>

      {/* 上段: 生産性スコア + AIインサイト */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <ProductivityScoreCard score={score} loading={scoreLoading} />
        <AIInsightCard insights={insights || []} loading={insightsLoading} />
      </div>

      {/* 中段: 週間トレンド + カテゴリ別 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <WeeklyTrendChart data={trends || []} loading={trendsLoading} />
        <CategoryPieChart data={categories || []} loading={categoriesLoading} />
      </div>

      {/* 下段: ヒートマップ（フル幅） */}
      <HeatmapChart data={heatmap || []} loading={heatmapLoading} />
    </div>
  );
}
