import { useState, useEffect } from 'react';
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

export default function NewSearchPage() {
  const { t } = useTranslation();
  const currentUser = useAuthStore((s) => s.user);
  const [activeTab, setActiveTab] = useState<SearchTab>('users');
  const [globalQuery, setGlobalQuery] = useState('');

  // デバウンス処理（300ms）
  const debouncedQuery = useDebounce(globalQuery, 300);

  const userSearch = useUserSearch();
  const postSearch = usePostSearch();
  const circleSearch = useCircleSearch();

  const handleGlobalQueryChange = (value: string) => {
    setGlobalQuery(value);
  };

  const handleGlobalSearch = () => {
    userSearch.handleSearch();
    postSearch.handleSearch();
    circleSearch.handleSearch();
  };

  // デバウンスされたクエリで自動検索
  useEffect(() => {
    if (debouncedQuery) {
      userSearch.setQuery(debouncedQuery);
      postSearch.setQuery(debouncedQuery);
      circleSearch.setQuery(debouncedQuery);
      handleGlobalSearch();
    }
  }, [debouncedQuery]);

  const counts = {
    users: userSearch.filteredUsers?.length || 0,
    posts: postSearch.results.length,
    circles: circleSearch.results.length,
  };

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
      <SearchBar
        value={globalQuery}
        onChange={handleGlobalQueryChange}
        onSearch={handleGlobalSearch}
        placeholder={t('search.placeholder')}
      />

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
          <PostResults posts={postSearch.results} query={globalQuery} />
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
        <UserCard key={user.id} user={user} currentUserId={currentUserId} t={t} />
      ))}
    </div>
  );
}

function UserCard({ user, currentUserId, t }: { user: User; currentUserId?: number; t: (key: string) => string }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors">
      <div className="flex items-start gap-4">
        <Link to={`/profile/${user.id}`}>
          <Avatar name={user.name} avatarUrl={user.avatar_url} />
        </Link>
        <div className="flex-1 min-w-0">
          <Link to={`/profile/${user.id}`} className="font-semibold text-sm hover:text-blue-400 transition-colors">
            {user.name}
          </Link>
          {user.bio && <p className="text-xs text-gray-400 mt-1 line-clamp-2">{user.bio}</p>}
        </div>
      </div>
      <div className="flex items-center gap-2 mt-4 pt-3 border-t border-gray-800/50">
        <Link
          to={`/profile/${user.id}`}
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

function PostResults({ posts, query }: { posts: Post[]; query: string }) {
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
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </div>
  );
}

function PostCard({ post }: { post: Post }) {
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
            <span>{post.like_count || 0} likes</span>
            <span>{post.comment_count || 0} comments</span>
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
