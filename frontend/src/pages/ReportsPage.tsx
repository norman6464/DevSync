import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { BarChart3, FileText, Users, Target, MessageCircle, Heart, type LucideIcon } from 'lucide-react';
import { useReport } from '../hooks';
import LoadingSpinner from '../components/common/LoadingSpinner';
import ReportCharts from '../components/reports/ReportCharts';

function formatDate(dateStr: string) {
  const date = new Date(dateStr);
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
}

function getTrendColor(diff: number) {
  if (diff > 0) return 'text-green-400';
  if (diff < 0) return 'text-red-400';
  return 'text-gray-400';
}

function getTrendIcon(d: number) {
  if (d > 0) return '↑';
  if (d < 0) return '↓';
  return '—';
}

export default function ReportsPage() {
  const { t } = useTranslation();
  const { report, comparison, loading, period, setPeriod } = useReport();

  const handleSetWeekly = useCallback(() => setPeriod('weekly'), [setPeriod]);
  const handleSetMonthly = useCallback(() => setPeriod('monthly'), [setPeriod]);

  const dailyContributions = report?.daily_contributions;
  const maxContribution = useMemo(
    () =>
      dailyContributions
        ? Math.max(...dailyContributions.map((d) => d.contributions), 1)
        : 1,
    [dailyContributions]
  );

  if (loading) {
    return (
      <div className="py-12">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t('reports.title')}</h1>
        <div className="flex bg-gray-800 rounded-lg p-1">
          <button
            onClick={handleSetWeekly}
            className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
              period === 'weekly'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('reports.weekly')}
          </button>
          <button
            onClick={handleSetMonthly}
            className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
              period === 'monthly'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('reports.monthly')}
          </button>
        </div>
      </div>

      {/* Period Info */}
      {report && (
        <p className="text-gray-400 text-sm">
          {formatDate(report.start_date)} - {formatDate(report.end_date)}
        </p>
      )}

      {/* Stats Grid */}
      {report && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatCard
            label={t('reports.contributions')}
            value={report.total_contributions}
            diff={comparison?.contributions_diff || 0}
            icon={BarChart3}
          />
          <StatCard
            label={t('reports.posts')}
            value={report.posts_created}
            diff={comparison?.posts_diff || 0}
            icon={FileText}
          />
          <StatCard
            label={t('reports.newFollowers')}
            value={report.new_followers}
            diff={comparison?.followers_diff || 0}
            icon={Users}
          />
          <StatCard
            label={t('reports.goalsCompleted')}
            value={report.goals_completed}
            diff={comparison?.goals_diff || 0}
            icon={Target}
          />
        </div>
      )}

      {/* Activity Overview */}
      {report && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <div className="flex items-center gap-2 mb-2">
              <MessageCircle className="w-5 h-5 text-gray-400" />
              <span className="text-sm text-gray-400">{t('reports.comments')}</span>
            </div>
            <p className="text-2xl font-bold">{report.comments_created}</p>
          </div>
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <div className="flex items-center gap-2 mb-2">
              <Heart className="w-5 h-5 text-gray-400" />
              <span className="text-sm text-gray-400">{t('reports.likesReceived')}</span>
            </div>
            <p className="text-2xl font-bold">{report.likes_received}</p>
          </div>
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <div className="flex items-center gap-2 mb-2">
              <MessageCircle className="w-5 h-5 text-gray-400" />
              <span className="text-sm text-gray-400">{t('reports.messages')}</span>
            </div>
            <p className="text-2xl font-bold">{report.messages_exchanged}</p>
          </div>
        </div>
      )}

      {/* Trend Indicator */}
      {comparison && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.trend')}</h2>
          <div className="flex items-center gap-4">
            <div
              className={`text-4xl font-bold ${getTrendColor(
                comparison.trend_percentage
              )}`}
            >
              {comparison.trend_percentage > 0 ? '+' : ''}
              {comparison.trend_percentage.toFixed(1)}%
            </div>
            <div className="text-gray-400 text-sm">
              {comparison.trend_percentage > 0
                ? t('reports.trendUp')
                : comparison.trend_percentage < 0
                ? t('reports.trendDown')
                : t('reports.trendStable')}
            </div>
          </div>
        </div>
      )}

      {/* Charts: Daily Activity, Top Languages, Goals Progress */}
      {report && <ReportCharts report={report} maxContribution={maxContribution} />}
    </div>
  );
}

interface StatCardProps {
  label: string;
  value: number;
  diff: number;
  icon: LucideIcon;
}

function StatCard({ label, value, diff, icon: Icon }: StatCardProps) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-center gap-2 mb-2">
        <Icon className="w-5 h-5 text-gray-400" />
        <span className="text-sm text-gray-400">{label}</span>
      </div>
      <div className="flex items-end gap-2">
        <p className="text-2xl font-bold">{value}</p>
        <span className={`text-sm ${getTrendColor(diff)}`}>
          {getTrendIcon(diff)} {Math.abs(diff)}
        </span>
      </div>
    </div>
  );
}
