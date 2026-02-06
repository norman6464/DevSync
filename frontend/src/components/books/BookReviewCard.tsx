import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { likeBookReview, unlikeBookReview, type BookReview } from '../../api/bookReviews';
import Avatar from '../common/Avatar';
import toast from 'react-hot-toast';

interface BookReviewCardProps {
  review: BookReview;
  onDelete?: () => void;
  showUser?: boolean;
  isOwner?: boolean;
}

export default function BookReviewCard({ review, onDelete, showUser = true, isOwner = false }: BookReviewCardProps) {
  const { t } = useTranslation();
  const [liked, setLiked] = useState(review.has_liked ?? false);
  const [likeCount, setLikeCount] = useState(review.like_count);

  const handleLike = async () => {
    try {
      if (liked) {
        await unlikeBookReview(review.id);
        setLikeCount((c) => c - 1);
      } else {
        await likeBookReview(review.id);
        setLikeCount((c) => c + 1);
      }
      setLiked(!liked);
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const renderStars = (rating: number) => {
    return (
      <div className="flex gap-0.5">
        {[1, 2, 3, 4, 5].map((star) => (
          <svg
            key={star}
            className={`w-4 h-4 ${star <= rating ? 'text-yellow-400' : 'text-gray-600'}`}
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
          </svg>
        ))}
      </div>
    );
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
  };

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 hover:border-gray-700 transition-colors">
      <div className="flex gap-4">
        {/* Book Cover */}
        {review.cover_url ? (
          <img
            src={review.cover_url}
            alt={review.title}
            className="w-20 h-28 object-cover rounded-lg flex-shrink-0"
          />
        ) : (
          <div className="w-20 h-28 bg-gradient-to-br from-purple-500/20 to-blue-500/20 rounded-lg flex-shrink-0 flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
            </svg>
          </div>
        )}

        {/* Content */}
        <div className="flex-1 min-w-0">
          {/* Header */}
          {showUser && review.user && (
            <div className="flex items-center gap-2 mb-2">
              <Link to={`/profile/${review.user.id}`}>
                <Avatar name={review.user.name} avatarUrl={review.user.avatar_url} size="sm" />
              </Link>
              <Link to={`/profile/${review.user.id}`} className="text-sm text-gray-400 hover:text-white transition-colors">
                {review.user.name}
              </Link>
              <span className="text-gray-600 text-xs">·</span>
              <span className="text-gray-500 text-xs">{formatDate(review.created_at)}</span>
            </div>
          )}

          {/* Book Info */}
          <h3 className="font-semibold text-white truncate">{review.title}</h3>
          <p className="text-sm text-gray-400 mb-2">{review.author}</p>

          {/* Rating */}
          <div className="mb-2">
            {renderStars(review.rating)}
          </div>

          {/* Review Content */}
          {review.content && (
            <p className="text-sm text-gray-300 line-clamp-3 mb-3">{review.content}</p>
          )}

          {/* Actions */}
          <div className="flex items-center gap-4">
            <button
              onClick={handleLike}
              className={`flex items-center gap-1.5 text-sm transition-colors ${
                liked ? 'text-red-400' : 'text-gray-400 hover:text-red-400'
              }`}
            >
              <svg
                className="w-4 h-4"
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
              {likeCount}
            </button>

            {isOwner && onDelete && (
              <button
                onClick={onDelete}
                className="text-sm text-gray-500 hover:text-red-400 transition-colors"
              >
                {t('common.delete')}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
