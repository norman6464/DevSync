import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Pencil, Trash2 } from 'lucide-react';
import type { BookReview, ReviewStatus } from '../../types/bookReview';
import Avatar from '../common/Avatar';
import StarRating from '../common/StarRating';
import { cardClass, iconButtonClass, deleteIconButtonClass } from '../../constants/styles';
import { sanitizeUrl } from '../../utils/url';

const statusColors: Record<ReviewStatus, string> = {
  not_started: 'bg-gray-600 text-gray-200',
  reading: 'bg-blue-600 text-blue-100',
  completed: 'bg-green-600 text-green-100',
};

interface BookReviewCardProps {
  review: BookReview;
  onEdit?: () => void;
  onDelete?: () => void;
  onStatusChange?: (id: number, status: ReviewStatus) => void;
  onArchive?: (id: number) => void;
  onUnarchive?: (id: number) => void;
  isOwner?: boolean;
  showUser?: boolean;
}

export default function BookReviewCard({
  review,
  onEdit,
  onDelete,
  onStatusChange,
  onArchive,
  onUnarchive,
  isOwner = false,
  showUser = true,
}: BookReviewCardProps) {
  const { t } = useTranslation();

  return (
    <div className={cardClass}>
      <div className="flex">
        {/* Book Cover */}
        {sanitizeUrl(review.image_url) && (
          <div className="w-24 min-h-32 bg-gray-700 flex-shrink-0">
            <img
              src={sanitizeUrl(review.image_url)!}
              alt={review.title}
              referrerPolicy="no-referrer"
              className="w-full h-full object-cover"
            />
          </div>
        )}

        <div className="flex-1 p-4">
          {/* Header */}
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-white line-clamp-1">{review.title}</h3>
              {review.author && (
                <p className="text-sm text-gray-400">{review.author}</p>
              )}
            </div>

            {isOwner && (
              <div className="flex gap-1">
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

          {/* Status & Rating */}
          <div className="mt-2 flex items-center gap-3">
            <StarRating rating={review.rating} />
            {review.status && (
              <span className={`text-xs px-2 py-0.5 rounded-full ${statusColors[review.status]}`}>
                {t(`bookReviews.status.${review.status}`)}
              </span>
            )}
            {review.is_archived && (
              <span className="text-xs px-2 py-0.5 rounded-full bg-yellow-700 text-yellow-200">
                {t('bookReviews.archivedLabel')}
              </span>
            )}
          </div>

          {/* Review Text */}
          {review.review && (
            <p className="text-gray-300 text-sm mt-2 line-clamp-2">{review.review}</p>
          )}

          {/* Owner Actions */}
          {isOwner && (onStatusChange || onArchive || onUnarchive) && (
            <div className="flex items-center gap-2 mt-2">
              {onStatusChange && (
                <select
                  value={review.status}
                  onChange={(e) => onStatusChange(review.id, e.target.value as ReviewStatus)}
                  className="text-xs bg-gray-700 text-gray-300 border border-gray-600 rounded px-2 py-1"
                >
                  <option value="not_started">{t('bookReviews.status.not_started')}</option>
                  <option value="reading">{t('bookReviews.status.reading')}</option>
                  <option value="completed">{t('bookReviews.status.completed')}</option>
                </select>
              )}
              {onArchive && !review.is_archived && (
                <button
                  onClick={() => onArchive(review.id)}
                  className="text-xs text-gray-400 hover:text-yellow-400 transition-colors"
                >
                  {t('bookReviews.archive')}
                </button>
              )}
              {onUnarchive && review.is_archived && (
                <button
                  onClick={() => onUnarchive(review.id)}
                  className="text-xs text-gray-400 hover:text-blue-400 transition-colors"
                >
                  {t('bookReviews.unarchive')}
                </button>
              )}
            </div>
          )}

          {/* Footer */}
          <div className="flex items-center justify-between mt-3 pt-2 border-t border-gray-700">
            {showUser && review.user && (
              <Link
                to={`/profile/${review.user?.username || review.user_id}`}
                className="flex items-center gap-2 hover:opacity-80 transition-opacity"
              >
                <Avatar avatarUrl={review.user.avatar_url} name={review.user.name} size="sm" />
                <span className="text-sm text-gray-400">{review.user.name}</span>
              </Link>
            )}

            {review.isbn && (
              <span className="text-xs text-gray-500">ISBN: {review.isbn}</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
