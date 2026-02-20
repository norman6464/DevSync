import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { BookReview, CreateBookReviewRequest, ReviewStatus } from '../types/bookReview';
import {
  getBookReviews, createBookReview, updateBookReview, deleteBookReview,
  archiveBookReview, unarchiveBookReview, updateBookReviewStatus,
} from '../api/bookReviews';
import { useAsyncData } from './useAsyncData';

export function useBookReviews() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);
  const [page, setPage] = useState(0);
  const [statusFilter, setStatusFilter] = useState<ReviewStatus | 'all'>('all');
  const [showArchived, setShowArchived] = useState(false);
  const limit = 20;

  const { data, loading, refetch } = useAsyncData(
    async () => {
      return await getBookReviews(limit, page * limit);
    },
    { deps: [page] }
  );

  const reviews = data?.reviews ?? [];
  const total = data?.total ?? 0;

  const [localReviews, setLocalReviews] = useState<BookReview[] | null>(null);
  const allReviews = localReviews ?? reviews;
  const currentReviews = allReviews.filter(r => {
    if (!showArchived && r.is_archived) return false;
    if (showArchived && !r.is_archived) return false;
    if (statusFilter !== 'all' && r.status !== statusFilter) return false;
    return true;
  });

  const handleCreate = useCallback(async (reqData: CreateBookReviewRequest) => {
    setSaving(true);
    try {
      const newReview = await createBookReview(reqData);
      setLocalReviews(prev => [newReview, ...(prev ?? reviews)]);
      toast.success(t('bookReviews.createSuccess'));
      return newReview;
    } catch {
      toast.error(t('bookReviews.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, reviews]);

  const handleUpdate = useCallback(async (reviewId: number, reqData: CreateBookReviewRequest) => {
    setSaving(true);
    try {
      const updated = await updateBookReview(reviewId, reqData);
      setLocalReviews(prev => (prev ?? reviews).map(r => r.id === updated.id ? updated : r));
      toast.success(t('bookReviews.updateSuccess'));
      return updated;
    } catch {
      toast.error(t('bookReviews.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, reviews]);

  const handleDelete = useCallback(async (review: BookReview) => {
    try {
      await deleteBookReview(review.id);
      setLocalReviews(prev => (prev ?? reviews).filter(r => r.id !== review.id));
      toast.success(t('bookReviews.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('bookReviews.deleteFailed'));
      return false;
    }
  }, [t, reviews]);

  const handleUpdateStatus = useCallback(async (reviewId: number, status: ReviewStatus) => {
    try {
      const updated = await updateBookReviewStatus(reviewId, status);
      setLocalReviews(prev => (prev ?? reviews).map(r => r.id === updated.id ? updated : r));
      toast.success(t('bookReviews.statusUpdated'));
      return updated;
    } catch {
      toast.error(t('bookReviews.statusUpdateFailed'));
      return null;
    }
  }, [t, reviews]);

  const handleArchive = useCallback(async (reviewId: number) => {
    try {
      const updated = await archiveBookReview(reviewId);
      setLocalReviews(prev => (prev ?? reviews).map(r => r.id === updated.id ? updated : r));
      toast.success(t('bookReviews.archived'));
      return updated;
    } catch {
      toast.error(t('bookReviews.archiveFailed'));
      return null;
    }
  }, [t, reviews]);

  const handleUnarchive = useCallback(async (reviewId: number) => {
    try {
      const updated = await unarchiveBookReview(reviewId);
      setLocalReviews(prev => (prev ?? reviews).map(r => r.id === updated.id ? updated : r));
      toast.success(t('bookReviews.unarchived'));
      return updated;
    } catch {
      toast.error(t('bookReviews.unarchiveFailed'));
      return null;
    }
  }, [t, reviews]);

  return {
    reviews: currentReviews,
    total,
    loading,
    saving,
    page,
    setPage,
    limit,
    statusFilter,
    setStatusFilter,
    showArchived,
    setShowArchived,
    createReview: handleCreate,
    updateReview: handleUpdate,
    deleteReview: handleDelete,
    updateStatus: handleUpdateStatus,
    archiveReview: handleArchive,
    unarchiveReview: handleUnarchive,
    refetch,
  };
}
