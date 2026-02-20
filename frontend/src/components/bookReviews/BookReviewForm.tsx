import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { BookReview, CreateBookReviewRequest } from '../../types/bookReview';
import { buttonSecondaryClass, inputClass, labelClass, textareaClass } from '../../constants/styles';
import StarRating from '../common/StarRating';

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

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Title */}
      <div>
        <label htmlFor="review-title" className={labelClass}>
          {t('bookReviews.bookTitle')} *
        </label>
        <input
          id="review-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={300}
          className={inputClass}
          placeholder={t('bookReviews.bookTitlePlaceholder')}
        />
      </div>

      {/* Author */}
      <div>
        <label htmlFor="review-author" className={labelClass}>
          {t('bookReviews.author')}
        </label>
        <input
          id="review-author"
          type="text"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          maxLength={200}
          className={inputClass}
          placeholder={t('bookReviews.authorPlaceholder')}
        />
      </div>

      {/* ISBN */}
      <div>
        <label htmlFor="review-isbn" className={labelClass}>
          ISBN
        </label>
        <input
          id="review-isbn"
          type="text"
          value={isbn}
          onChange={(e) => setIsbn(e.target.value)}
          maxLength={20}
          className={inputClass}
          placeholder="978-4-..."
        />
      </div>

      {/* Rating */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          {t('bookReviews.rating')} *
        </label>
        <StarRating rating={rating} onChange={setRating} size="lg" />
      </div>

      {/* Review */}
      <div>
        <label htmlFor="review-text" className={labelClass}>
          {t('bookReviews.review')}
        </label>
        <textarea
          id="review-text"
          value={reviewText}
          onChange={(e) => setReviewText(e.target.value)}
          rows={4}
          maxLength={2000}
          className={textareaClass}
          placeholder={t('bookReviews.reviewPlaceholder')}
        />
        <p className="text-xs text-gray-500 text-right mt-1">{reviewText.length}/2000</p>
      </div>

      {/* Image URL */}
      <div>
        <label htmlFor="review-image-url" className={labelClass}>
          {t('bookReviews.coverImage')}
        </label>
        <input
          id="review-image-url"
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          maxLength={2048}
          className={inputClass}
          placeholder="https://..."
        />
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className={`flex-1 ${buttonSecondaryClass}`}
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className={`flex-1 ${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
