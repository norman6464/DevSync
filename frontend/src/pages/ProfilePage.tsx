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

  if (loading) return <div className="py-12"><LoadingSpinner /></div>;
  if (!user) return <div className="text-center text-gray-400 py-12">User not found</div>;

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Profile Header */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-6">
        <div className="flex items-start gap-5">
          <Avatar name={user.name} avatarUrl={user.avatar_url} size="lg" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-3 flex-wrap">
              <h1 className="text-2xl font-bold">{user.name}</h1>
              {!isOwnProfile && <FollowButton userId={user.id} />}
            </div>
            {user.bio && <p className="text-gray-400 mt-1 text-sm">{user.bio}</p>}
            {user.github_username && (
              <div className="flex items-center gap-1.5 mt-2 text-sm text-gray-500">
                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                @{user.github_username}
              </div>
            )}
            <div className="flex gap-4 mt-3 text-sm">
              <span className="text-gray-400"><strong className="text-white">{followerCount}</strong> followers</span>
              <span className="text-gray-400"><strong className="text-white">{followingCount}</strong> following</span>
            </div>
          </div>
        </div>
      </div>

      {/* GitHub Data */}
      {user.github_connected && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-800">
              <h2 className="text-sm font-semibold">Contribution activity</h2>
            </div>
            <div className="p-6">
              <ContributionCalendar contributions={contributions} />
            </div>
          </div>

          {languages.length > 0 && (
            <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-800">
                <h2 className="text-sm font-semibold">Languages</h2>
              </div>
              <div className="p-6">
                <LanguageChart languages={languages} />
              </div>
            </div>
          )}
        </div>
      )}

      {/* Posts */}
      <div>
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">Posts</h2>
        {posts.length === 0 ? (
          <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center text-gray-400 text-sm">
            No posts yet
          </div>
        ) : (
          <div className="space-y-3">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
