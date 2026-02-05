import { useEffect, useState } from 'react';
import { followUser, unfollowUser, getFollowers } from '../../api/users';
import { useAuthStore } from '../../store/authStore';

interface FollowButtonProps {
  userId: number;
}

export default function FollowButton({ userId }: FollowButtonProps) {
  const currentUser = useAuthStore((s) => s.user);
  const [isFollowing, setIsFollowing] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!currentUser) return;
    getFollowers(userId)
      .then(({ data }) => {
        setIsFollowing((data || []).some((u) => u.id === currentUser.id));
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [userId, currentUser]);

  const handleToggle = async () => {
    setLoading(true);
    try {
      if (isFollowing) {
        await unfollowUser(userId);
        setIsFollowing(false);
      } else {
        await followUser(userId);
        setIsFollowing(true);
      }
    } catch {
      // handle error
    } finally {
      setLoading(false);
    }
  };

  if (loading) return null;

  return (
    <button
      onClick={handleToggle}
      className={`px-4 py-1.5 rounded-md text-sm font-medium transition-colors ${
        isFollowing
          ? 'bg-gray-700 text-gray-300 hover:bg-red-600 hover:text-white'
          : 'bg-blue-600 text-white hover:bg-blue-700'
      }`}
    >
      {isFollowing ? 'Following' : 'Follow'}
    </button>
  );
}
