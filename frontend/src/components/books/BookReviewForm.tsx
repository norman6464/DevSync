import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { createBookReview, updateBookReview, type BookReview, type CreateBookReviewRequest } from '../../api/bookReviews';

interface BookReviewFormProps {
  review?: BookReview;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function BookReviewForm({ review, onSuccess, onCancel }: BookReviewFormProps) {
  const { t } = useTranslation();
  const isEditing = !!review;

  const [title, setTitle] = useState(review?.title || '');
  const [author, setAuthor] = useState(review?.author || '');
  const [coverUrl, setCoverUrl] = useState(review?.cover_url || '');
  const [rating, setRating] = useState(review?.rating || 5);
  const [content, setContent] = useState(review?.content || '');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !author.trim()) {
      toast.error(t('bookReviews.titleAuthorRequired'));
      return;
    }

    setSubmitting(true);
    try {
      const data: CreateBookReviewRequest = {
        title: title.trim(),
        author: author.trim(),
        cover_url: coverUrl.trim() || undefined,
        rating,
        content: content.trim(),
      };

      if (isEditing && review) {
        await updateBookReview(review.id, data);
        toast.success(t('bookReviews.updateSuccess'));
      } else {
        await createBookReview(data);
        toast.success(t('bookReviews.createSuccess'));
      }
      onSuccess();
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setSubmitting(false);
    }
  };

  const renderStarInput = () => {
    return (
      <div className="flex gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <button
            key={star}
            type="button"
            onClick={() => setRating(star)}
            className={`w-8 h-8 rounded transition-colors ${
              star <= rating ? 'text-yellow-400' : 'text-gray-600 hover:text-yellow-400/50'
            }`}
          >
            <svg fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
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
          placeholder={t('bookReviews.titlePlaceholder')}
          className="w-full px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
          required
        />
      </div>

      {/* Author */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.bookAuthor')} *
        </label>
        <input
          type="text"
          value={author}
          onChange={(e) => setAuthor(e.target.value)}
          placeholder={t('bookReviews.authorPlaceholder')}
          className="w-full px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
          required
        />
      </div>

      {/* Cover URL */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.coverUrl')}
        </label>
        <input
          type="url"
          value={coverUrl}
          onChange={(e) => setCoverUrl(e.target.value)}
          placeholder={t('bookReviews.coverUrlPlaceholder')}
          className="w-full px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500"
        />
      </div>

      {/* Rating */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-2">
          {t('bookReviews.rating')}
        </label>
        {renderStarInput()}
      </div>

      {/* Content */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('bookReviews.reviewContent')}
        </label>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder={t('bookReviews.contentPlaceholder')}
          rows={5}
          className="w-full px-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-purple-500 resize-none"
        />
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-3 pt-2">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-gray-400 hover:text-white transition-colors"
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={submitting}
          className="px-6 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors disabled:opacity-50"
        >
          {submitting ? t('common.loading') : isEditing ? t('common.save') : t('bookReviews.submit')}
        </button>
      </div>
    </form>
  );
}
