import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Question, CreateQuestionRequest } from '../../types/qa';

interface QuestionFormProps {
  question?: Question;
  onSubmit: (data: CreateQuestionRequest) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

export default function QuestionForm({ question, onSubmit, onCancel, loading }: QuestionFormProps) {
  const { t } = useTranslation();

  const parseTags = (tags: string | undefined): string => {
    if (!tags) return '';
    try {
      return (JSON.parse(tags) as string[]).join(', ');
    } catch {
      return '';
    }
  };

  const [title, setTitle] = useState(question?.title || '');
  const [body, setBody] = useState(question?.body || '');
  const [tagsInput, setTagsInput] = useState(parseTags(question?.tags));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const tags = tagsInput.trim()
      ? JSON.stringify(tagsInput.split(',').map(t => t.trim()).filter(Boolean))
      : '';
    await onSubmit({ title, body, tags });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Title */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('qa.questionTitle')} *
        </label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={500}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent"
          placeholder={t('qa.questionTitlePlaceholder')}
        />
      </div>

      {/* Body */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('qa.questionBody')} *
        </label>
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          required
          rows={8}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent resize-none"
          placeholder={t('qa.questionBodyPlaceholder')}
        />
      </div>

      {/* Tags */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('qa.tags')}
        </label>
        <input
          type="text"
          value={tagsInput}
          onChange={(e) => setTagsInput(e.target.value)}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent"
          placeholder={t('qa.tagsPlaceholder')}
        />
        <p className="text-xs text-gray-500 mt-1">{t('qa.tagsHint')}</p>
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim() || !body.trim()}
          className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
