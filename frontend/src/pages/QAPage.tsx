import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { HelpCircle } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { Question } from '../types/qa';
import { useQuestions, useDebounce, useConfirm } from '../hooks';
import QuestionCard from '../components/qa/QuestionCard';
import QuestionForm from '../components/qa/QuestionForm';
import { QuestionCardSkeleton } from '../components/common/Skeleton';
import { EmptyState, Modal, Pagination, SearchInput, PageHeader } from '../components/common';
import ConfirmDialog from '../components/common/ConfirmDialog';

export default function QAPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const {
    questions, total, loading, saving,
    searchQuery, setSearchQuery,
    sort, setSort,
    solvedFilter, setSolvedFilter,
    page, setPage, limit,
    handleSearch,
    createQuestion, updateQuestion, deleteQuestion,
  } = useQuestions();

  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingQuestion, setEditingQuestion] = useState<Question | null>(null);

  const handleDeleteQuestion = useCallback(async (question: Question) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('qa.confirmDelete'), variant: 'danger' });
    if (ok) deleteQuestion(question);
  }, [confirm, t, deleteQuestion]);

  const handleFormClose = useCallback(() => {
    setShowForm(false);
    setEditingQuestion(null);
  }, []);

  const handleFormSubmit = useCallback(async (data: Parameters<typeof createQuestion>[0]) => {
    if (editingQuestion) {
      const result = await updateQuestion(editingQuestion.id, data);
      if (result) setEditingQuestion(null);
    } else {
      const result = await createQuestion(data);
      if (result) setShowForm(false);
    }
  }, [editingQuestion, updateQuestion, createQuestion]);

  // デバウンス処理（300ms）
  const debouncedQuery = useDebounce(searchQuery, 300);

  // デバウンスされたクエリで自動検索
  useEffect(() => {
    if (debouncedQuery !== undefined) {
      handleSearch();
    }
  }, [debouncedQuery, handleSearch]);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <PageHeader
        title={t('qa.pageTitle')}
        subtitle={t('qa.pageSubtitle')}
        actionLabel={t('qa.askQuestion')}
        onAction={() => setShowForm(true)}
      />

      {/* Search & Sort */}
      <div className="flex flex-wrap gap-4 mb-6">
        <div className="flex-1 min-w-[200px]">
          <SearchInput
            value={searchQuery}
            onChange={setSearchQuery}
            onSearch={handleSearch}
            placeholder={t('qa.searchPlaceholder')}
            showButton
          />
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
        <div className="flex gap-2">
          {(['all', 'solved', 'unsolved'] as const).map(f => (
            <button
              key={f}
              onClick={() => setSolvedFilter(f)}
              className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                solvedFilter === f ? 'bg-green-700/50 text-green-300' : 'bg-gray-800 text-gray-400 hover:text-white'
              }`}
            >
              {f === 'all' ? t('common.all') : t(`qa.filter.${f}`)}
            </button>
          ))}
        </div>
      </div>

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingQuestion}
        onClose={handleFormClose}
        title={editingQuestion ? t('qa.editQuestion') : t('qa.newQuestion')}
      >
        <QuestionForm
          question={editingQuestion || undefined}
          onSubmit={handleFormSubmit}
          onCancel={handleFormClose}
          loading={saving}
        />
      </Modal>

      {/* Content */}
      {loading ? (
        <div className="space-y-4">
          <QuestionCardSkeleton />
          <QuestionCardSkeleton />
          <QuestionCardSkeleton />
        </div>
      ) : questions.length === 0 ? (
        <EmptyState
          icon={HelpCircle}
          title={t('qa.noQuestions')}
          actionLabel={t('qa.askFirstQuestion')}
          onAction={() => setShowForm(true)}
        />
      ) : (
        <>
          <div className="space-y-4">
            {questions.map(question => (
              <QuestionCard
                key={question.id}
                question={question}
                isOwner={user?.id === question.user_id}
                onEdit={() => setEditingQuestion(question)}
                onDelete={() => handleDeleteQuestion(question)}
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
