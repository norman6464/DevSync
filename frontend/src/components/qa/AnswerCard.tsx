import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Pencil, Trash2 } from 'lucide-react';
import type { Answer } from '../../types/qa';
import Avatar from '../common/Avatar';
import { iconButtonClass, deleteIconButtonClass } from '../../constants/styles';
import { formatDate } from '../../utils/timeFormat';

interface AnswerCardProps {
  answer: Answer;
  isOwner?: boolean;
  isQuestionOwner?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
  onSetBest?: () => void;
  onVote?: (value: 1 | -1) => void;
  onRemoveVote?: () => void;
}

export default function AnswerCard({
  answer,
  isOwner = false,
  isQuestionOwner = false,
  onEdit,
  onDelete,
  onSetBest,
  onVote,
  onRemoveVote,
}: AnswerCardProps) {
  const { t } = useTranslation();
  const [voting, setVoting] = useState(false);

  const handleVote = async (value: 1 | -1) => {
    if (voting) return;
    setVoting(true);
    try {
      await onVote?.(value);
    } finally {
      setVoting(false);
    }
  };

  return (
    <div className={`bg-gray-800 rounded-md p-5 border transition-colors ${
      answer.is_best ? 'border-green-600/50 bg-green-900/10' : 'border-gray-700'
    }`}>
      {answer.is_best && (
        <div className="flex items-center gap-1.5 mb-3 text-green-400 text-sm font-medium">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
          {t('qa.bestAnswer')}
        </div>
      )}

      <div className="flex gap-4">
        {/* Vote buttons */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={() => handleVote(1)}
            disabled={voting}
            className="p-1 text-gray-400 hover:text-green-400 transition-colors disabled:opacity-50"
            title={t('qa.upvote')}
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 15.75l7.5-7.5 7.5 7.5" />
            </svg>
          </button>
          <span className={`text-sm font-bold ${
            answer.vote_count > 0 ? 'text-green-400' : answer.vote_count < 0 ? 'text-red-400' : 'text-gray-400'
          }`}>
            {answer.vote_count}
          </span>
          <button
            onClick={() => handleVote(-1)}
            disabled={voting}
            className="p-1 text-gray-400 hover:text-red-400 transition-colors disabled:opacity-50"
            title={t('qa.downvote')}
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <p className="text-gray-200 whitespace-pre-wrap">{answer.body}</p>

          {/* Footer */}
          <div className="flex items-center justify-between mt-4 pt-3 border-t border-gray-700/50">
            <div className="flex items-center gap-3">
              {answer.user && (
                <Link
                  to={`/profile/${answer.user?.username || answer.user_id}`}
                  className="flex items-center gap-2 hover:opacity-80 transition-opacity"
                >
                  <Avatar avatarUrl={answer.user.avatar_url} name={answer.user.name} size="sm" />
                  <span className="text-sm text-gray-400">{answer.user.name}</span>
                </Link>
              )}
              <span className="text-xs text-gray-500">
                {formatDate(answer.created_at)}
              </span>
            </div>

            <div className="flex items-center gap-2">
              {isQuestionOwner && !answer.is_best && (
                <button
                  onClick={onSetBest}
                  className="px-3 py-1 text-xs border border-green-600/50 text-green-400 hover:bg-green-600/20 rounded-lg transition-colors"
                >
                  {t('qa.markBest')}
                </button>
              )}
              {isOwner && (
                <>
                  <button
                    onClick={onEdit}
                    className={iconButtonClass}
                    title={t('common.edit')}
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button
                    onClick={onDelete}
                    className={deleteIconButtonClass}
                    title={t('common.delete')}
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
