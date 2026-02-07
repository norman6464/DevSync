import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { Question, CreateQuestionRequest } from '../types/qa';
import {
  getQuestions,
  searchQuestions,
  createQuestion,
  updateQuestion,
  deleteQuestion,
} from '../api/qa';
import { useAsyncData } from './useAsyncData';

type SortType = 'newest' | 'votes' | 'unanswered';

export function useQuestions() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [tagFilter, setTagFilter] = useState('');
  const [sort, setSort] = useState<SortType>('newest');
  const [page, setPage] = useState(0);
  const limit = 20;

  const { data, loading, refetch } = useAsyncData(
    async () => {
      if (searchQuery.trim()) {
        return await searchQuestions(searchQuery, limit, page * limit);
      }
      return await getQuestions(limit, page * limit, tagFilter, sort);
    },
    { deps: [tagFilter, sort, page] }
  );

  const questions = data?.questions ?? [];
  const total = data?.total ?? 0;

  const [localQuestions, setLocalQuestions] = useState<Question[] | null>(null);
  const currentQuestions = localQuestions ?? questions;

  const handleSearch = useCallback(() => {
    setPage(0);
    refetch();
  }, [refetch]);

  const handleCreate = useCallback(async (reqData: CreateQuestionRequest) => {
    setSaving(true);
    try {
      const newQuestion = await createQuestion(reqData);
      setLocalQuestions(prev => [newQuestion, ...(prev ?? questions)]);
      toast.success(t('qa.createSuccess'));
      return newQuestion;
    } catch {
      toast.error(t('qa.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, questions]);

  const handleUpdate = useCallback(async (questionId: number, reqData: CreateQuestionRequest) => {
    setSaving(true);
    try {
      const updated = await updateQuestion(questionId, reqData);
      setLocalQuestions(prev => (prev ?? questions).map(q => q.id === updated.id ? updated : q));
      toast.success(t('qa.updateSuccess'));
      return updated;
    } catch {
      toast.error(t('qa.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, questions]);

  const handleDelete = useCallback(async (question: Question) => {
    if (!confirm(t('qa.confirmDelete'))) return false;
    try {
      await deleteQuestion(question.id);
      setLocalQuestions(prev => (prev ?? questions).filter(q => q.id !== question.id));
      toast.success(t('qa.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('qa.deleteFailed'));
      return false;
    }
  }, [t, questions]);

  const changeSort = useCallback((newSort: SortType) => {
    setSort(newSort);
    setPage(0);
    setLocalQuestions(null);
  }, []);

  return {
    questions: currentQuestions,
    total,
    loading,
    saving,
    searchQuery,
    setSearchQuery,
    tagFilter,
    setTagFilter: (v: string) => { setTagFilter(v); setPage(0); setLocalQuestions(null); },
    sort,
    setSort: changeSort,
    page,
    setPage,
    limit,
    handleSearch,
    createQuestion: handleCreate,
    updateQuestion: handleUpdate,
    deleteQuestion: handleDelete,
    refetch,
  };
}
