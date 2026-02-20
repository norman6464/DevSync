import { useState, useMemo } from 'react';
import { getContributionRanking, getLanguageRanking, getLevelRanking, getAvailableLanguages } from '../api/rankings';
import type { RankingEntry } from '../types/ranking';
import { useAsyncData } from './useAsyncData';

const DEFAULT_LANGUAGES = [
  'JavaScript', 'TypeScript', 'Python', 'Java', 'Go', 'Rust', 'C++', 'C#',
  'Ruby', 'PHP', 'Swift', 'Kotlin', 'Scala', 'C', 'Shell', 'HTML', 'CSS'
];

export function useRankings() {
  const [tab, setTab] = useState<'contributions' | 'languages' | 'level'>('contributions');
  const [period, setPeriod] = useState<'weekly' | 'monthly'>('weekly');
  const [language, setLanguage] = useState('JavaScript');

  const { data: languages } = useAsyncData(
    async () => {
      const { data } = await getAvailableLanguages();
      if (data && data.length > 0) {
        setLanguage(data[0]);
        return data;
      }
      return DEFAULT_LANGUAGES;
    },
    { initialData: DEFAULT_LANGUAGES }
  );

  const [searchQuery, setSearchQuery] = useState('');

  const { data: rankings, loading } = useAsyncData(
    async () => {
      if (tab === 'contributions') {
        const { data } = await getContributionRanking(period);
        return data || [];
      }
      if (tab === 'level') {
        const { data } = await getLevelRanking();
        return data || [];
      }
      if (language) {
        const { data } = await getLanguageRanking(language, period);
        return data || [];
      }
      return [];
    },
    { initialData: [] as RankingEntry[], deps: [tab, period, language] }
  );

  const filteredRankings = useMemo(() => {
    const q = searchQuery.toLowerCase().trim();
    if (!q) return rankings;
    return rankings.filter(
      (entry) => entry.name.toLowerCase().includes(q) || entry.username.toLowerCase().includes(q)
    );
  }, [rankings, searchQuery]);

  return {
    rankings: filteredRankings,
    languages,
    loading,
    tab,
    setTab,
    period,
    setPeriod,
    language,
    setLanguage,
    searchQuery,
    setSearchQuery,
  };
}
