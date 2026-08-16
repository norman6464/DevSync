import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Code2, Pencil, Trash2 } from 'lucide-react';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { cardDarkClass, iconButtonClass, deleteIconButtonClass } from '../../constants/styles';
import { parseJsonArray } from '../../utils/json';
import PostCardContent from './PostCardContent';
import PostCardActions from './PostCardActions';
import { isWithinLast } from '../../utils/timeFormat';

interface PostCardProps {
  post: Post;
  isOwner?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
  onUpdate?: () => void;
}

export default function PostCard({ post, isOwner = false, onEdit, onDelete, onUpdate }: PostCardProps) {
  const { t } = useTranslation();

  const imageUrls = parseJsonArray(post.image_urls);
  const isNew = isWithinLast(post.created_at, 24 * 60 * 60 * 1000);

  return (
    <div className={cardDarkClass}>
      <div className="flex items-center gap-3 mb-3">
        <Link to={`/profile/${post.user?.username || post.user_id}`}>
          <Avatar name={post.user?.name || 'U'} avatarUrl={post.user?.avatar_url} size="sm" />
        </Link>
        <div className="min-w-0">
          <Link to={`/profile/${post.user?.username || post.user_id}`} className="font-medium text-sm hover:text-blue-400 transition-colors">
            {post.user?.name}
          </Link>
          <p className="text-xs text-gray-500">
            {format(new Date(post.created_at), 'MMM d, yyyy')}
          </p>
        </div>
      </div>

      <div className="flex items-start justify-between gap-2 mb-2">
        <Link to={`/posts/${post.id}`} className="flex-1 group">
          <div className="flex items-center gap-2">
            <h3 className="text-base font-semibold group-hover:text-blue-400 transition-colors">{post.title}</h3>
            {isNew && (
              <span className="text-[10px] font-bold px-1.5 py-0.5 bg-green-500/20 text-green-400 rounded shrink-0">
                {t('post.newBadge')}
              </span>
            )}
          </div>
          <div className="flex items-center gap-2 mt-0.5">
            {post.estimated_read_time > 0 && (
              <span className="text-xs text-gray-500">{t('post.readTime', { minutes: post.estimated_read_time })}</span>
            )}
            {post.code_snippets && post.code_snippets.length > 0 && (
              <span className="inline-flex items-center gap-0.5 text-[10px] px-1.5 py-0.5 bg-purple-500/15 text-purple-400 rounded">
                <Code2 className="w-3 h-3" />
                {post.code_snippets.length}
              </span>
            )}
          </div>
        </Link>

        {isOwner && (
          <div className="flex gap-1 flex-shrink-0">
            <button
              onClick={onEdit}
              className={iconButtonClass}
              aria-label={t('common.edit')}
            >
              <Pencil className="w-4 h-4" />
            </button>
            <button
              onClick={onDelete}
              className={deleteIconButtonClass}
              aria-label={t('common.delete')}
            >
              <Trash2 className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>

      {post.tags && post.tags.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mb-2">
          {post.tags.slice(0, 3).map((tag) => (
            <span key={tag} className="text-xs px-2 py-0.5 bg-blue-500/10 text-blue-400 rounded">
              #{tag}
            </span>
          ))}
          {post.tags.length > 3 && (
            <span className="text-xs text-gray-500">+{post.tags.length - 3}</span>
          )}
        </div>
      )}

      <PostCardContent post={post} imageUrls={imageUrls} />

      <PostCardActions
        postId={post.id}
        initialLiked={post.liked || false}
        initialLikeCount={post.like_count}
        initialBookmarked={post.bookmarked || false}
        bookmarkCount={post.bookmark_count}
        commentCount={post.comment_count}
        viewCount={post.view_count}
        onUpdate={onUpdate}
      />
    </div>
  );
}
