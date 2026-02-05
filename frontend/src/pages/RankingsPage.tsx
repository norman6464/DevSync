import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getContributionRanking, getLanguageRanking, getAvailableLanguages } from '../api/rankings';
import type { RankingEntry } from '../types/ranking';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function RankingsPage() {
  const [tab, setTab] = useState<'contributions' | 'languages'>('contributions');
  const [period, setPeriod] = useState<'weekly' | 'monthly'>('weekly');
  const [language, setLanguage] = useState('');
  const [languages, setLanguages] = useState<string[]>([]);
  const [rankings, setRankings] = useState<RankingEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getAvailableLanguages()
      .then(({ data }) => {
        setLanguages(data || []);
        if (data && data.length > 0) setLanguage(data[0]);
      })
      .catch(() => {});
  }, []);

  useEffect(() => {
    const fetchRankings = async () => {
      setLoading(true);
      try {
        if (tab === 'contributions') {
          const { data } = await getContributionRanking(period);
          setRankings(data || []);
        } else if (language) {
          const { data } = await getLanguageRanking(language, period);
          setRankings(data || []);
        }
      } catch {
        setRankings([]);
      } finally {
        setLoading(false);
      }
    };
    fetchRankings();
  }, [tab, period, language]);

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Rankings</h1>

      <div className="flex flex-wrap gap-2">
        <button
          onClick={() => setTab('contributions')}
          className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
            tab === 'contributions' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          Contributions
        </button>
        <button
          onClick={() => setTab('languages')}
          className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
            tab === 'languages' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          By Language
        </button>

        <div className="ml-auto flex gap-2">
          <button
            onClick={() => setPeriod('weekly')}
            className={`px-3 py-1 rounded text-sm ${
              period === 'weekly' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400'
            }`}
          >
            Weekly
          </button>
          <button
            onClick={() => setPeriod('monthly')}
            className={`px-3 py-1 rounded text-sm ${
              period === 'monthly' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400'
            }`}
          >
            Monthly
          </button>
        </div>
      </div>

      {tab === 'languages' && (
        <div className="flex flex-wrap gap-2">
          {languages.map((lang) => (
            <button
              key={lang}
              onClick={() => setLanguage(lang)}
              className={`px-3 py-1 rounded text-sm ${
                language === lang ? 'bg-purple-600 text-white' : 'bg-gray-800 text-gray-400'
              }`}
            >
              {lang}
            </button>
          ))}
        </div>
      )}

      {loading ? (
        <LoadingSpinner />
      ) : rankings.length === 0 ? (
        <div className="bg-gray-900 rounded-lg p-8 text-center text-gray-400">
          No ranking data yet
        </div>
      ) : (
        <div className="bg-gray-900 rounded-lg divide-y divide-gray-800">
          {rankings.map((entry, index) => (
            <Link
              key={entry.user_id}
              to={`/profile/${entry.user_id}`}
              className="flex items-center gap-4 p-4 hover:bg-gray-800 transition-colors"
            >
              <span className="w-8 text-center font-bold text-lg text-gray-500">
                {index + 1}
              </span>
              <Avatar name={entry.name} avatarUrl={entry.avatar_url} />
              <span className="flex-1 font-medium">{entry.name}</span>
              <span className="text-blue-400 font-semibold">{entry.score.toLocaleString()}</span>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
