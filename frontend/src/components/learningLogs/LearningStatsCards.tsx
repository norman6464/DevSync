import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Clock, BookOpen, TrendingUp, PieChart } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useAsyncData } from '../../hooks/useAsyncData';
import { getMyLogs } from '../../api/learningLogs';
import { getCategoryInfo, getCategoryColor } from '../../constants/learningLogs';
import type { LearningLog, LogCategory } from '../../types/learningLog';

export default function LearningStatsCards() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const { data: logs, loading } = useAsyncData(
    async () => {
      if (!user) return [];
      const res = await getMyLogs();
      return res.data;
    },
    { initialData: [] as LearningLog[], deps: [user?.id], enabled: !!user }
  );

  const stats = useMemo(() => {
    const totalDuration = logs.reduce((sum, log) => sum + log.duration, 0);
    const totalHours = Math.round((totalDuration / 60) * 10) / 10;
    const avgDuration = logs.length > 0 ? Math.round((totalDuration / logs.length) * 10) / 10 : 0;
    const avgHours = Math.round((avgDuration / 60) * 10) / 10;

    // カテゴリ別集計
    const categoryStats = logs.reduce((acc, log) => {
      const cat = log.category;
      if (!acc[cat]) {
        acc[cat] = { count: 0, duration: 0 };
      }
      acc[cat].count++;
      acc[cat].duration += log.duration;
      return acc;
    }, {} as Record<LogCategory, { count: number; duration: number }>);

    // カテゴリ別割合
    const categoryPercentages = Object.entries(categoryStats).map(([category, stat]) => ({
      category: category as LogCategory,
      count: stat.count,
      duration: stat.duration,
      percentage: totalDuration > 0 ? Math.round((stat.duration / totalDuration) * 100) : 0,
    }));

    // 割合でソート
    categoryPercentages.sort((a, b) => b.percentage - a.percentage);

    return {
      totalDuration,
      totalHours,
      totalLogs: logs.length,
      avgDuration,
      avgHours,
      categoryPercentages,
    };
  }, [logs]);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-2" />
          <div className="h-8 bg-gray-800 rounded animate-pulse w-3/4" />
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-2" />
          <div className="h-8 bg-gray-800 rounded animate-pulse w-3/4" />
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-2" />
          <div className="h-8 bg-gray-800 rounded animate-pulse w-3/4" />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-white flex items-center gap-2">
        <TrendingUp className="w-4 h-4 text-blue-400" />
        学習統計
      </h3>

      {/* Main Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Total Duration */}
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-2">
            <Clock className="w-4 h-4 text-blue-400" />
            <span className="text-xs text-gray-400">総学習時間</span>
          </div>
          <div className="text-2xl font-bold text-white">
            {stats.totalHours}
            <span className="text-sm text-gray-400 ml-1">時間</span>
          </div>
          <div className="text-xs text-gray-500 mt-1">
            {stats.totalDuration}分
          </div>
        </div>

        {/* Total Logs */}
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-2">
            <BookOpen className="w-4 h-4 text-green-400" />
            <span className="text-xs text-gray-400">学習ログ数</span>
          </div>
          <div className="text-2xl font-bold text-white">
            {stats.totalLogs}
            <span className="text-sm text-gray-400 ml-1">件</span>
          </div>
        </div>

        {/* Average Duration */}
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-2">
            <TrendingUp className="w-4 h-4 text-purple-400" />
            <span className="text-xs text-gray-400">平均学習時間</span>
          </div>
          <div className="text-2xl font-bold text-white">
            {stats.avgHours}
            <span className="text-sm text-gray-400 ml-1">時間</span>
          </div>
          <div className="text-xs text-gray-500 mt-1">
            {stats.avgDuration}分/回
          </div>
        </div>
      </div>

      {/* Category Stats */}
      {stats.categoryPercentages.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-4">
          <div className="flex items-center gap-2 mb-4">
            <PieChart className="w-4 h-4 text-orange-400" />
            <span className="text-sm font-medium text-white">カテゴリ別学習時間</span>
          </div>

          <div className="space-y-3">
            {stats.categoryPercentages.map(({ category, duration, percentage }) => {
              const categoryInfo = getCategoryInfo(category);
              const Icon = categoryInfo.Icon;
              const colorClass = getCategoryColor(category);
              const hours = Math.round((duration / 60) * 10) / 10;

              return (
                <div key={category} className="space-y-1">
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-2">
                      <Icon className="w-3.5 h-3.5 text-gray-400" />
                      <span className="text-white">{t(categoryInfo.label)}</span>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className="text-gray-400">{hours}h</span>
                      <span className="text-blue-400 font-medium min-w-[3rem] text-right">
                        {percentage}%
                      </span>
                    </div>
                  </div>
                  <div className="h-2 bg-gray-800 rounded-full overflow-hidden">
                    <div
                      className={`h-full ${colorClass.split(' ')[1]} transition-all`}
                      style={{ width: `${percentage}%` }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {stats.totalLogs === 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-8 text-center">
          <BookOpen className="w-12 h-12 mx-auto mb-3 text-gray-700" />
          <p className="text-gray-400 text-sm">
            まだ学習ログがありません
          </p>
          <p className="text-gray-500 text-xs mt-1">
            学習ログを記録して統計を確認しましょう
          </p>
        </div>
      )}
    </div>
  );
}
