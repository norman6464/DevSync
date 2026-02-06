import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { BookReview } from '../../types/bookReview';
import Avatar from '../common/Avatar';

interface BookReviewCardProps {
  review: BookReview;
  onEdit?: () => void;
  onDelete?: () => void;
  isOwner?: boolean;
  showUser?: boolean;
}

export default function BookReviewCard({
  review,
  onEdit,
  onDelete,
  isOwner = false,
  showUser = true,
}: BookReviewCardProps) {
  const { t } = useTranslation();

  const renderStars = (rating: number) => {
    return (
      <div className="flex gap-0.5">
        {[1, 2, 3, 4, 5].map((star) => (
          <svg
            key={star}
            className={`w-4 h-4 ${star <= rating ? 'text-yellow-400' : 'text-gray-600'}`}
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
        ))}
      </div>
    );
  };

  return (
    <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700 hover:border-gray-600 transition-colors">
      <div className="flex">
        {/* Book Cover */}
        {review.image_url && (
          <div className="w-24 min-h-32 bg-gray-700 flex-shrink-0">
            <img
              src={review.image_url}
              alt={review.title}
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
                  className="p-1.5 text-gray-400 hover:text-white transition-colors"
                  title={t('common.edit')}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                  </svg>
                </button>
                <button
                  onClick={onDelete}
                  className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
                  title={t('common.delete')}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                  </svg>
                </button>
              </div>
            )}
          </div>

          {/* Rating */}
          <div className="mt-2">
            {renderStars(review.rating)}
          </div>

          {/* Review Text */}
          {review.review && (
            <p className="text-gray-300 text-sm mt-2 line-clamp-2">{review.review}</p>
          )}

          {/* Footer */}
          <div className="flex items-center justify-between mt-3 pt-2 border-t border-gray-700">
            {showUser && review.user && (
              <Link
                to={`/profile/${review.user_id}`}
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
