import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Star } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { BookReview, ReviewStatus } from '../types/bookReview';
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
    statusFilter, setStatusFilter, showArchived, setShowArchived,
    sortBy, setSortBy, ratingFilter, setRatingFilter,
    createReview, updateReview, deleteReview,
    updateStatus, archiveReview, unarchiveReview,
  } = useBookReviews();

  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingReview, setEditingReview] = useState<BookReview | null>(null);

  const handleDeleteReview = useCallback(async (review: BookReview) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('bookReviews.confirmDelete'), variant: 'danger' });
    if (ok) deleteReview(review);
  }, [confirm, t, deleteReview]);

  const handleFormClose = useCallback(() => {
    setShowForm(false);
    setEditingReview(null);
  }, []);

  const handleFormSubmit = useCallback(async (data: Parameters<typeof createReview>[0]) => {
    if (editingReview) {
      const result = await updateReview(editingReview.id, data);
      if (result) setEditingReview(null);
    } else {
      const result = await createReview(data);
      if (result) setShowForm(false);
    }
  }, [editingReview, updateReview, createReview]);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <PageHeader
        title={t('bookReviews.title')}
        subtitle={t('bookReviews.subtitle')}
        actionLabel={t('bookReviews.addReview')}
        onAction={() => setShowForm(true)}
      />

      {/* Filter Bar */}
      <div className="flex flex-wrap items-center gap-3 mb-6">
        {(['all', 'not_started', 'reading', 'completed'] as const).map(status => (
          <button
            key={status}
            onClick={() => setStatusFilter(status)}
            className={`px-3 py-1.5 text-sm rounded-lg transition-colors ${
              statusFilter === status
                ? 'bg-blue-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            {status === 'all' ? t('bookReviews.statusAll') : t(`bookReviews.status.${status}`)}
          </button>
        ))}
        <button
          onClick={() => setShowArchived(!showArchived)}
          className={`px-3 py-1.5 text-sm rounded-lg transition-colors ${
            showArchived
              ? 'bg-yellow-700 text-yellow-200'
              : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
          }`}
        >
          {t('bookReviews.archivedLabel')}
        </button>

        <div className="flex-1" />
        {(['newest', 'oldest', 'ratingDesc', 'ratingAsc'] as const).map(sort => (
          <button
            key={sort}
            onClick={() => setSortBy(sort)}
            className={`px-3 py-1.5 text-sm rounded-lg transition-colors ${
              sortBy === sort
                ? 'bg-orange-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            {t(`bookReviews.sort.${sort}`)}
          </button>
        ))}
      </div>

      {/* Rating Filter */}
      <div className="flex flex-wrap items-center gap-2 mb-6">
        <Star className="w-4 h-4 text-yellow-400" />
        <button
          onClick={() => setRatingFilter(0)}
          className={`px-3 py-1.5 text-sm rounded-lg transition-colors ${
            ratingFilter === 0
              ? 'bg-yellow-600 text-white'
              : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
          }`}
        >
          {t('bookReviews.ratingFilter.all')}
        </button>
        {[3, 4, 5].map(star => (
          <button
            key={star}
            onClick={() => setRatingFilter(star)}
            className={`px-3 py-1.5 text-sm rounded-lg transition-colors inline-flex items-center gap-1 ${
              ratingFilter === star
                ? 'bg-yellow-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            }`}
          >
            {t('bookReviews.ratingFilter.minRating', { stars: star })}
          </button>
        ))}
      </div>

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingReview}
        onClose={handleFormClose}
        title={editingReview ? t('bookReviews.editReview') : t('bookReviews.newReview')}
        maxWidth="max-w-lg"
      >
        <BookReviewForm
          review={editingReview || undefined}
          onSubmit={handleFormSubmit}
          onCancel={handleFormClose}
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
                onStatusChange={updateStatus}
                onArchive={archiveReview}
                onUnarchive={unarchiveReview}
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
