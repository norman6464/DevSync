import { useTranslation } from 'react-i18next';
import { BarChart3, FileText, Users, Target, MessageCircle, Heart, type LucideIcon } from 'lucide-react';
import { useReport } from '../hooks';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function ReportsPage() {
  const { t } = useTranslation();
  const { report, comparison, loading, period, setPeriod } = useReport();

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
  };

  const formatBytes = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  const getTrendColor = (diff: number) => {
    if (diff > 0) return 'text-green-400';
    if (diff < 0) return 'text-red-400';
    return 'text-gray-400';
  };

  const maxContribution = report?.daily_contributions
    ? Math.max(...report.daily_contributions.map((d) => d.contributions), 1)
    : 1;

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
            onClick={() => setPeriod('weekly')}
            className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
              period === 'weekly'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('reports.weekly')}
          </button>
          <button
            onClick={() => setPeriod('monthly')}
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
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <MessageCircle className="w-5 h-5 text-gray-400" />
              <span className="text-sm text-gray-400">{t('reports.comments')}</span>
            </div>
            <p className="text-2xl font-bold">{report.comments_created}</p>
          </div>
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Heart className="w-5 h-5 text-gray-400" />
              <span className="text-sm text-gray-400">{t('reports.likesReceived')}</span>
            </div>
            <p className="text-2xl font-bold">{report.likes_received}</p>
          </div>
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
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
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
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

      {/* Daily Activity Chart */}
      {report && report.daily_contributions && report.daily_contributions.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.dailyActivity')}</h2>
          <div className="flex items-end gap-1 h-40">
            {report.daily_contributions.map((day, index) => {
              const height = (day.contributions / maxContribution) * 100;
              return (
                <div
                  key={index}
                  className="flex-1 flex flex-col items-center gap-1 group"
                >
                  <div className="relative w-full">
                    <div
                      className="w-full bg-blue-500/80 rounded-t hover:bg-blue-400 transition-colors cursor-pointer"
                      style={{ height: `${Math.max(height, 4)}%`, minHeight: '4px' }}
                    />
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 hidden group-hover:block bg-gray-800 text-white text-xs px-2 py-1 rounded whitespace-nowrap">
                      {day.contributions} {t('reports.contributions')}
                      <br />
                      {day.posts} {t('reports.posts')}
                    </div>
                  </div>
                  <span className="text-[10px] text-gray-500">
                    {new Date(day.date).toLocaleDateString(undefined, {
                      weekday: 'short',
                    })}
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Top Languages */}
      {report && report.top_languages && report.top_languages.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.topLanguages')}</h2>
          <div className="space-y-3">
            {report.top_languages.map((lang, index) => {
              const maxBytes = report.top_languages[0].bytes;
              const width = (lang.bytes / maxBytes) * 100;
              return (
                <div key={index}>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="font-medium">{lang.language}</span>
                    <span className="text-gray-400">
                      {formatBytes(lang.bytes)} · {lang.repos} {t('reports.repos')}
                    </span>
                  </div>
                  <div className="h-2 bg-gray-800 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-500 rounded-full transition-all"
                      style={{ width: `${width}%` }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Goals Progress */}
      {report && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.goalsProgress')}</h2>
          <div className="flex items-center gap-4">
            <div className="relative w-24 h-24">
              <svg className="w-24 h-24 transform -rotate-90">
                <circle
                  cx="48"
                  cy="48"
                  r="40"
                  stroke="currentColor"
                  strokeWidth="8"
                  fill="none"
                  className="text-gray-800"
                />
                <circle
                  cx="48"
                  cy="48"
                  r="40"
                  stroke="currentColor"
                  strokeWidth="8"
                  fill="none"
                  strokeDasharray={`${report.goals_progress * 2.51} 251`}
                  className="text-blue-500"
                />
              </svg>
              <div className="absolute inset-0 flex items-center justify-center">
                <span className="text-lg font-bold">{report.goals_progress}%</span>
              </div>
            </div>
            <div>
              <p className="text-gray-400 text-sm">{t('reports.avgProgress')}</p>
              <p className="text-sm mt-1">
                <span className="font-medium">{report.goals_completed}</span>{' '}
                <span className="text-gray-400">{t('reports.goalsCompletedThisPeriod')}</span>
              </p>
            </div>
          </div>
        </div>
      )}
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
  const getTrendIcon = (d: number) => {
    if (d > 0) return '↑';
    if (d < 0) return '↓';
    return '—';
  };

  const getTrendColor = (d: number) => {
    if (d > 0) return 'text-green-400';
    if (d < 0) return 'text-red-400';
    return 'text-gray-400';
  };

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
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
