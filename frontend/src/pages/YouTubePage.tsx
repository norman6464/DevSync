import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Youtube, Sparkles, Search } from 'lucide-react';
import { useYoutube, useDebounce } from '../hooks';
import YouTubeVideoCard from '../components/youtube/YouTubeVideoCard';
import { Skeleton } from '../components/common/Skeleton';
import EmptyState from '../components/common/EmptyState';
import PageHeader from '../components/common/PageHeader';
import SearchInput from '../components/common/SearchInput';

export default function YouTubePage() {
  const { t } = useTranslation();
  const {
    tab, setTab, available,
    recommendations, recommendSkills, recommendLoading,
    searchQuery, setSearchQuery, searchResults, searching, cached,
    handleSearch,
  } = useYoutube();

  const debouncedQuery = useDebounce(searchQuery, 500);

  useEffect(() => {
    if (tab === 'search' && debouncedQuery) {
      handleSearch();
    }
  }, [debouncedQuery, tab, handleSearch]);

  if (!available) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        <EmptyState
          icon={Youtube}
          title={t('youtube.unavailable')}
        />
      </div>
    );
  }

  const videos = tab === 'recommend' ? recommendations : searchResults;
  const loading = tab === 'recommend' ? recommendLoading : searching;

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <PageHeader
        title={t('youtube.pageTitle')}
        subtitle={t('youtube.pageSubtitle')}
      />

      <div className="flex gap-2 mb-6">
        <button onClick={() => setTab('recommend')}
          className={`flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'recommend' ? 'bg-red-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}>
          <Sparkles className="w-4 h-4" />
          {t('youtube.recommendTab')}
        </button>
        <button onClick={() => setTab('search')}
          className={`flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'search' ? 'bg-red-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}>
          <Search className="w-4 h-4" />
          {t('youtube.searchTab')}
        </button>
      </div>

      {tab === 'recommend' && recommendSkills.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-4">
          <span className="text-sm text-gray-400">{t('youtube.basedOnSkills')}:</span>
          {recommendSkills.map(skill => (
            <span key={skill} className="px-2 py-0.5 bg-red-500/20 text-red-400 text-xs rounded-full">
              {skill}
            </span>
          ))}
        </div>
      )}

      {tab === 'search' && (
        <div className="mb-6">
          <SearchInput
            value={searchQuery}
            onChange={setSearchQuery}
            onSearch={handleSearch}
            placeholder={t('youtube.searchPlaceholder')}
            showButton
          />
          {cached && searchResults.length > 0 && (
            <p className="text-xs text-gray-500 mt-1">{t('youtube.cachedResult')}</p>
          )}
        </div>
      )}

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="bg-gray-800 rounded-md overflow-hidden border border-gray-700">
              <Skeleton className="w-full aspect-video" />
              <div className="p-4 space-y-2">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-3 w-1/2" />
              </div>
            </div>
          ))}
        </div>
      ) : videos.length === 0 ? (
        <EmptyState
          icon={Youtube}
          title={tab === 'recommend' ? t('youtube.noRecommendations') : t('youtube.noResults')}
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {videos.map(video => (
            <YouTubeVideoCard key={video.video_id} video={video} />
          ))}
        </div>
      )}
    </div>
  );
}
