import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { BookReview } from '../types/bookReview';
import { useBookReviews, useConfirm } from '../hooks';
import BookReviewCard from '../components/bookReviews/BookReviewCard';
import BookReviewForm from '../components/bookReviews/BookReviewForm';
import { EmptyState, Modal, Pagination, PageHeader, PageLoader } from '../components/common';
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
      <PageHeader
        title={t('bookReviews.title')}
        subtitle={t('bookReviews.subtitle')}
        actionLabel={t('bookReviews.addReview')}
        onAction={() => setShowForm(true)}
      />

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
        <PageLoader />
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
