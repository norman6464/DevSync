import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Question, CreateQuestionRequest } from '../../types/qa';
import { buttonSecondaryClass, inputClass, labelClass, textareaClass } from '../../constants/styles';

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
        <label htmlFor="question-title" className={labelClass}>
          {t('qa.questionTitle')} *
        </label>
        <input
          id="question-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={500}
          className={inputClass}
          placeholder={t('qa.questionTitlePlaceholder')}
        />
      </div>

      {/* Body */}
      <div>
        <label htmlFor="question-body" className={labelClass}>
          {t('qa.questionBody')} *
        </label>
        <textarea
          id="question-body"
          value={body}
          onChange={(e) => setBody(e.target.value)}
          required
          rows={8}
          maxLength={5000}
          className={textareaClass}
          placeholder={t('qa.questionBodyPlaceholder')}
        />
        <p className="text-xs text-gray-500 text-right mt-1">{body.length}/5000</p>
      </div>

      {/* Tags */}
      <div>
        <label htmlFor="question-tags" className={labelClass}>
          {t('qa.tags')}
        </label>
        <input
          id="question-tags"
          type="text"
          value={tagsInput}
          onChange={(e) => setTagsInput(e.target.value)}
          maxLength={500}
          className={inputClass}
          placeholder={t('qa.tagsPlaceholder')}
        />
        <p className="text-xs text-gray-500 mt-1">{t('qa.tagsHint')}</p>
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className={`flex-1 ${buttonSecondaryClass}`}
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim() || !body.trim()}
          className={`flex-1 ${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
