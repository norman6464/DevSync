import { useTranslation } from 'react-i18next';
import { Star, Timer, Zap, Pencil, Trash2 } from 'lucide-react';
import type { LearningLog } from '../../types/learningLog';
import { formatDate } from '../../utils/timeFormat';
import { panelClass, badgeBaseClass, editIconButtonClass, deleteIconButtonLargeClass } from '../../constants/styles';
import { getCategoryInfo, getCategoryColor } from '../../constants/learningLogs';

interface LogCardProps {
  log: LearningLog;
  onEdit: (log: LearningLog) => void;
  onDelete: (id: number) => void;
  onToggleFavorite: (id: number) => void;
}

export default function LogCard({ log, onEdit, onDelete, onToggleFavorite }: LogCardProps) {
  const { t } = useTranslation();
  const catInfo = getCategoryInfo(log.category);
  const CatIcon = catInfo.Icon;

  return (
    <div className={panelClass}>
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3 min-w-0 flex-1">
          <CatIcon className="w-5 h-5 text-purple-400 flex-shrink-0 mt-0.5" />
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="font-medium">{log.title}</h3>
              <span className={`${badgeBaseClass} ${getCategoryColor(log.category)}`}>
                {t(catInfo.label)}
              </span>
              {log.source === 'pomodoro' && (
                <span className={`${badgeBaseClass} inline-flex items-center gap-0.5 text-red-400 bg-red-400/10`}>
                  <Timer className="w-3 h-3" />
                  {t('learningLogs.sourcePomodoro')}
                </span>
              )}
            </div>
            <p className="text-sm text-gray-400 mt-1 whitespace-pre-wrap">{log.content}</p>
            <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
              <span>{formatDate(log.created_at)}</span>
              {log.duration > 0 && (
                <span>
                  {log.duration >= 60
                    ? t('learningLogs.durationHoursMinutes', { hours: Math.floor(log.duration / 60), minutes: log.duration % 60 })
                    : t('learningLogs.durationMinutes', { minutes: log.duration })}
                </span>
              )}
              {log.duration >= 120 && (
                <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-purple-400/10 text-purple-400">
                  <Zap className="w-3 h-3" />
                  {t('learningLogs.deepFocus')}
                </span>
              )}
            </div>
          </div>
        </div>

        <div className="flex items-center gap-1">
          <button
            onClick={() => onToggleFavorite(log.id)}
            className={`p-2 transition-colors ${log.is_favorite ? 'text-yellow-400' : 'text-gray-400 hover:text-yellow-400'}`}
            title={t('learningLogs.toggleFavorite')}
          >
            <Star className={`w-4 h-4 ${log.is_favorite ? 'fill-yellow-400' : ''}`} />
          </button>
          <button
            onClick={() => onEdit(log)}
            className={editIconButtonClass}
            title={t('common.edit')}
          >
            <Pencil className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(log.id)}
            className={deleteIconButtonLargeClass}
            title={t('common.delete')}
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
