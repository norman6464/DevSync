import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import type { LogCategory } from '../../types/learningLog';
import type { LearningLog } from '../../types/learningLog';
import { CATEGORIES } from './LogCard';
import { Modal } from '../common';
import { inputClass, buttonSecondaryClass, labelClass } from '../../constants/styles';
import { MAX_LENGTH } from '../../utils/formValidation';
import { CharCount } from '../common';

interface LogFormModalProps {
  isOpen: boolean;
  editingLog: LearningLog | null;
  title: string;
  setTitle: (v: string) => void;
  content: string;
  setContent: (v: string) => void;
  category: LogCategory;
  setCategory: (v: LogCategory) => void;
  duration: string;
  setDuration: (v: string) => void;
  saving: boolean;
  onSubmit: (e: React.FormEvent) => void;
  onClose: () => void;
}

export default function LogFormModal({
  isOpen, editingLog, title, setTitle, content, setContent,
  category, setCategory, duration, setDuration, saving, onSubmit, onClose,
}: LogFormModalProps) {
  const { t } = useTranslation();

  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTitle(e.target.value), [setTitle]);
  const handleContentChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setContent(e.target.value), [setContent]);
  const handleCategoryChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => setCategory(e.target.value as LogCategory), [setCategory]);
  const handleDurationChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setDuration(e.target.value), [setDuration]);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={editingLog ? t('learningLogs.editLog') : t('learningLogs.addLog')}
      maxWidth="max-w-md"
    >
      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label htmlFor="log-title" className={labelClass}>
            {t('learningLogs.logTitle')}
          </label>
          <input
            id="log-title"
            type="text"
            value={title}
            onChange={handleTitleChange}
            placeholder={t('learningLogs.titlePlaceholder')}
            maxLength={MAX_LENGTH.logTitle}
            className={inputClass}
            required
          />
        </div>
        <div>
          <label htmlFor="log-content" className={labelClass}>
            {t('learningLogs.content')}
          </label>
          <textarea
            id="log-content"
            value={content}
            onChange={handleContentChange}
            placeholder={t('learningLogs.contentPlaceholder')}
            rows={4}
            maxLength={MAX_LENGTH.logContent}
            className={`${inputClass} resize-none`}
            required
          />
          <CharCount value={content} max={MAX_LENGTH.logContent} />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="log-category" className={labelClass}>
              {t('learningLogs.category')}
            </label>
            <select
              id="log-category"
              value={category}
              onChange={handleCategoryChange}
              className={inputClass}
            >
              {CATEGORIES.map((cat) => (
                <option key={cat.value} value={cat.value}>
                  {t(cat.label)}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="log-duration" className={labelClass}>
              {t('learningLogs.duration')}
            </label>
            <input
              id="log-duration"
              type="number"
              value={duration}
              onChange={handleDurationChange}
              placeholder={t('learningLogs.durationPlaceholder')}
              min="0"
              className={inputClass}
            />
          </div>
        </div>
        <div className="flex gap-3 justify-end pt-2">
          <button
            type="button"
            onClick={onClose}
            className={`${buttonSecondaryClass} text-sm font-medium`}
          >
            {t('common.cancel')}
          </button>
          <button
            type="submit"
            disabled={saving || !title.trim() || !content.trim()}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
          >
            {saving ? t('common.loading') : editingLog ? t('common.save') : t('learningLogs.addLog')}
          </button>
        </div>
      </form>
    </Modal>
  );
}
