import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';

interface PostSearchCardProps {
  post: Post;
}

export default function PostSearchCard({ post }: PostSearchCardProps) {
  const { t } = useTranslation();

  return (
    <Link
      to={`/posts/${post.id}`}
      className="block bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors"
    >
      <div className="flex items-start gap-3">
        <Avatar name={post.user?.name || 'User'} avatarUrl={post.user?.avatar_url} size="sm" />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 text-xs text-gray-500 mb-1">
            <span className="font-medium text-gray-300">{post.user?.name}</span>
            <span>•</span>
            <span>{new Date(post.created_at).toLocaleDateString()}</span>
          </div>
          <h3 className="font-semibold text-white mb-1">{post.title}</h3>
          <p className="text-sm text-gray-400 line-clamp-2">{post.content}</p>
          <div className="flex items-center gap-4 mt-3 text-xs text-gray-500">
            <span>{t('search.likesCount', { count: post.like_count || 0 })}</span>
            <span>{t('search.commentsCount', { count: post.comment_count || 0 })}</span>
          </div>
        </div>
      </div>
    </Link>
  );
}
