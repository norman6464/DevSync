import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { textareaClass, buttonSecondaryClass } from '../../constants/styles';
import { useSubmitShortcut } from '../../hooks/useSubmitShortcut';

interface AnswerFormProps {
  initialBody?: string;
  onSubmit: (body: string) => Promise<boolean>;
  onCancel?: () => void;
  loading?: boolean;
  isEdit?: boolean;
}

export default function AnswerForm({ initialBody = '', onSubmit, onCancel, loading, isEdit }: AnswerFormProps) {
  const { t } = useTranslation();
  const [body, setBody] = useState(initialBody);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await onSubmit(body);
    if (success && !isEdit) {
      setBody('');
    }
  };

  const handleKeyDown = useSubmitShortcut(() => {
    onSubmit(body).then((success) => {
      if (success && !isEdit) setBody('');
    });
  }, !!body.trim() && !loading);

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        onKeyDown={handleKeyDown}
        required
        rows={isEdit ? 4 : 6}
        maxLength={5000}
        className={textareaClass}
        placeholder={t('qa.answerPlaceholder')}
      />
      <p className="text-xs text-gray-500">{t('qa.submitShortcut')}</p>
      <div className="flex gap-3 justify-end">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className={buttonSecondaryClass}
          >
            {t('common.cancel')}
          </button>
        )}
        <button
          type="submit"
          disabled={loading || !body.trim()}
          className="px-6 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
        >
          {loading
            ? t('common.saving')
            : isEdit
            ? t('qa.updateAnswer')
            : t('qa.postAnswer')
          }
        </button>
      </div>
    </form>
  );
}
