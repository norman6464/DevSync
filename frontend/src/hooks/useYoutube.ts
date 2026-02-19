import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { YouTubeVideo } from '../types/youtube';
import {
  searchYouTubeVideos,
  getYouTubeRecommendations,
  getYouTubeStatus,
} from '../api/youtube';
import { useAsyncData } from './useAsyncData';

type TabType = 'recommend' | 'search';

export function useYoutube() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabType>('recommend');
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<YouTubeVideo[]>([]);
  const [searching, setSearching] = useState(false);
  const [cached, setCached] = useState(false);

  const { data: status } = useAsyncData(
    () => getYouTubeStatus(),
    { deps: [] }
  );

  const { data: recommendations, loading: recommendLoading } = useAsyncData(
    () => getYouTubeRecommendations(),
    { deps: [], enabled: status?.available === true }
  );

  const handleSearch = useCallback(async () => {
    if (!searchQuery.trim()) return;
    setSearching(true);
    try {
      const result = await searchYouTubeVideos(searchQuery.trim());
      setSearchResults(result.videos ?? []);
      setCached(result.cached);
    } catch {
      toast.error(t('youtube.searchFailed'));
    } finally {
      setSearching(false);
    }
  }, [searchQuery, t]);

  return {
    tab,
    setTab,
    available: status?.available ?? false,
    recommendations: recommendations?.videos ?? [],
    recommendSkills: recommendations?.skills ?? [],
    recommendLoading,
    searchQuery,
    setSearchQuery,
    searchResults,
    searching,
    cached,
    handleSearch,
  };
}
