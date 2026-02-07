import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { useQuestionDetail } from '../hooks';
import type { Answer } from '../types/qa';
import AnswerCard from '../components/qa/AnswerCard';
import AnswerForm from '../components/qa/AnswerForm';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function QADetailPage() {
  const { id } = useParams<{ id: string }>();
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const {
    question, userVote, answers,
    loading, submitting,
    voteQuestion, removeQuestionVote,
    createAnswer, updateAnswer, deleteAnswer,
    setBestAnswer, voteAnswer, removeAnswerVote,
  } = useQuestionDetail(id);

  const [editingAnswer, setEditingAnswer] = useState<Answer | null>(null);

  const tags: string[] = question?.tags ? (() => {
    try { return JSON.parse(question.tags); } catch { return []; }
  })() : [];

  const isQuestionOwner = user?.id === question?.user_id;

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <LoadingSpinner />
      </div>
    );
  }

  if (!question) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8 text-center">
        <p className="text-gray-400">{t('qa.questionNotFound')}</p>
        <Link to="/qa" className="text-green-400 hover:text-green-300 mt-2 inline-block">
          {t('qa.backToList')}
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      {/* Back link */}
      <Link
        to="/qa"
        className="inline-flex items-center gap-1 text-gray-400 hover:text-white transition-colors mb-6"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18" />
        </svg>
        {t('qa.backToList')}
      </Link>

      {/* Question */}
      <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
        <div className="flex gap-4">
          {/* Vote buttons */}
          <div className="flex flex-col items-center gap-1">
            <button
              onClick={() => userVote === 1 ? removeQuestionVote() : voteQuestion(1)}
              className={`p-1 transition-colors ${userVote === 1 ? 'text-green-400' : 'text-gray-400 hover:text-green-400'}`}
              title={t('qa.upvote')}
            >
              <svg className="w-8 h-8" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 15.75l7.5-7.5 7.5 7.5" />
              </svg>
            </button>
            <span className={`text-xl font-bold ${
              question.vote_count > 0 ? 'text-green-400' : question.vote_count < 0 ? 'text-red-400' : 'text-gray-400'
            }`}>
              {question.vote_count}
            </span>
            <button
              onClick={() => userVote === -1 ? removeQuestionVote() : voteQuestion(-1)}
              className={`p-1 transition-colors ${userVote === -1 ? 'text-red-400' : 'text-gray-400 hover:text-red-400'}`}
              title={t('qa.downvote')}
            >
              <svg className="w-8 h-8" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
              </svg>
            </button>
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <h1 className="text-2xl font-bold text-white">
              {question.is_solved && (
                <span className="inline-flex items-center mr-2 text-green-400">
                  <svg className="w-6 h-6" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                </span>
              )}
              {question.title}
            </h1>

            <p className="text-gray-200 mt-4 whitespace-pre-wrap">{question.body}</p>

            {/* Tags */}
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-1.5 mt-4">
                {tags.map(tag => (
                  <span key={tag} className="px-2 py-0.5 bg-blue-600/20 text-blue-400 text-xs rounded-md">
                    {tag}
                  </span>
                ))}
              </div>
            )}

            {/* Meta */}
            <div className="flex items-center justify-between mt-4 pt-3 border-t border-gray-700/50">
              {question.user && (
                <Link
                  to={`/profile/${question.user_id}`}
                  className="flex items-center gap-2 hover:opacity-80 transition-opacity"
                >
                  <Avatar avatarUrl={question.user.avatar_url} name={question.user.name} size="sm" />
                  <div>
                    <span className="text-sm text-gray-300">{question.user.name}</span>
                    <span className="text-xs text-gray-500 ml-2">
                      {new Date(question.created_at).toLocaleDateString()}
                    </span>
                  </div>
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Answers Section */}
      <div className="mt-8">
        <h2 className="text-lg font-semibold text-white mb-4">
          {t('qa.answersCount', { count: answers.length })}
        </h2>

        <div className="space-y-4">
          {answers.map(answer => (
            editingAnswer?.id === answer.id ? (
              <div key={answer.id} className="bg-gray-800 rounded-xl p-5 border border-gray-700">
                <AnswerForm
                  initialBody={answer.body}
                  onSubmit={async (body) => {
                    const success = await updateAnswer(answer.id, { body });
                    if (success) setEditingAnswer(null);
                    return success;
                  }}
                  onCancel={() => setEditingAnswer(null)}
                  loading={submitting}
                  isEdit
                />
              </div>
            ) : (
              <AnswerCard
                key={answer.id}
                answer={answer}
                isOwner={user?.id === answer.user_id}
                isQuestionOwner={isQuestionOwner}
                onEdit={() => setEditingAnswer(answer)}
                onDelete={() => deleteAnswer(answer.id)}
                onSetBest={() => setBestAnswer(answer.id)}
                onVote={(value) => voteAnswer(answer.id, value)}
                onRemoveVote={() => removeAnswerVote(answer.id)}
              />
            )
          ))}
        </div>

        {/* Answer Form */}
        <div className="mt-8">
          <h3 className="text-md font-semibold text-white mb-3">{t('qa.yourAnswer')}</h3>
          <AnswerForm
            onSubmit={async (body) => {
              return await createAnswer({ body });
            }}
            loading={submitting}
          />
        </div>
      </div>
    </div>
  );
}
