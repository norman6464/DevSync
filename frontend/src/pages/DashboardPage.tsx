import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { getTimeline, getPosts, createPost } from '../api/posts';
import type { Post } from '../types/post';
import PostCard from '../components/posts/PostCard';
import PostForm from '../components/posts/PostForm';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function DashboardPage() {
  const user = useAuthStore((s) => s.user);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<'timeline' | 'all'>('timeline');

  const fetchPosts = async () => {
    setLoading(true);
    try {
      const { data } = tab === 'timeline' ? await getTimeline() : await getPosts();
      setPosts(data || []);
    } catch {
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, [tab]);

  const handleCreatePost = async (title: string, content: string) => {
    await createPost({ title, content });
    fetchPosts();
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div className="lg:col-span-2 space-y-6">
        <PostForm onSubmit={handleCreatePost} />

        <div className="flex gap-2">
          <button
            onClick={() => setTab('timeline')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              tab === 'timeline'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            Timeline
          </button>
          <button
            onClick={() => setTab('all')}
            className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
              tab === 'all'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            All Posts
          </button>
        </div>

        {loading ? (
          <LoadingSpinner />
        ) : posts.length === 0 ? (
          <div className="bg-gray-900 rounded-lg p-8 text-center text-gray-400">
            {tab === 'timeline'
              ? 'No posts from people you follow yet. Try following some users!'
              : 'No posts yet. Be the first to share!'}
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} onUpdate={fetchPosts} />
            ))}
          </div>
        )}
      </div>

      <div className="space-y-4">
        {user && (
          <div className="bg-gray-900 rounded-lg p-4">
            <h3 className="font-medium mb-2">Welcome, {user.name}!</h3>
            <div className="space-y-2 text-sm text-gray-400">
              <Link to={`/profile/${user.id}`} className="block hover:text-blue-400">
                View Profile
              </Link>
              <Link to="/settings" className="block hover:text-blue-400">
                Settings
              </Link>
              {!user.github_connected && (
                <Link to="/settings" className="block text-yellow-400 hover:text-yellow-300">
                  Connect GitHub to get started
                </Link>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
