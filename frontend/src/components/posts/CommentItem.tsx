import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { format } from 'date-fns';
import type { Comment } from '../../types/post';
import Avatar from '../common/Avatar';
import MentionText from '../common/MentionText';
import CommentLikeButton from './CommentLikeButton';

interface CommentItemProps {
  comment: Comment;
  onReply?: (commentId: number) => void;
  isReplying?: boolean;
  size?: 'sm' | 'xs';
}

export default function CommentItem({
  comment,
  onReply,
  isReplying = false,
  size = 'sm',
}: CommentItemProps) {
  const { t } = useTranslation();
  const isXs = size === 'xs';

  return (
    <div className={`flex gap-${isXs ? '2.5' : '3'}`}>
      <Link
        to={`/profile/${comment.user?.username || comment.user_id}`}
        className="flex-shrink-0"
      >
        <Avatar
          name={comment.user?.name || 'U'}
          avatarUrl={comment.user?.avatar_url}
          size={isXs ? 'xs' : 'sm'}
        />
      </Link>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2 flex-wrap">
          <Link
            to={`/profile/${comment.user?.username || comment.user_id}`}
            className={`font-medium ${isXs ? 'text-xs' : 'text-sm'} hover:text-blue-400 transition-colors`}
          >
            {comment.user?.name}
          </Link>
          <span className="text-xs text-gray-600">
            {format(new Date(comment.created_at), 'MMM d, yyyy · HH:mm')}
          </span>
        </div>
        <p className={`${isXs ? 'text-sm mt-0.5' : 'text-sm mt-1'} text-gray-300 leading-relaxed`}>
          <MentionText text={comment.content} />
        </p>
        <div className="flex items-center gap-3 mt-1.5">
          <CommentLikeButton
            commentId={comment.id}
            initialLiked={comment.liked ?? false}
            initialCount={comment.like_count ?? 0}
          />
          {onReply && (
            <button
              onClick={() => onReply(comment.id)}
              className={`text-xs transition-colors ${
                isReplying ? 'text-blue-400' : 'text-gray-500 hover:text-blue-400'
              }`}
            >
              {t('post.reply')}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
