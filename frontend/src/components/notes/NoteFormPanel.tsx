import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { buttonPrimaryClass, buttonSecondaryClass } from '../../constants/styles';

interface NoteFormPanelProps {
  editingNote: boolean;
  title: string;
  setTitle: (v: string) => void;
  content: string;
  setContent: (v: string) => void;
  tags: string;
  setTags: (v: string) => void;
  saving: boolean;
  onSubmit: (e: React.FormEvent) => void;
  onCancel: () => void;
}

export default function NoteFormPanel({
  editingNote,
  title,
  setTitle,
  content,
  setContent,
  tags,
  setTags,
  saving,
  onSubmit,
  onCancel,
}: NoteFormPanelProps) {
  const { t } = useTranslation();
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTitle(e.target.value), [setTitle]);
  const handleContentChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setContent(e.target.value), [setContent]);
  const handleTagsChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTags(e.target.value), [setTags]);

  return (
    <div className="mb-8 p-6 bg-gray-800 border border-gray-700 rounded-md">
      <h2 className="text-xl font-semibold mb-4">
        {editingNote ? t('notes.editNote') : t('notes.createNote')}
      </h2>
      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label htmlFor="note-title" className="block text-sm font-medium mb-2">{t('notes.noteTitle')}</label>
          <input
            id="note-title"
            type="text"
            value={title}
            onChange={handleTitleChange}
            maxLength={200}
            className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>
        <div>
          <label htmlFor="note-content" className="block text-sm font-medium mb-2">{t('notes.noteContent')}</label>
          <textarea
            id="note-content"
            value={content}
            onChange={handleContentChange}
            rows={8}
            maxLength={10000}
            className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
          <p className="text-xs text-gray-500 text-right mt-1">{content.length}/10000</p>
        </div>
        <div>
          <label htmlFor="note-tags" className="block text-sm font-medium mb-2">{t('notes.tags')}</label>
          <input
            id="note-tags"
            type="text"
            value={tags}
            onChange={handleTagsChange}
            maxLength={500}
            placeholder={t('notes.tagsPlaceholder')}
            className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
        <div className="flex gap-3">
          <button
            type="submit"
            disabled={saving}
            className={`${buttonPrimaryClass} disabled:opacity-50`}
          >
            {saving ? t('common.saving') : t('common.save')}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className={buttonSecondaryClass}
          >
            {t('common.cancel')}
          </button>
        </div>
      </form>
    </div>
  );
}
