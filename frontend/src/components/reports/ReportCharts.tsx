import { useTranslation } from 'react-i18next';
import { type ActivityReport } from '../../api/reports';

function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

interface ReportChartsProps {
  report: ActivityReport;
  maxContribution: number;
}

export default function ReportCharts({ report, maxContribution }: ReportChartsProps) {
  const { t } = useTranslation();

  return (
    <>
      {/* Daily Activity Chart */}
      {report.daily_contributions && report.daily_contributions.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.dailyActivity')}</h2>
          <div className="flex items-end gap-1 h-40">
            {report.daily_contributions.map((day) => {
              const height = (day.contributions / maxContribution) * 100;
              return (
                <div
                  key={day.date}
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
      {report.top_languages && report.top_languages.length > 0 && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <h2 className="text-lg font-semibold mb-4">{t('reports.topLanguages')}</h2>
          <div className="space-y-3">
            {report.top_languages.map((lang) => {
              const maxBytes = report.top_languages[0].bytes;
              const width = (lang.bytes / maxBytes) * 100;
              return (
                <div key={lang.language}>
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
      <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
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
    </>
  );
}
