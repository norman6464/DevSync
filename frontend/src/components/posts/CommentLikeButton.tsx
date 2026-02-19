import { useTranslation } from 'react-i18next';
import { useCommentLike } from '../../hooks/useCommentLike';

interface CommentLikeButtonProps {
  commentId: number;
  initialLiked: boolean;
  initialCount: number;
}

export default function CommentLikeButton({
  commentId,
  initialLiked,
  initialCount,
}: CommentLikeButtonProps) {
  const { t } = useTranslation();
  const { liked, likeCount, toggleLike, loading } = useCommentLike(
    commentId,
    initialLiked,
    initialCount,
  );

  return (
    <button
      onClick={toggleLike}
      disabled={loading}
      aria-pressed={liked}
      aria-label={liked ? t('post.unlike') : t('post.like')}
      className={`flex items-center gap-1 text-xs transition-colors disabled:opacity-50 ${
        liked
          ? 'text-pink-400 hover:text-pink-300'
          : 'text-gray-500 hover:text-pink-400'
      }`}
    >
      <svg
        aria-hidden="true"
        className="w-3.5 h-3.5"
        fill={liked ? 'currentColor' : 'none'}
        stroke="currentColor"
        strokeWidth="2"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z"
        />
      </svg>
      {likeCount > 0 && <span>{likeCount}</span>}
    </button>
  );
}
