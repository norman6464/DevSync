import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { Question } from '../../types/qa';
import Avatar from '../common/Avatar';

interface QuestionCardProps {
  question: Question;
  isOwner?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

export default function QuestionCard({ question, isOwner = false, onEdit, onDelete }: QuestionCardProps) {
  const { t } = useTranslation();

  const tags: string[] = question.tags ? (() => {
    try { return JSON.parse(question.tags); } catch { return []; }
  })() : [];

  return (
    <div className="bg-gray-800 rounded-xl p-5 border border-gray-700 hover:border-gray-600 transition-colors">
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
                  <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                </span>
              )}
              {question.title}
            </Link>

            {isOwner && (
              <div className="flex gap-1 flex-shrink-0">
                <button
                  onClick={onEdit}
                  className="p-1.5 text-gray-400 hover:text-white transition-colors"
                  title={t('common.edit')}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                  </svg>
                </button>
                <button
                  onClick={onDelete}
                  className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
                  title={t('common.delete')}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                  </svg>
                </button>
              </div>
            )}
          </div>

          <p className="text-gray-400 text-sm mt-1 line-clamp-2">{question.body}</p>

          {/* Tags */}
          {tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5 mt-2">
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
                to={`/profile/${question.user_id}`}
                className="flex items-center gap-2 hover:opacity-80 transition-opacity"
              >
                <Avatar avatarUrl={question.user.avatar_url} name={question.user.name} size="sm" />
                <span className="text-sm text-gray-400">{question.user.name}</span>
              </Link>
            )}
            <span className="text-xs text-gray-500">
              {new Date(question.created_at).toLocaleDateString()}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
