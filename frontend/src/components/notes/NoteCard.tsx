import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Star, Edit, Trash2, FileText, BookOpen, Sparkles } from 'lucide-react';
import type { Note } from '../../api/notes';
import { formatDistanceToNow } from '../../utils/timeFormat';
import { linkSmallClass } from '../../constants/styles';
import { isWithinLast } from '../../utils/timeFormat';

interface NoteCardProps {
  note: Note;
  onToggleFavorite: (id: number) => void;
  onEdit: (note: Note) => void;
  onDelete: (id: number) => void;
}

export default function NoteCard({ note, onToggleFavorite, onEdit, onDelete }: NoteCardProps) {
  const { t } = useTranslation();
  const [expanded, setExpanded] = useState(false);
  const isLong = note.content.length > 100;
  const isNew = isWithinLast(note.created_at, 24 * 60 * 60 * 1000);

  return (
    <div className="p-6 bg-gray-800 border border-gray-700 rounded-md hover:border-gray-600 transition-colors">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-xl font-semibold">{note.title}</h3>
            {note.is_favorite && (
              <Star className="w-5 h-5 fill-yellow-500 text-yellow-500" />
            )}
            {isNew && (
              <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-emerald-400/10 text-emerald-400 font-medium">
                <Sparkles className="w-3 h-3" />
                {t('notes.new')}
              </span>
            )}
          </div>
          <p className={`text-gray-400 mb-1 whitespace-pre-wrap ${!expanded ? 'line-clamp-2' : ''}`}>{note.content}</p>
          {isLong && (
            <button
              onClick={() => setExpanded(!expanded)}
              className={`${linkSmallClass} mb-2`}
            >
              {expanded ? t('notes.collapse') : t('notes.expand')}
            </button>
          )}
          {!isLong && <div className="mb-2" />}
          {note.tags && (
            <div className="flex flex-wrap gap-2 mb-3">
              {note.tags.split(',').map((tag, idx) => (
                <span
                  key={idx}
                  className="px-2 py-1 text-xs bg-gray-700 text-gray-300 rounded"
                >
                  {tag.trim()}
                </span>
              ))}
            </div>
          )}
          <div className="flex items-center gap-3 text-sm text-gray-500">
            <span>{t('notes.lastUpdated')}: {formatDistanceToNow(note.updated_at)}</span>
            <span className="text-gray-600">·</span>
            <span>{note.content.length.toLocaleString()} {t('notes.chars')}</span>
            {note.content.length >= 500 ? (
              <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-purple-400/10 text-purple-400">
                <BookOpen className="w-3 h-3" />
                {t('notes.volumeDetailed')}
              </span>
            ) : note.content.length >= 200 ? (
              <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-yellow-400/10 text-yellow-400">
                <FileText className="w-3 h-3" />
                {t('notes.volumeRegular')}
              </span>
            ) : null}
          </div>
        </div>
        <div className="flex gap-2 ml-4">
          <button
            onClick={() => onToggleFavorite(note.id)}
            className="p-2 text-gray-400 hover:text-yellow-500 transition-colors"
            aria-label={t('notes.favorites')}
          >
            <Star aria-hidden="true" className={`w-5 h-5 ${note.is_favorite ? 'fill-yellow-500 text-yellow-500' : ''}`} />
          </button>
          <button
            onClick={() => onEdit(note)}
            className="p-2 text-gray-400 hover:text-blue-500 transition-colors"
            aria-label={t('common.edit')}
          >
            <Edit aria-hidden="true" className="w-5 h-5" />
          </button>
          <button
            onClick={() => onDelete(note.id)}
            className="p-2 text-gray-400 hover:text-red-500 transition-colors"
            aria-label={t('common.delete')}
          >
            <Trash2 aria-hidden="true" className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
