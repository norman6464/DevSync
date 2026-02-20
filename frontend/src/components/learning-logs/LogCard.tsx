import { useTranslation } from 'react-i18next';
import { Code, BookOpen, GraduationCap, Users, FileText, Star, type LucideIcon } from 'lucide-react';
import type { LearningLog, LogCategory } from '../../types/learningLog';
import { formatDate } from '../../utils/timeFormat';
import { panelClass } from '../../constants/styles';

export const CATEGORIES: { value: LogCategory; label: string; Icon: LucideIcon }[] = [
  { value: 'coding', label: 'learningLogs.categoryCoding', Icon: Code },
  { value: 'reading', label: 'learningLogs.categoryReading', Icon: BookOpen },
  { value: 'course', label: 'learningLogs.categoryCourse', Icon: GraduationCap },
  { value: 'meetup', label: 'learningLogs.categoryMeetup', Icon: Users },
  { value: 'other', label: 'learningLogs.categoryOther', Icon: FileText },
];

export const getCategoryInfo = (cat: LogCategory) =>
  CATEGORIES.find((c) => c.value === cat) || CATEGORIES[4];

export const getCategoryColor = (cat: LogCategory) => {
  switch (cat) {
    case 'coding': return 'text-blue-400 bg-blue-400/10';
    case 'reading': return 'text-green-400 bg-green-400/10';
    case 'course': return 'text-purple-400 bg-purple-400/10';
    case 'meetup': return 'text-orange-400 bg-orange-400/10';
    default: return 'text-gray-400 bg-gray-400/10';
  }
};

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
              <span className={`px-2 py-0.5 text-xs rounded-full ${getCategoryColor(log.category)}`}>
                {t(catInfo.label)}
              </span>
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
            className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
            title={t('common.edit')}
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
            </svg>
          </button>
          <button
            onClick={() => onDelete(log.id)}
            className="p-2 text-gray-400 hover:text-red-400 transition-colors"
            title={t('common.delete')}
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
