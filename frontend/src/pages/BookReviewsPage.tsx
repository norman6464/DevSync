import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { BookReview } from '../types/bookReview';
import { useBookReviews, useConfirm } from '../hooks';
import BookReviewCard from '../components/bookReviews/BookReviewCard';
import BookReviewForm from '../components/bookReviews/BookReviewForm';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { EmptyState, Modal, Pagination } from '../components/common';
import ConfirmDialog from '../components/common/ConfirmDialog';

export default function BookReviewsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const {
    reviews, total, loading, saving, page, setPage, limit,
    createReview, updateReview, deleteReview,
  } = useBookReviews();

  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingReview, setEditingReview] = useState<BookReview | null>(null);

  const handleDeleteReview = async (review: BookReview) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('bookReviews.confirmDelete'), variant: 'danger' });
    if (ok) deleteReview(review);
  };

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('bookReviews.title')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('bookReviews.subtitle')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('bookReviews.addReview')}
        </button>
      </div>

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingReview}
        onClose={() => { setShowForm(false); setEditingReview(null); }}
        title={editingReview ? t('bookReviews.editReview') : t('bookReviews.newReview')}
        maxWidth="max-w-lg"
      >
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
      </Modal>

      {/* Content */}
      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <LoadingSpinner />
        </div>
      ) : reviews.length === 0 ? (
        <EmptyState
          icon={BookOpen}
          message={t('bookReviews.noReviews')}
          actionLabel={t('bookReviews.addFirstReview')}
          onAction={() => setShowForm(true)}
        />
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
                onDelete={() => handleDeleteReview(review)}
              />
            ))}
          </div>

          <Pagination
            currentPage={page}
            totalItems={total}
            itemsPerPage={limit}
            onPageChange={setPage}
          />
        </>
      )}
      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
