import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { User } from '../../types/user';
import Avatar from '../common/Avatar';
import FollowButton from '../profile/FollowButton';

interface UserSearchCardProps {
  user: User;
  currentUserId?: number;
}

export default function UserSearchCard({ user, currentUserId }: UserSearchCardProps) {
  const { t } = useTranslation();

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
