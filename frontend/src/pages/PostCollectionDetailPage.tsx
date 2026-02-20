import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAsyncData } from '../hooks';
import { getPostCollection } from '../api/postCollections';
import { useCollectionPosts } from '../hooks/usePostCollections';
import type { PostCollection } from '../types/post';
import { PageLoader } from '../components/common';
import PostCard from '../components/posts/PostCard';
import { useAuthStore } from '../store/authStore';
import { Globe, Lock } from 'lucide-react';
import { emptyStateClass } from '../constants/styles';
import { formatDate } from '../utils/timeFormat';

export default function PostCollectionDetailPage() {
  const { t } = useTranslation();
  const { id } = useParams<{ id: string }>();
  const collectionId = Number(id);
  const currentUser = useAuthStore((s) => s.user);

  const { data: collection, loading: collectionLoading } = useAsyncData(
    async () => {
      const { data } = await getPostCollection(collectionId);
      return data;
    },
    { deps: [collectionId], enabled: !!collectionId }
  );

  const { items, loading: itemsLoading } = useCollectionPosts(collectionId);

  if (collectionLoading || itemsLoading) return <PageLoader />;
  if (!collection) return <div className="text-center text-gray-400 py-12">{t('errors.notFound')}</div>;

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold text-white">{collection.title}</h1>
          {collection.is_public ? (
            <Globe className="w-4 h-4 text-green-400" />
          ) : (
            <Lock className="w-4 h-4 text-gray-500" />
          )}
        </div>
        {collection.description && (
          <p className="text-gray-400 mt-2">{collection.description}</p>
        )}
        {(collection as PostCollection & { user?: { name: string; username: string } }).user && (
          <Link
            to={`/profile/${(collection as PostCollection & { user: { username: string } }).user.username}`}
            className="text-sm text-blue-400 hover:text-blue-300 mt-2 inline-block"
          >
            {(collection as PostCollection & { user: { name: string } }).user.name}
          </Link>
        )}
        <p className="text-xs text-gray-500 mt-2">
          {formatDate(collection.created_at)}
        </p>
      </div>

      <div>
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
          {t('collections.posts')} ({items.length})
        </h2>
        {items.length === 0 ? (
          <div className={`${emptyStateClass} text-gray-400 text-sm`}>
            {t('collections.noPosts')}
          </div>
        ) : (
          <div className="space-y-3">
            {items.map((item) => (
              <div key={item.id}>
                {item.note && (
                  <p className="text-xs text-gray-500 mb-1 ml-1">{item.note}</p>
                )}
                {item.post && (
                  <PostCard
                    post={item.post}
                    isOwner={currentUser?.id === item.post.user_id}
                  />
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
