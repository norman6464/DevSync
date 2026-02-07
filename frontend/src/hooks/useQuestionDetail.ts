import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { Answer, CreateAnswerRequest } from '../types/qa';
import {
  getQuestionById,
  voteQuestion,
  removeQuestionVote,
  createAnswer,
  updateAnswer,
  deleteAnswer,
  setBestAnswer,
  voteAnswer,
  removeAnswerVote,
} from '../api/qa';
import client from '../api/client';
import { useAsyncData } from './useAsyncData';

export function useQuestionDetail(id: string | undefined) {
  const { t } = useTranslation();
  const questionId = id ? parseInt(id) : 0;
  const [submitting, setSubmitting] = useState(false);

  const { data, loading, refetch } = useAsyncData(
    async () => {
      const [questionResult, answersRes] = await Promise.all([
        getQuestionById(questionId),
        client.get(`/questions/${questionId}/answers`),
      ]);
      return {
        ...questionResult,
        answers: (answersRes.data || []) as Answer[],
      };
    },
    { deps: [questionId], enabled: !!questionId }
  );

  const question = data?.question ?? null;
  const userVote = data?.user_vote ?? 0;
  const fetchedAnswers = data?.answers ?? [];

  const [localAnswers, setLocalAnswers] = useState<Answer[] | null>(null);
  const answers = localAnswers ?? fetchedAnswers;

  const handleVoteQuestion = useCallback(async (value: 1 | -1) => {
    try {
      await voteQuestion(questionId, { value });
      await refetch();
    } catch {
      toast.error(t('qa.voteFailed'));
    }
  }, [questionId, refetch, t]);

  const handleRemoveQuestionVote = useCallback(async () => {
    try {
      await removeQuestionVote(questionId);
      await refetch();
    } catch {
      toast.error(t('qa.voteFailed'));
    }
  }, [questionId, refetch, t]);

  const handleCreateAnswer = useCallback(async (reqData: CreateAnswerRequest) => {
    if (!reqData.body.trim() || !questionId) return false;
    setSubmitting(true);
    try {
      const newAnswer = await createAnswer(questionId, reqData);
      setLocalAnswers(prev => [...(prev ?? fetchedAnswers), newAnswer]);
      toast.success(t('qa.answerCreated'));
      await refetch();
      return true;
    } catch {
      toast.error(t('qa.answerFailed'));
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [questionId, refetch, t, fetchedAnswers]);

  const handleUpdateAnswer = useCallback(async (answerId: number, reqData: CreateAnswerRequest) => {
    setSubmitting(true);
    try {
      const updated = await updateAnswer(questionId, answerId, reqData);
      setLocalAnswers(prev => (prev ?? fetchedAnswers).map(a => a.id === updated.id ? updated : a));
      toast.success(t('qa.answerUpdated'));
      return true;
    } catch {
      toast.error(t('qa.answerUpdateFailed'));
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [questionId, t, fetchedAnswers]);

  const handleDeleteAnswer = useCallback(async (answerId: number) => {
    if (!confirm(t('qa.confirmDeleteAnswer'))) return false;
    try {
      await deleteAnswer(questionId, answerId);
      setLocalAnswers(prev => (prev ?? fetchedAnswers).filter(a => a.id !== answerId));
      await refetch();
      toast.success(t('qa.answerDeleted'));
      return true;
    } catch {
      toast.error(t('qa.answerDeleteFailed'));
      return false;
    }
  }, [questionId, refetch, t, fetchedAnswers]);

  const handleSetBestAnswer = useCallback(async (answerId: number) => {
    try {
      await setBestAnswer(questionId, answerId);
      setLocalAnswers(prev => (prev ?? fetchedAnswers).map(a => ({
        ...a,
        is_best: a.id === answerId,
      })));
      await refetch();
      toast.success(t('qa.bestAnswerSet'));
    } catch {
      toast.error(t('qa.bestAnswerFailed'));
    }
  }, [questionId, refetch, t, fetchedAnswers]);

  const handleVoteAnswer = useCallback(async (answerId: number, value: 1 | -1) => {
    try {
      await voteAnswer(questionId, answerId, { value });
      await refetch();
      setLocalAnswers(null); // reset to use fetched data
    } catch {
      toast.error(t('qa.voteFailed'));
    }
  }, [questionId, refetch, t]);

  const handleRemoveAnswerVote = useCallback(async (answerId: number) => {
    try {
      await removeAnswerVote(questionId, answerId);
      await refetch();
      setLocalAnswers(null);
    } catch {
      toast.error(t('qa.voteFailed'));
    }
  }, [questionId, refetch, t]);

  return {
    question,
    userVote,
    answers,
    loading,
    submitting,
    voteQuestion: handleVoteQuestion,
    removeQuestionVote: handleRemoveQuestionVote,
    createAnswer: handleCreateAnswer,
    updateAnswer: handleUpdateAnswer,
    deleteAnswer: handleDeleteAnswer,
    setBestAnswer: handleSetBestAnswer,
    voteAnswer: handleVoteAnswer,
    removeAnswerVote: handleRemoveAnswerVote,
    refetch,
  };
}
