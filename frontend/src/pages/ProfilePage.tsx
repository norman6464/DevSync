import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getUser, getFollowers, getFollowing } from '../api/users';
import { getUserPosts } from '../api/posts';
import { getContributions, getLanguages } from '../api/github';
import { useAuthStore } from '../store/authStore';
import type { User } from '../types/user';
import type { Post } from '../types/post';
import type { GitHubContribution, GitHubLanguageStat } from '../types/github';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';
import FollowButton from '../components/profile/FollowButton';
import ContributionCalendar from '../components/profile/ContributionCalendar';
import LanguageChart from '../components/profile/LanguageChart';
import PostCard from '../components/posts/PostCard';

export default function ProfilePage() {
  const { id } = useParams<{ id: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const [user, setUser] = useState<User | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [contributions, setContributions] = useState<GitHubContribution[]>([]);
  const [languages, setLanguages] = useState<GitHubLanguageStat[]>([]);
  const [followerCount, setFollowerCount] = useState(0);
  const [followingCount, setFollowingCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    const userId = parseInt(id);

    const fetchData = async () => {
      setLoading(true);
      try {
        const [userRes, postsRes, followersRes, followingRes] = await Promise.all([
          getUser(userId),
          getUserPosts(userId),
          getFollowers(userId),
          getFollowing(userId),
        ]);
        setUser(userRes.data);
        setPosts(postsRes.data || []);
        setFollowerCount((followersRes.data || []).length);
        setFollowingCount((followingRes.data || []).length);

        if (userRes.data.github_connected) {
          const [contribRes, langRes] = await Promise.all([
            getContributions(userId),
            getLanguages(userId),
          ]);
          setContributions(contribRes.data || []);
          setLanguages(langRes.data || []);
        }
      } catch {
        // handle error
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading) return <LoadingSpinner />;
  if (!user) return <div className="text-center text-gray-400 py-12">User not found</div>;

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="space-y-6">
      <div className="bg-gray-900 rounded-lg p-6">
        <div className="flex items-start gap-4">
          <Avatar name={user.name} avatarUrl={user.avatar_url} size="lg" />
          <div className="flex-1">
            <div className="flex items-center gap-4">
              <h1 className="text-2xl font-bold">{user.name}</h1>
              {!isOwnProfile && <FollowButton userId={user.id} />}
            </div>
            {user.bio && <p className="text-gray-400 mt-1">{user.bio}</p>}
            {user.github_username && (
              <p className="text-sm text-gray-500 mt-1">
                GitHub: @{user.github_username}
              </p>
            )}
            <div className="flex gap-4 mt-3 text-sm text-gray-400">
              <span><strong className="text-white">{followerCount}</strong> followers</span>
              <span><strong className="text-white">{followingCount}</strong> following</span>
            </div>
          </div>
        </div>
      </div>

      {user.github_connected && (
        <>
          <div className="bg-gray-900 rounded-lg p-6">
            <h2 className="text-lg font-semibold mb-4">Contributions</h2>
            <ContributionCalendar contributions={contributions} />
          </div>

          {languages.length > 0 && (
            <div className="bg-gray-900 rounded-lg p-6">
              <h2 className="text-lg font-semibold mb-4">Languages</h2>
              <LanguageChart languages={languages} />
            </div>
          )}
        </>
      )}

      <div>
        <h2 className="text-lg font-semibold mb-4">Posts</h2>
        {posts.length === 0 ? (
          <div className="bg-gray-900 rounded-lg p-8 text-center text-gray-400">
            No posts yet
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
