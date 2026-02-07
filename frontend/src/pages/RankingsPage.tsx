import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useRankings } from '../hooks';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function RankingsPage() {
  const { t } = useTranslation();
  const {
    rankings, languages, loading,
    tab, setTab, period, setPeriod, language, setLanguage,
  } = useRankings();

  const medalColor = (index: number) => {
    if (index === 0) return 'text-yellow-400';
    if (index === 1) return 'text-gray-300';
    if (index === 2) return 'text-amber-600';
    return 'text-gray-600';
  };

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">{t('rankings.title')}</h1>

      {/* Tab & Period Controls */}
      <div className="flex items-center gap-3">
        <div className="flex bg-gray-900 border border-gray-800 rounded-lg p-1">
          <button
            onClick={() => setTab('contributions')}
            className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
              tab === 'contributions'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            <span className="flex items-center gap-1.5">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
              </svg>
              {t('rankings.contributions')}
            </span>
          </button>
          <button
            onClick={() => setTab('languages')}
            className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
              tab === 'languages'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            <span className="flex items-center gap-1.5">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M17.25 6.75 22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3-4.5 16.5" />
              </svg>
              {t('rankings.byLanguage')}
            </span>
          </button>
        </div>

        <div className="ml-auto flex bg-gray-900 border border-gray-800 rounded-lg p-1">
          <button
            onClick={() => setPeriod('weekly')}
            className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
              period === 'weekly'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('rankings.thisWeek')}
          </button>
          <button
            onClick={() => setPeriod('monthly')}
            className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
              period === 'monthly'
                ? 'bg-gray-700 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('rankings.thisMonth')}
          </button>
        </div>
      </div>

      {/* Language Filter */}
      {tab === 'languages' && (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
          <h3 className="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
            {t('rankings.selectLanguage')}
          </h3>
          <div className="flex flex-wrap gap-2">
            {languages.map((lang) => (
              <button
                key={lang}
                onClick={() => setLanguage(lang)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all border ${
                  language === lang
                    ? 'bg-purple-600 text-white border-purple-500 shadow-lg shadow-purple-500/20'
                    : 'bg-gray-800 text-gray-400 border-gray-700 hover:text-white hover:border-gray-500 hover:bg-gray-700'
                }`}
              >
                {lang}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Rankings Table */}
      {loading ? (
        <div className="py-12"><LoadingSpinner /></div>
      ) : rankings.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center text-gray-400 text-sm">
          {t('rankings.noData')}
        </div>
      ) : (
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
          <div className="px-6 py-3 border-b border-gray-800 grid grid-cols-[3rem_1fr_auto] gap-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
            <span className="text-center">{t('rankings.rank')}</span>
            <span>{t('rankings.developer')}</span>
            <span className="text-right">{tab === 'contributions' ? t('rankings.contributions') : t('rankings.total')}</span>
          </div>
          {rankings.map((entry, index) => (
            <Link
              key={entry.user_id}
              to={`/profile/${entry.user_id}`}
              className="grid grid-cols-[3rem_1fr_auto] gap-4 items-center px-6 py-3 hover:bg-gray-800/50 transition-colors border-b border-gray-800/50 last:border-b-0"
            >
              <span className={`text-center font-bold text-lg ${medalColor(index)}`}>
                {index < 3 ? (
                  <svg className="w-6 h-6 mx-auto" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 2l2.4 7.4H22l-6.2 4.5 2.4 7.4L12 16.8l-6.2 4.5 2.4-7.4L2 9.4h7.6z" />
                  </svg>
                ) : (
                  index + 1
                )}
              </span>
              <div className="flex items-center gap-3 min-w-0">
                <Avatar name={entry.name} avatarUrl={entry.avatar_url} size="sm" />
                <span className="font-medium text-sm truncate">{entry.name}</span>
              </div>
              <span className="text-green-400 font-semibold text-sm tabular-nums">
                {entry.score.toLocaleString()}
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
