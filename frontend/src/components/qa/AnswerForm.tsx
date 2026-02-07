import { useState } from 'react';
import { useTranslation } from 'react-i18next';

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

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        required
        rows={isEdit ? 4 : 6}
        className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent resize-none"
        placeholder={t('qa.answerPlaceholder')}
      />
      <div className="flex gap-3 justify-end">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            {t('common.cancel')}
          </button>
        )}
        <button
          type="submit"
          disabled={loading || !body.trim()}
          className="px-6 py-2 bg-green-600 hover:bg-green-500 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
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
