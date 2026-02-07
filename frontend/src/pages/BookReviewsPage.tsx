import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import type { BookReview } from '../types/bookReview';
import { useBookReviews } from '../hooks';
import BookReviewCard from '../components/bookReviews/BookReviewCard';
import BookReviewForm from '../components/bookReviews/BookReviewForm';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function BookReviewsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const {
    reviews, total, loading, saving, page, setPage, limit,
    createReview, updateReview, deleteReview,
  } = useBookReviews();

  const [showForm, setShowForm] = useState(false);
  const [editingReview, setEditingReview] = useState<BookReview | null>(null);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('bookReviews.title')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('bookReviews.subtitle')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('bookReviews.addReview')}
        </button>
      </div>

      {/* Form Modal */}
      {(showForm || editingReview) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingReview ? t('bookReviews.editReview') : t('bookReviews.newReview')}
            </h2>
            <BookReviewForm
              review={editingReview || undefined}
              onSubmit={async (data) => {
                if (editingReview) {
                  const result = await updateReview(editingReview.id, data);
                  if (result) setEditingReview(null);
                } else {
                  const result = await createReview(data);
                  if (result) setShowForm(false);
                }
              }}
              onCancel={() => {
                setShowForm(false);
                setEditingReview(null);
              }}
              loading={saving}
            />
          </div>
        </div>
      )}

      {/* Content */}
      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <LoadingSpinner />
        </div>
      ) : reviews.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
            </svg>
          </div>
          <p className="text-gray-400">{t('bookReviews.noReviews')}</p>
          <button
            onClick={() => setShowForm(true)}
            className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
          >
            {t('bookReviews.addFirstReview')}
          </button>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {reviews.map(review => (
              <BookReviewCard
                key={review.id}
                review={review}
                isOwner={user?.id === review.user_id}
                showUser={true}
                onEdit={() => setEditingReview(review)}
                onDelete={() => deleteReview(review)}
              />
            ))}
          </div>

          {/* Pagination */}
          {total > limit && (
            <div className="flex justify-center gap-2 mt-8">
              <button
                onClick={() => setPage(p => Math.max(0, p - 1))}
                disabled={page === 0}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.previous')}
              </button>
              <span className="px-4 py-2 text-gray-400">
                {page + 1} / {Math.ceil(total / limit)}
              </span>
              <button
                onClick={() => setPage(p => p + 1)}
                disabled={(page + 1) * limit >= total}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.next')}
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
