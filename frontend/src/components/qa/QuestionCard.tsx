import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { TrendingUp, Clock, CheckCircle2, Pencil, Trash2 } from 'lucide-react';
import type { Question } from '../../types/qa';
import Avatar from '../common/Avatar';
import { cardPaddedClass, iconButtonClass, deleteIconButtonClass } from '../../constants/styles';
import { parseJsonArray } from '../../utils/json';
import { formatDate } from '../../utils/timeFormat';
import { isWithinLast } from '../../utils/timeFormat';

interface QuestionCardProps {
  question: Question;
  isOwner?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

export default function QuestionCard({ question, isOwner = false, onEdit, onDelete }: QuestionCardProps) {
  const { t } = useTranslation();

  const tags = parseJsonArray(question.tags);
  const isNew = isWithinLast(question.created_at, 24 * 60 * 60 * 1000);

  const statusBorderClass = question.is_solved
    ? 'border-l-4 border-l-green-500'
    : question.answer_count > 0
    ? 'border-l-4 border-l-yellow-500'
    : '';

  return (
    <div className={`${cardPaddedClass} ${statusBorderClass}`}>
      <div className="flex gap-4">
        {/* Vote & Answer counts */}
        <div className="flex flex-col items-center gap-2 text-center min-w-[60px]">
          <div className={`text-sm ${question.vote_count > 0 ? 'text-green-400' : question.vote_count < 0 ? 'text-red-400' : 'text-gray-400'}`}>
            <span className="block text-lg font-bold">{question.vote_count}</span>
            <span className="text-xs">{t('qa.votes')}</span>
          </div>
          <div className={`px-2 py-1 rounded text-xs font-medium ${
            question.is_solved
              ? 'bg-green-600/20 text-green-400 border border-green-600/30'
              : question.answer_count > 0
              ? 'bg-gray-700 text-gray-300'
              : 'text-gray-500'
          }`}>
            <span className="block text-sm font-bold">{question.answer_count}</span>
            {t('qa.answers')}
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <Link
              to={`/qa/${question.id}`}
              className="text-lg font-semibold text-white hover:text-blue-400 transition-colors line-clamp-2"
            >
              {question.is_solved && (
                <span className="inline-flex items-center mr-2 text-green-400">
                  <CheckCircle2 className="w-4 h-4" aria-hidden="true" />
                </span>
              )}
              {question.title}
            </Link>

            {isOwner && (
              <div className="flex gap-1 flex-shrink-0">
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
              </div>
            )}
          </div>

          <p className="text-gray-400 text-sm mt-1 line-clamp-2">{question.body}</p>

          {/* Tags & Badges */}
          {(tags.length > 0 || question.vote_count >= 5 || isNew) && (
            <div className="flex flex-wrap gap-1.5 mt-2">
              {isNew && (
                <span className="inline-flex items-center gap-0.5 px-2 py-0.5 bg-cyan-400/10 text-cyan-400 text-xs rounded-md font-medium">
                  <Clock className="w-3 h-3" />
                  {t('qa.newQuestion')}
                </span>
              )}
              {question.vote_count >= 5 && (
                <span className="inline-flex items-center gap-0.5 px-2 py-0.5 bg-yellow-400/10 text-yellow-400 text-xs rounded-md font-medium">
                  <TrendingUp className="w-3 h-3" />
                  {t('qa.popularBadge')}
                </span>
              )}
              {tags.map(tag => (
                <span key={tag} className="px-2 py-0.5 bg-blue-600/20 text-blue-400 text-xs rounded-md">
                  {tag}
                </span>
              ))}
            </div>
          )}

          {/* Footer */}
          <div className="flex items-center justify-between mt-3 pt-2 border-t border-gray-700/50">
            {question.user && (
              <Link
                to={`/profile/${question.user?.username || question.user_id}`}
                className="flex items-center gap-2 hover:opacity-80 transition-opacity"
              >
                <Avatar avatarUrl={question.user.avatar_url} name={question.user.name} size="sm" />
                <span className="text-sm text-gray-400">{question.user.name}</span>
              </Link>
            )}
            <span className="text-xs text-gray-500">
              {formatDate(question.created_at)}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
