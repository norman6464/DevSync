import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { TrendingUp, Heart, MessageCircle, Bookmark } from 'lucide-react';
import { useTrendingPosts, useTrendingResources } from '../../hooks';

export default function TrendingWidget() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<'posts' | 'resources'>('posts');
  const { posts, loading: postsLoading } = useTrendingPosts();
  const { resources, loading: resourcesLoading } = useTrendingResources();

  const loading = tab === 'posts' ? postsLoading : resourcesLoading;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
        <TrendingUp className="w-4 h-4 text-orange-400" />
        {t('recommendations.trending')}
      </h3>

      {/* Tabs */}
      <div className="flex items-center border-b border-gray-800 mb-3">
        <button
          onClick={() => setTab('posts')}
          className={`px-3 py-1.5 text-xs font-medium transition-colors relative ${
            tab === 'posts' ? 'text-white' : 'text-gray-400 hover:text-white'
          }`}
        >
          {t('recommendations.trendingPosts')}
          {tab === 'posts' && (
            <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-t" />
          )}
        </button>
        <button
          onClick={() => setTab('resources')}
          className={`px-3 py-1.5 text-xs font-medium transition-colors relative ${
            tab === 'resources' ? 'text-white' : 'text-gray-400 hover:text-white'
          }`}
        >
          {t('recommendations.trendingResources')}
          {tab === 'resources' && (
            <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-t" />
          )}
        </button>
      </div>

      {loading ? (
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-10 bg-gray-800 rounded animate-pulse" />
          ))}
        </div>
      ) : tab === 'posts' ? (
        posts.length === 0 ? (
          <p className="text-xs text-gray-500 text-center py-4">{t('recommendations.noTrending')}</p>
        ) : (
          <div className="space-y-1.5">
            {posts.slice(0, 5).map((post) => (
              <Link
                key={post.id}
                to={`/posts/${post.id}`}
                className="block p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
              >
                <p className="text-xs text-gray-200 truncate mb-1">{post.title}</p>
                <div className="flex items-center gap-3 text-[10px] text-gray-500">
                  <span className="flex items-center gap-1">
                    <Heart className="w-3 h-3" />
                    {post.like_count}
                  </span>
                  <span className="flex items-center gap-1">
                    <MessageCircle className="w-3 h-3" />
                    {post.comment_count}
                  </span>
                  <span className="truncate">{post.user?.name}</span>
                </div>
              </Link>
            ))}
          </div>
        )
      ) : resources.length === 0 ? (
        <p className="text-xs text-gray-500 text-center py-4">{t('recommendations.noTrending')}</p>
      ) : (
        <div className="space-y-1.5">
          {resources.slice(0, 5).map((resource) => (
            <Link
              key={resource.id}
              to={`/resources/${resource.id}`}
              className="block p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              <p className="text-xs text-gray-200 truncate mb-1">{resource.title}</p>
              <div className="flex items-center gap-3 text-[10px] text-gray-500">
                <span className="flex items-center gap-1">
                  <Heart className="w-3 h-3" />
                  {resource.like_count}
                </span>
                <span className="flex items-center gap-1">
                  <Bookmark className="w-3 h-3" />
                  {resource.save_count}
                </span>
                <span className="truncate">{resource.user?.name}</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
