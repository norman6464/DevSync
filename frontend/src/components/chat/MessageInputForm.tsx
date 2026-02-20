import { useTranslation } from 'react-i18next';
import { Send } from 'lucide-react';
import { messageInputClass } from '../../constants/styles';

interface MessageInputFormProps {
  value: string;
  onChange: (value: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  placeholder?: string;
}

export default function MessageInputForm({
  value,
  onChange,
  onSubmit,
  placeholder,
}: MessageInputFormProps) {
  const { t } = useTranslation();

  return (
    <form onSubmit={onSubmit} className="p-4 border-t border-gray-800 flex gap-3">
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className={messageInputClass}
      />
      <button
        type="submit"
        disabled={!value.trim()}
        className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors flex items-center gap-2"
      >
        <Send className="w-4 h-4" />
        {t('chat.send')}
      </button>
    </form>
  );
}
