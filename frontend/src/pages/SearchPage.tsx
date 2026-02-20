import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import SearchBar from '../components/search/SearchBar';
import SearchTabs from '../components/search/SearchTabs';
import type { SearchTab } from '../components/search/SearchTabs';
import PostFilterPanel from '../components/search/PostFilterPanel';
import { LoadingResults, SearchEmptyState, UserResults, PostResults, CircleResults } from '../components/search/SearchResults';
import { useUserSearch, usePostSearch, useCircleSearch, useDebounce } from '../hooks';
import { useAuthStore } from '../store/authStore';
import type { PostSearchFilters } from '../api/posts';

export default function NewSearchPage() {
  const { t } = useTranslation();
  const currentUser = useAuthStore((s) => s.user);
  const [activeTab, setActiveTab] = useState<SearchTab>('users');
  const [globalQuery, setGlobalQuery] = useState('');
  const [showFilters, setShowFilters] = useState(false);

  // デバウンス処理（300ms）
  const debouncedQuery = useDebounce(globalQuery, 300);

  const userSearch = useUserSearch();
  const postSearch = usePostSearch();
  const circleSearch = useCircleSearch();

  const handleGlobalQueryChange = useCallback((value: string) => {
    setGlobalQuery(value);
  }, []);

  const handleGlobalSearch = useCallback(() => {
    userSearch.handleSearch(globalQuery);
    postSearch.handleSearch(globalQuery);
    circleSearch.handleSearch(globalQuery);
  }, [globalQuery, userSearch, postSearch, circleSearch]);

  // デバウンスされたクエリで自動検索
  useEffect(() => {
    if (debouncedQuery) {
      userSearch.setQuery(debouncedQuery);
      postSearch.setQuery(debouncedQuery);
      circleSearch.setQuery(debouncedQuery);
      userSearch.handleSearch(debouncedQuery);
      postSearch.handleSearch(debouncedQuery);
      circleSearch.handleSearch(debouncedQuery);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedQuery]);

  const handleFiltersChange = useCallback((filters: PostSearchFilters) => {
    postSearch.setFilters(filters);
    if (debouncedQuery) {
      postSearch.handleSearch(debouncedQuery, filters);
    }
  }, [debouncedQuery, postSearch]);

  const counts = useMemo(() => ({
    users: userSearch.filteredUsers?.length || 0,
    posts: postSearch.results.length,
    circles: circleSearch.results.length,
  }), [userSearch.filteredUsers, postSearch.results, circleSearch.results]);

  const isLoading = userSearch.loading || postSearch.loading || circleSearch.loading;
  const hasSearched = userSearch.searched || postSearch.searched || circleSearch.searched;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold mb-2">{t('search.title')}</h1>
        <p className="text-sm text-gray-400">{t('search.description')}</p>
      </div>

      {/* Search Bar */}
      <div className="flex gap-2">
        <div className="flex-1">
          <SearchBar
            value={globalQuery}
            onChange={handleGlobalQueryChange}
            onSearch={handleGlobalSearch}
            placeholder={t('search.placeholder')}
          />
        </div>
        {activeTab === 'posts' && (
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`px-4 py-2 rounded-lg border text-sm font-medium transition-colors ${
              showFilters
                ? 'border-blue-500 text-blue-400 bg-blue-500/10'
                : 'border-gray-700 text-gray-400 hover:border-gray-600'
            }`}
          >
            {t('search.filters')}
          </button>
        )}
      </div>

      {/* Filter Panel (posts tab only) */}
      {activeTab === 'posts' && showFilters && (
        <PostFilterPanel
          filters={postSearch.filters}
          onFiltersChange={handleFiltersChange}
        />
      )}

      {/* Tabs */}
      <SearchTabs activeTab={activeTab} onTabChange={setActiveTab} counts={counts} />

      {/* Results */}
      <div>
        {isLoading ? (
          <LoadingResults tab={activeTab} />
        ) : !hasSearched ? (
          <SearchEmptyState message={t('search.emptyInitial')} />
        ) : activeTab === 'users' ? (
          <UserResults users={userSearch.filteredUsers} currentUserId={currentUser?.id} query={globalQuery} />
        ) : activeTab === 'posts' ? (
          <PostResults posts={postSearch.results} total={postSearch.total} query={globalQuery} />
        ) : (
          <CircleResults circles={circleSearch.results} query={globalQuery} />
        )}
      </div>
    </div>
  );
}
