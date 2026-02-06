import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getBookReviews, deleteBookReview, type BookReview } from '../api/bookReviews';
import { useAuthStore } from '../store/authStore';
import BookReviewCard from '../components/books/BookReviewCard';
import BookReviewForm from '../components/books/BookReviewForm';
import LoadingSpinner from '../components/common/LoadingSpinner';
import toast from 'react-hot-toast';

export default function BookReviewsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [reviews, setReviews] = useState<BookReview[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  const fetchReviews = async () => {
    try {
      const res = await getBookReviews();
      setReviews(res.data || []);
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReviews();
  }, []);

  const handleDelete = async (id: number) => {
    if (!window.confirm(t('bookReviews.confirmDelete'))) return;
    try {
      await deleteBookReview(id);
      setReviews((prev) => prev.filter((r) => r.id !== id));
      toast.success(t('bookReviews.deleteSuccess'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleFormSuccess = () => {
    setShowForm(false);
    fetchReviews();
  };

  if (loading) {
    return (
      <div className="py-12">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t('bookReviews.title')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('bookReviews.description')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('bookReviews.addReview')}
        </button>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={() => setShowForm(false)} />
          <div className="relative bg-gray-900 border border-gray-800 rounded-xl shadow-2xl max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <h2 className="text-lg font-semibold mb-4">{t('bookReviews.newReview')}</h2>
              <BookReviewForm onSuccess={handleFormSuccess} onCancel={() => setShowForm(false)} />
            </div>
          </div>
        </div>
      )}

      {/* Reviews List */}
      {reviews.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-xl p-12 text-center">
          <svg className="w-16 h-16 mx-auto text-gray-700 mb-4" fill="none" stroke="currentColor" strokeWidth="1" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
          </svg>
          <p className="text-gray-400">{t('bookReviews.noReviews')}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {reviews.map((review) => (
            <BookReviewCard
              key={review.id}
              review={review}
              isOwner={user?.id === review.user_id}
              onDelete={() => handleDelete(review.id)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
