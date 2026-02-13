import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { FileEdit, Eye, Trash2, Send } from 'lucide-react';
import { useAsyncData } from '../hooks/useAsyncData';
import { getDrafts, publishPost, deletePost } from '../api/posts';
import { PostCardSkeleton } from '../components/common/Skeleton';
import { formatDistanceToNow } from '../utils/timeFormat';

export default function DraftsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [refreshKey, setRefreshKey] = useState(0);

  const { data: drafts, loading } = useAsyncData(
    async () => {
      const { data } = await getDrafts();
      return data || [];
    },
    { initialData: [], deps: [refreshKey] }
  );

  const handlePublish = async (id: number) => {
    if (!confirm(t('post.confirmPublish'))) return;
    try {
      await publishPost(id);
      toast.success(t('post.published'));
      setRefreshKey((k) => k + 1);
    } catch {
      toast.error(t('post.publishFailed'));
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm(t('post.confirmDelete'))) return;
    try {
      await deletePost(id);
      toast.success(t('post.deleted'));
      setRefreshKey((k) => k + 1);
    } catch {
      toast.error(t('post.deleteFailed'));
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <FileEdit className="w-6 h-6 text-orange-500" />
        <h1 className="text-2xl font-bold text-white">{t('post.drafts')}</h1>
      </div>

      {loading ? (
        <div className="space-y-3">
          <PostCardSkeleton />
          <PostCardSkeleton />
        </div>
      ) : drafts.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
          <FileEdit className="w-16 h-16 mx-auto mb-4 text-gray-700" />
          <p className="text-gray-400">{t('post.noDrafts')}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {drafts.map((draft) => (
            <div
              key={draft.id}
              className="bg-gray-900 border border-gray-800 rounded-md p-4 hover:border-gray-700 transition-colors"
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1 min-w-0">
                  <h3 className="text-lg font-semibold text-white truncate">{draft.title}</h3>
                  <p className="text-sm text-gray-400 line-clamp-2 mt-1">{draft.content}</p>
                </div>
                <span className="ml-3 px-2 py-1 text-xs font-medium bg-yellow-500/10 text-yellow-400 rounded shrink-0">
                  {t('post.draft')}
                </span>
              </div>

              <div className="flex items-center justify-between mt-4 pt-3 border-t border-gray-800">
                <span className="text-xs text-gray-500">
                  {t('post.lastEdited')}: {formatDistanceToNow(draft.updated_at)}
                </span>

                <div className="flex items-center gap-2">
                  <button
                    onClick={() => navigate(`/posts/${draft.id}`)}
                    className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
                    title={t('common.view')}
                  >
                    <Eye className="w-4 h-4" />
                    {t('common.view')}
                  </button>
                  <button
                    onClick={() => handlePublish(draft.id)}
                    className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-blue-400 hover:text-blue-300 hover:bg-blue-500/10 rounded-lg transition-colors"
                    title={t('post.publish')}
                  >
                    <Send className="w-4 h-4" />
                    {t('post.publish')}
                  </button>
                  <button
                    onClick={() => handleDelete(draft.id)}
                    className="flex items-center gap-1.5 px-3 py-1.5 text-sm text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-colors"
                    title={t('common.delete')}
                  >
                    <Trash2 className="w-4 h-4" />
                    {t('common.delete')}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
