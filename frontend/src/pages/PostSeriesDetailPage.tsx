import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAsyncData } from '../hooks';
import { getPostSeries } from '../api/postSeries';
import { useSeriesPosts } from '../hooks/usePostSeries';
import type { PostSeries } from '../types/post';
import { PageLoader } from '../components/common';
import PostCard from '../components/posts/PostCard';
import { useAuthStore } from '../store/authStore';

export default function PostSeriesDetailPage() {
  const { t } = useTranslation();
  const { id } = useParams<{ id: string }>();
  const seriesId = Number(id);
  const currentUser = useAuthStore((s) => s.user);

  const { data: series, loading: seriesLoading } = useAsyncData(
    async () => {
      const { data } = await getPostSeries(seriesId);
      return data;
    },
    { deps: [seriesId], enabled: !!seriesId }
  );

  const { items, loading: itemsLoading } = useSeriesPosts(seriesId);

  if (seriesLoading || itemsLoading) return <PageLoader />;
  if (!series) return <div className="text-center text-gray-400 py-12">{t('errors.notFound')}</div>;

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
        <h1 className="text-2xl font-bold text-white">{series.title}</h1>
        {series.description && (
          <p className="text-gray-400 mt-2">{series.description}</p>
        )}
        {(series as PostSeries & { user?: { name: string; username: string } }).user && (
          <Link
            to={`/profile/${(series as PostSeries & { user: { username: string } }).user.username}`}
            className="text-sm text-blue-400 hover:text-blue-300 mt-2 inline-block"
          >
            {(series as PostSeries & { user: { name: string } }).user.name}
          </Link>
        )}
        <p className="text-xs text-gray-500 mt-2">
          {new Date(series.created_at).toLocaleDateString()}
        </p>
      </div>

      <div>
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
          {t('series.posts')} ({items.length})
        </h2>
        {items.length === 0 ? (
          <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center text-gray-400 text-sm">
            {t('series.noPosts')}
          </div>
        ) : (
          <div className="space-y-3">
            {items.map((item, index) => (
              <div key={item.id} className="flex gap-3 items-start">
                <div className="flex-shrink-0 w-8 h-8 bg-gray-800 border border-gray-700 rounded-full flex items-center justify-center text-sm text-gray-400 mt-3">
                  {index + 1}
                </div>
                <div className="flex-1">
                  {item.post && (
                    <PostCard
                      post={item.post}
                      isOwner={currentUser?.id === item.post.user_id}
                    />
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
