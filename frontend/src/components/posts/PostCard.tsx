import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { cardDarkClass } from '../../constants/styles';
import { parseJsonArray } from '../../utils/json';
import PostCardContent from './PostCardContent';
import PostCardActions from './PostCardActions';

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
          <h3 className="text-base font-semibold group-hover:text-blue-400 transition-colors">{post.title}</h3>
          {post.estimated_read_time > 0 && (
            <p className="text-xs text-gray-500 mt-0.5">{t('post.readTime', { minutes: post.estimated_read_time })}</p>
          )}
        </Link>

        {isOwner && (
          <div className="flex gap-1 flex-shrink-0">
            <button
              onClick={onEdit}
              className="p-1.5 text-gray-400 hover:text-white transition-colors"
              aria-label={t('common.edit')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
              </svg>
            </button>
            <button
              onClick={onDelete}
              className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
              aria-label={t('common.delete')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
              </svg>
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
