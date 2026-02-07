import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { BookReview, CreateBookReviewRequest } from '../types/bookReview';
import { getBookReviews, createBookReview, updateBookReview, deleteBookReview } from '../api/bookReviews';
import { useAsyncData } from './useAsyncData';

export function useBookReviews() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);
  const [page, setPage] = useState(0);
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
  const currentReviews = localReviews ?? reviews;

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
    if (!confirm(t('bookReviews.confirmDelete'))) return false;
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

  return {
    reviews: currentReviews,
    total,
    loading,
    saving,
    page,
    setPage,
    limit,
    createReview: handleCreate,
    updateReview: handleUpdate,
    deleteReview: handleDelete,
    refetch,
  };
}
