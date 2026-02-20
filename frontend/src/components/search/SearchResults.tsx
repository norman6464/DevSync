import { useTranslation } from 'react-i18next';
import type { SearchTab } from './SearchTabs';
import UserSearchCard from './UserSearchCard';
import PostSearchCard from './PostSearchCard';
import CircleSearchCard from './CircleSearchCard';
import { UserCardSkeleton, PostCardSkeleton } from '../common/Skeleton';
import type { User } from '../../types/user';
import type { Post } from '../../types/post';
import type { StudyCircle } from '../../types/studyCircle';
import { emptyStateClass } from '../../constants/styles';

export function LoadingResults({ tab }: { tab: SearchTab }) {
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

export function SearchEmptyState({ message }: { message: string }) {
  return (
    <div className={emptyStateClass}>
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

function NoResults({ query }: { query: string }) {
  const { t } = useTranslation();
  return (
    <div className={emptyStateClass}>
      <p className="text-gray-400 mb-1">{t('search.noResults')}</p>
      <p className="text-gray-500 text-sm">"{query}"</p>
    </div>
  );
}

export function UserResults({ users, currentUserId, query }: { users: User[]; currentUserId?: number; query: string }) {
  if (users.length === 0) return <NoResults query={query} />;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {users.map((user) => (
        <UserSearchCard key={user.id} user={user} currentUserId={currentUserId} />
      ))}
    </div>
  );
}

export function PostResults({ posts, total, query }: { posts: Post[]; total: number; query: string }) {
  const { t } = useTranslation();

  if (posts.length === 0) return <NoResults query={query} />;

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

export function CircleResults({ circles, query }: { circles: StudyCircle[]; query: string }) {
  if (circles.length === 0) return <NoResults query={query} />;

  return (
    <div className="space-y-4">
      {circles.map((circle) => (
        <CircleSearchCard key={circle.id} circle={circle} />
      ))}
    </div>
  );
}
