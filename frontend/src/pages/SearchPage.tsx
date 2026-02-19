import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import SearchBar from '../components/search/SearchBar';
import SearchTabs from '../components/search/SearchTabs';
import type { SearchTab } from '../components/search/SearchTabs';
import { useUserSearch, usePostSearch, useCircleSearch, useDebounce } from '../hooks';
import { useAuthStore } from '../store/authStore';
import Avatar from '../components/common/Avatar';
import FollowButton from '../components/profile/FollowButton';
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

// フィルターパネル（投稿タブ用）
function PostFilterPanel({
  filters,
  onFiltersChange,
}: {
  filters: PostSearchFilters;
  onFiltersChange: (f: PostSearchFilters) => void;
}) {
  const { t } = useTranslation();
  const [tagInput, setTagInput] = useState('');

  const handleSortChange = (sortBy: PostSearchFilters['sortBy']) => {
    onFiltersChange({ ...filters, sortBy });
  };

  const handleAddTag = () => {
    const tag = tagInput.trim();
    if (tag && !filters.tags?.includes(tag)) {
      onFiltersChange({ ...filters, tags: [...(filters.tags || []), tag] });
      setTagInput('');
    }
  };

  const handleRemoveTag = (tag: string) => {
    onFiltersChange({ ...filters, tags: filters.tags?.filter((t) => t !== tag) });
  };

  const handleDateFromChange = (value: string) => {
    onFiltersChange({ ...filters, dateFrom: value || undefined });
  };

  const handleDateToChange = (value: string) => {
    onFiltersChange({ ...filters, dateTo: value || undefined });
  };

  const sortOptions: { value: PostSearchFilters['sortBy']; label: string }[] = [
    { value: 'latest', label: t('search.sortLatest') },
    { value: 'popular', label: t('search.sortPopular') },
    { value: 'views', label: t('search.sortViews') },
  ];

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg p-4 space-y-4">
      {/* ソート順 */}
      <div>
        <label className="block text-xs font-medium text-gray-400 mb-2">{t('search.sortBy')}</label>
        <div className="flex gap-2">
          {sortOptions.map((opt) => (
            <button
              key={opt.value}
              onClick={() => handleSortChange(opt.value)}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                filters.sortBy === opt.value
                  ? 'bg-blue-500/20 text-blue-400 border border-blue-500/50'
                  : 'bg-gray-800 text-gray-400 border border-gray-700 hover:border-gray-600'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      {/* タグフィルター */}
      <div>
        <label className="block text-xs font-medium text-gray-400 mb-2">{t('search.tagFilter')}</label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleAddTag()}
            placeholder={t('search.tagPlaceholder')}
            className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <button
            onClick={handleAddTag}
            className="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-300 hover:border-gray-600 transition-colors"
          >
            {t('common.add')}
          </button>
        </div>
        {filters.tags && filters.tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {filters.tags.map((tag) => (
              <span
                key={tag}
                className="inline-flex items-center gap-1 px-2 py-1 bg-blue-500/10 text-blue-400 border border-blue-500/30 rounded-full text-xs"
              >
                #{tag}
                <button
                  onClick={() => handleRemoveTag(tag)}
                  className="hover:text-blue-300 ml-0.5"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
      </div>

      {/* 日付範囲フィルター */}
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-1">{t('search.dateFrom')}</label>
          <input
            type="date"
            value={filters.dateFrom || ''}
            onChange={(e) => handleDateFromChange(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 focus:outline-none focus:border-blue-500"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-1">{t('search.dateTo')}</label>
          <input
            type="date"
            value={filters.dateTo || ''}
            onChange={(e) => handleDateToChange(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 focus:outline-none focus:border-blue-500"
          />
        </div>
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
        <UserCard key={user.id} user={user} currentUserId={currentUserId} t={t} />
      ))}
    </div>
  );
}

function UserCard({ user, currentUserId, t }: { user: User; currentUserId?: number; t: (key: string) => string }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors">
      <div className="flex items-start gap-4">
        <Link to={`/profile/${user.username}`}>
          <Avatar name={user.name} avatarUrl={user.avatar_url} />
        </Link>
        <div className="flex-1 min-w-0">
          <Link to={`/profile/${user.username}`} className="font-semibold text-sm hover:text-blue-400 transition-colors">
            {user.name}
          </Link>
          {user.bio && <p className="text-xs text-gray-400 mt-1 line-clamp-2">{user.bio}</p>}
        </div>
      </div>
      <div className="flex items-center gap-2 mt-4 pt-3 border-t border-gray-800/50">
        <Link
          to={`/profile/${user.username}`}
          className="flex-1 px-3 py-1.5 rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-300 text-xs font-medium text-center transition-colors"
        >
          {t('search.viewProfile')}
        </Link>
        {currentUserId && currentUserId !== user.id && (
          <div className="flex-shrink-0">
            <FollowButton userId={user.id} />
          </div>
        )}
      </div>
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
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}

function PostCard({ post }: { post: Post }) {
  const { t } = useTranslation();
  return (
    <Link
      to={`/posts/${post.id}`}
      className="block bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors"
    >
      <div className="flex items-start gap-3">
        <Avatar name={post.user?.name || 'User'} avatarUrl={post.user?.avatar_url} size="sm" />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 text-xs text-gray-500 mb-1">
            <span className="font-medium text-gray-300">{post.user?.name}</span>
            <span>•</span>
            <span>{new Date(post.created_at).toLocaleDateString()}</span>
          </div>
          <h3 className="font-semibold text-white mb-1">{post.title}</h3>
          <p className="text-sm text-gray-400 line-clamp-2">{post.content}</p>
          <div className="flex items-center gap-4 mt-3 text-xs text-gray-500">
            <span>{t('search.likesCount', { count: post.like_count || 0 })}</span>
            <span>{t('search.commentsCount', { count: post.comment_count || 0 })}</span>
          </div>
        </div>
      </div>
    </Link>
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
        <CircleCard key={circle.id} circle={circle} />
      ))}
    </div>
  );
}

function CircleCard({ circle }: { circle: StudyCircle }) {
  const { t } = useTranslation();

  return (
    <Link
      to={`/study-circles/${circle.id}`}
      className="block bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <h3 className="font-semibold text-white mb-1">{circle.name}</h3>
          <p className="text-sm text-blue-400 mb-2">{circle.topic}</p>
          {circle.description && <p className="text-sm text-gray-400 line-clamp-2">{circle.description}</p>}
        </div>
        <div className="text-right flex-shrink-0">
          <div className="text-sm text-gray-400">
            {circle.member_count || 0} / {circle.max_members || '∞'}
          </div>
          <div className="text-xs text-gray-500 mt-1">{t('studyCircles.members')}</div>
        </div>
      </div>
    </Link>
  );
}
