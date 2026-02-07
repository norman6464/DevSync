import { useState } from 'react';
import { getContributionRanking, getLanguageRanking, getAvailableLanguages } from '../api/rankings';
import type { RankingEntry } from '../types/ranking';
import { useAsyncData } from './useAsyncData';

const DEFAULT_LANGUAGES = [
  'JavaScript', 'TypeScript', 'Python', 'Java', 'Go', 'Rust', 'C++', 'C#',
  'Ruby', 'PHP', 'Swift', 'Kotlin', 'Scala', 'C', 'Shell', 'HTML', 'CSS'
];

export function useRankings() {
  const [tab, setTab] = useState<'contributions' | 'languages'>('contributions');
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

  const { data: rankings, loading } = useAsyncData(
    async () => {
      if (tab === 'contributions') {
        const { data } = await getContributionRanking(period);
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

  return {
    rankings,
    languages,
    loading,
    tab,
    setTab,
    period,
    setPeriod,
    language,
    setLanguage,
  };
}
