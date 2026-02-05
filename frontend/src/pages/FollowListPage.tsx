import { useEffect, useState } from 'react';
import { useParams, Link, useLocation } from 'react-router-dom';
import { getUser, getFollowers, getFollowing } from '../api/users';
import { useAuthStore } from '../store/authStore';
import type { User } from '../types/user';
import Avatar from '../components/common/Avatar';
import FollowButton from '../components/profile/FollowButton';
import LoadingSpinner from '../components/common/LoadingSpinner';

type Tab = 'followers' | 'following';

export default function FollowListPage() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const currentUser = useAuthStore((s) => s.user);

  const initialTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
  const [tab, setTab] = useState<Tab>(initialTab);
  const [profileUser, setProfileUser] = useState<User | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    getUser(parseInt(id))
      .then(({ data }) => setProfileUser(data))
      .catch(() => {});
  }, [id]);

  useEffect(() => {
    const newTab: Tab = location.pathname.endsWith('/following') ? 'following' : 'followers';
    setTab(newTab);
  }, [location.pathname]);

  useEffect(() => {
    if (!id) return;
    const userId = parseInt(id);
    setLoading(true);
    const fetcher = tab === 'followers' ? getFollowers : getFollowing;
    fetcher(userId)
      .then(({ data }) => setUsers(data || []))
      .catch(() => setUsers([]))
      .finally(() => setLoading(false));
  }, [id, tab]);

  if (!profileUser && loading) return <div className="py-12"><LoadingSpinner /></div>;

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header with back link */}
      <div className="flex items-center gap-3">
        <Link
          to={`/profile/${id}`}
          className="p-1.5 rounded-lg hover:bg-gray-800 transition-colors text-gray-400 hover:text-white"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18" />
          </svg>
        </Link>
        <div>
          <h1 className="text-lg font-bold">{profileUser?.name}</h1>
          <p className="text-xs text-gray-500">@{profileUser?.github_username || profileUser?.name}</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-gray-800">
        <Link
          to={`/profile/${id}/followers`}
          className={`flex-1 text-center py-3 text-sm font-medium transition-colors relative ${
            tab === 'followers' ? 'text-white' : 'text-gray-500 hover:text-gray-300'
          }`}
        >
          Followers
          {tab === 'followers' && (
            <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-full" />
          )}
        </Link>
        <Link
          to={`/profile/${id}/following`}
          className={`flex-1 text-center py-3 text-sm font-medium transition-colors relative ${
            tab === 'following' ? 'text-white' : 'text-gray-500 hover:text-gray-300'
          }`}
        >
          Following
          {tab === 'following' && (
            <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-full" />
          )}
        </Link>
      </div>

      {/* User List */}
      {loading ? (
        <div className="py-12"><LoadingSpinner /></div>
      ) : users.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center text-gray-500 text-sm">
          {tab === 'followers' ? 'No followers yet' : 'Not following anyone yet'}
        </div>
      ) : (
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden divide-y divide-gray-800/50">
          {users.map((user) => (
            <div key={user.id} className="flex items-center gap-4 px-5 py-4 hover:bg-gray-800/40 transition-colors">
              <Link to={`/profile/${user.id}`} className="flex-shrink-0">
                <Avatar name={user.name} avatarUrl={user.avatar_url} size="sm" />
              </Link>

              <div className="flex-1 min-w-0">
                <Link
                  to={`/profile/${user.id}`}
                  className="font-medium text-sm hover:text-blue-400 transition-colors"
                >
                  {user.name}
                </Link>
                {user.github_username && (
                  <div className="flex items-center gap-1 mt-0.5 text-xs text-gray-500">
                    <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    @{user.github_username}
                  </div>
                )}
                {user.bio && (
                  <p className="text-xs text-gray-400 mt-1 line-clamp-1">{user.bio}</p>
                )}
              </div>

              <div className="flex items-center gap-2 flex-shrink-0">
                {/* Chat button - only show for other users */}
                {currentUser?.id !== user.id && (
                  <Link
                    to={`/chat/${user.id}`}
                    className="p-2 rounded-lg border border-gray-700 text-gray-400 hover:text-white hover:border-gray-500 hover:bg-gray-800 transition-colors"
                    title={`Chat with ${user.name}`}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" d="M8.625 12a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 0 1-2.555-.337A5.972 5.972 0 0 1 5.41 20.97a5.969 5.969 0 0 1-.474-.065 4.48 4.48 0 0 0 .978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25Z" />
                    </svg>
                  </Link>
                )}
                {/* Follow button - only show for other users */}
                {currentUser?.id !== user.id && (
                  <FollowButton userId={user.id} />
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
