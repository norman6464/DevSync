import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import type { Question } from '../types/qa';
import { useQuestions } from '../hooks';
import QuestionCard from '../components/qa/QuestionCard';
import QuestionForm from '../components/qa/QuestionForm';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function QAPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const {
    questions, total, loading, saving,
    searchQuery, setSearchQuery,
    sort, setSort,
    page, setPage, limit,
    handleSearch,
    createQuestion, updateQuestion, deleteQuestion,
  } = useQuestions();

  const [showForm, setShowForm] = useState(false);
  const [editingQuestion, setEditingQuestion] = useState<Question | null>(null);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('qa.pageTitle')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('qa.pageSubtitle')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('qa.askQuestion')}
        </button>
      </div>

      {/* Search & Sort */}
      <div className="flex flex-wrap gap-4 mb-6">
        <div className="flex-1 min-w-[200px]">
          <div className="flex gap-2">
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              placeholder={t('qa.searchPlaceholder')}
              className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent"
            />
            <button
              onClick={handleSearch}
              className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            >
              {t('common.search')}
            </button>
          </div>
        </div>
        <div className="flex gap-2">
          {(['newest', 'votes', 'unanswered'] as const).map(s => (
            <button
              key={s}
              onClick={() => setSort(s)}
              className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                sort === s ? 'bg-gray-700 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
              }`}
            >
              {t(`qa.sort.${s}`)}
            </button>
          ))}
        </div>
      </div>

      {/* Form Modal */}
      {(showForm || editingQuestion) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingQuestion ? t('qa.editQuestion') : t('qa.newQuestion')}
            </h2>
            <QuestionForm
              question={editingQuestion || undefined}
              onSubmit={async (data) => {
                if (editingQuestion) {
                  const result = await updateQuestion(editingQuestion.id, data);
                  if (result) setEditingQuestion(null);
                } else {
                  const result = await createQuestion(data);
                  if (result) setShowForm(false);
                }
              }}
              onCancel={() => {
                setShowForm(false);
                setEditingQuestion(null);
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
      ) : questions.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z" />
            </svg>
          </div>
          <p className="text-gray-400">{t('qa.noQuestions')}</p>
          <button
            onClick={() => setShowForm(true)}
            className="mt-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            {t('qa.askFirstQuestion')}
          </button>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {questions.map(question => (
              <QuestionCard
                key={question.id}
                question={question}
                isOwner={user?.id === question.user_id}
                onEdit={() => setEditingQuestion(question)}
                onDelete={() => deleteQuestion(question)}
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
