import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { BookReview, CreateBookReviewRequest } from '../../types/bookReview';

interface BookReviewFormProps {
  review?: BookReview;
  onSubmit: (data: CreateBookReviewRequest) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

export default function BookReviewForm({ review, onSubmit, onCancel, loading }: BookReviewFormProps) {
  const { t } = useTranslation();
  const [title, setTitle] = useState(review?.title || '');
  const [author, setAuthor] = useState(review?.author || '');
  const [isbn, setIsbn] = useState(review?.isbn || '');
  const [rating, setRating] = useState(review?.rating || 5);
  const [reviewText, setReviewText] = useState(review?.review || '');
  const [imageUrl, setImageUrl] = useState(review?.image_url || '');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit({
      title,
      author,
      isbn,
      rating,
      review: reviewText,
      image_url: imageUrl,
    });
  };

  const renderStarInput = () => {
    return (
      <div className="flex gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            onClick={() => setRating(star)}
            className="focus:outline-none"
          >
            <svg
              className={`w-8 h-8 transition-colors ${
                star <= rating ? 'text-yellow-400' : 'text-gray-600 hover:text-gray-500'
              }`}
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          </button>
        ))}
      </div>
    );
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Title */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.bookTitle')} *
        </label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={300}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder={t('bookReviews.bookTitlePlaceholder')}
        />
      </div>

      {/* Author */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.author')}
        </label>
        <input
          type="text"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          maxLength={200}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder={t('bookReviews.authorPlaceholder')}
        />
      </div>

      {/* ISBN */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          ISBN
        </label>
        <input
          type="text"
          value={isbn}
          onChange={(e) => setIsbn(e.target.value)}
          maxLength={20}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder="978-4-..."
        />
      </div>

      {/* Rating */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          {t('bookReviews.rating')} *
        </label>
        {renderStarInput()}
      </div>

      {/* Review */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.review')}
        </label>
        <textarea
          value={reviewText}
          onChange={(e) => setReviewText(e.target.value)}
          rows={4}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent resize-none"
          placeholder={t('bookReviews.reviewPlaceholder')}
        />
      </div>

      {/* Image URL */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.coverImage')}
        </label>
        <input
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder="https://..."
        />
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-500 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
