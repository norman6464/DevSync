import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users, UserPlus, Check } from 'lucide-react';
import { useRecommendedUsers } from '../../hooks';
import { followUser } from '../../api/users';
import Avatar from '../common/Avatar';

export default function RecommendedUsersWidget() {
  const { t } = useTranslation();
  const { users, loading } = useRecommendedUsers();
  const [followedIds, setFollowedIds] = useState<Set<number>>(new Set());
  const [followingId, setFollowingId] = useState<number | null>(null);

  const handleFollow = async (userId: number) => {
    setFollowingId(userId);
    try {
      await followUser(userId);
      setFollowedIds((prev) => new Set(prev).add(userId));
    } catch (e) {
      console.warn('Failed to follow user:', e);
    } finally {
      setFollowingId(null);
    }
  };

  if (loading) {
    return (
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <div className="h-5 bg-gray-800 rounded animate-pulse w-1/2 mb-3" />
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-3">
              <div className="w-8 h-8 bg-gray-800 rounded-full animate-pulse" />
              <div className="flex-1">
                <div className="h-3 bg-gray-800 rounded animate-pulse w-2/3 mb-1" />
                <div className="h-2 bg-gray-800 rounded animate-pulse w-1/2" />
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (!users || users.length === 0) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="flex items-center gap-2 text-sm font-medium text-white">
          <Users aria-hidden="true" className="w-4 h-4 text-purple-400" />
          {t('recommendations.recommendedUsers')}
        </h3>
        <Link to="/search" className="text-xs text-gray-400 hover:text-blue-400 transition-colors">
          {t('recommendations.viewAll')}
        </Link>
      </div>

      <div className="space-y-2">
        {users.slice(0, 5).map((rec) => {
          const isFollowed = followedIds.has(rec.user.id);
          const isFollowing = followingId === rec.user.id;

          return (
            <div key={rec.user.id} className="flex items-center gap-2.5 p-1.5 rounded-lg hover:bg-gray-800/50 transition-colors">
              <Link to={`/profile/${rec.user.id}`}>
                <Avatar name={rec.user.name} avatarUrl={rec.user.avatar_url} size="sm" />
              </Link>
              <div className="flex-1 min-w-0">
                <Link
                  to={`/profile/${rec.user.id}`}
                  className="text-xs font-medium text-gray-200 hover:text-white block truncate"
                >
                  {rec.user.name}
                </Link>
                <div className="flex flex-wrap gap-1 mt-0.5">
                  {rec.common_skills.slice(0, 3).map((skill) => (
                    <span
                      key={skill}
                      className="text-[10px] px-1.5 py-0.5 bg-purple-500/10 text-purple-400 rounded"
                    >
                      {skill}
                    </span>
                  ))}
                  {rec.common_skills.length > 3 && (
                    <span className="text-[10px] text-gray-500">
                      +{rec.common_skills.length - 3}
                    </span>
                  )}
                </div>
              </div>
              <button
                onClick={() => !isFollowed && !isFollowing && handleFollow(rec.user.id)}
                disabled={isFollowed || isFollowing}
                aria-label={isFollowed ? t('recommendations.following') : t('recommendations.follow')}
                className={`shrink-0 flex items-center gap-1 px-2 py-1 rounded text-[10px] font-medium transition-colors ${
                  isFollowed
                    ? 'bg-gray-800 text-gray-500 cursor-default'
                    : 'bg-blue-600 hover:bg-blue-500 text-white'
                }`}
              >
                {isFollowed ? (
                  <>
                    <Check aria-hidden="true" className="w-3 h-3" />
                  </>
                ) : (
                  <>
                    <UserPlus aria-hidden="true" className="w-3 h-3" />
                    {t('recommendations.follow')}
                  </>
                )}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
