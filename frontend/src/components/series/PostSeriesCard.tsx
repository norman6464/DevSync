import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import type { PostSeries } from '../../types/post';
import { cardPaddedClass } from '../../constants/styles';

interface PostSeriesCardProps {
  series: PostSeries;
  onEdit?: () => void;
  onDelete?: () => void;
  isOwner?: boolean;
}

export default function PostSeriesCard({ series, onEdit, onDelete, isOwner }: PostSeriesCardProps) {
  const { t } = useTranslation();

  return (
    <div className={cardPaddedClass}>
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1">
          <Link
            to={`/series/${series.id}`}
            className="text-lg font-semibold text-white hover:text-blue-400 transition-colors"
          >
            {series.title}
          </Link>
        </div>

        {isOwner && (
          <div className="flex gap-1">
            <button
              onClick={onEdit}
              className="p-1.5 text-gray-400 hover:text-white transition-colors"
              aria-label={t('common.edit')}
              title={t('common.edit')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
              </svg>
            </button>
            <button
              onClick={onDelete}
              className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
              aria-label={t('common.delete')}
              title={t('common.delete')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
              </svg>
            </button>
          </div>
        )}
      </div>

      {series.description && (
        <p className="text-gray-400 text-sm mt-2 line-clamp-2">{series.description}</p>
      )}

      <p className="text-xs text-gray-500 mt-3">
        {new Date(series.created_at).toLocaleDateString()}
      </p>
    </div>
  );
}
