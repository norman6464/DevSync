import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Pencil, Trash2 } from 'lucide-react';
import type { PostSeries } from '../../types/post';
import { cardPaddedClass, iconButtonClass, deleteIconButtonClass } from '../../constants/styles';
import { formatDate } from '../../utils/timeFormat';

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
              className={iconButtonClass}
              aria-label={t('common.edit')}
              title={t('common.edit')}
            >
              <Pencil className="w-4 h-4" />
            </button>
            <button
              onClick={onDelete}
              className={deleteIconButtonClass}
              aria-label={t('common.delete')}
              title={t('common.delete')}
            >
              <Trash2 className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>

      {series.description && (
        <p className="text-gray-400 text-sm mt-2 line-clamp-2">{series.description}</p>
      )}

      <p className="text-xs text-gray-500 mt-3">
        {formatDate(series.created_at)}
      </p>
    </div>
  );
}
