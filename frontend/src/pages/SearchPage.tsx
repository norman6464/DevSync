import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import SearchBar from '../components/search/SearchBar';
import SearchTabs from '../components/search/SearchTabs';
import type { SearchTab } from '../components/search/SearchTabs';
import PostFilterPanel from '../components/search/PostFilterPanel';
import UserSearchCard from '../components/search/UserSearchCard';
import PostSearchCard from '../components/search/PostSearchCard';
import CircleSearchCard from '../components/search/CircleSearchCard';
import { useUserSearch, usePostSearch, useCircleSearch, useDebounce } from '../hooks';
import { useAuthStore } from '../store/authStore';
import { UserCardSkeleton, PostCardSkeleton } from '../components/common/Skeleton';
import type { User } from '../types/user';
import type { Post } from '../types/post';
import type { StudyCircle } from '../types/studyCircle';
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
          <EmptyState message={t('search.emptyInitial')} />
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

function LoadingResults({ tab }: { tab: SearchTab }) {
  if (tab === 'users') {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {[1, 2, 3, 4].map((i) => (
          <UserCardSkeleton key={i} />
        ))}
      </div>
    );
  }
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <PostCardSkeleton key={i} />
      ))}
    </div>
  );
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
      <svg
        className="w-16 h-16 mx-auto mb-4 text-gray-700"
        fill="none"
        stroke="currentColor"
        strokeWidth="1"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
        />
      </svg>
      <p className="text-gray-400">{message}</p>
    </div>
  );
}

function UserResults({ users, currentUserId, query }: { users: User[]; currentUserId?: number; query: string }) {
  const { t } = useTranslation();

  if (users.length === 0) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
        <p className="text-gray-400 mb-1">{t('search.noResults')}</p>
        <p className="text-gray-500 text-sm">"{query}"</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {users.map((user) => (
        <UserSearchCard key={user.id} user={user} currentUserId={currentUserId} />
      ))}
    </div>
  );
}

function PostResults({ posts, total, query }: { posts: Post[]; total: number; query: string }) {
  const { t } = useTranslation();

  if (posts.length === 0) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
        <p className="text-gray-400 mb-1">{t('search.noResults')}</p>
        <p className="text-gray-500 text-sm">"{query}"</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {total > posts.length && (
        <p className="text-sm text-gray-400">
          {t('search.totalResults', { count: total })}
        </p>
      )}
      {posts.map((post) => (
        <PostSearchCard key={post.id} post={post} />
      ))}
    </div>
  );
}

function CircleResults({ circles, query }: { circles: StudyCircle[]; query: string }) {
  const { t } = useTranslation();

  if (circles.length === 0) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
        <p className="text-gray-400 mb-1">{t('search.noResults')}</p>
        <p className="text-gray-500 text-sm">"{query}"</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {circles.map((circle) => (
        <CircleSearchCard key={circle.id} circle={circle} />
      ))}
    </div>
  );
}
