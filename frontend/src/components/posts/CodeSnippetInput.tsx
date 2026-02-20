import { useTranslation } from 'react-i18next';
import { X } from 'lucide-react';
import LanguageSelector from './LanguageSelector';

export interface SnippetInputData {
  language: string;
  file_name: string;
  code: string;
}

interface CodeSnippetInputProps {
  value: SnippetInputData;
  onChange: (value: SnippetInputData) => void;
  onRemove: () => void;
  index: number;
}

export default function CodeSnippetInput({ value, onChange, onRemove, index }: CodeSnippetInputProps) {
  const { t } = useTranslation();

  return (
    <div className="border border-gray-700 rounded-lg overflow-hidden bg-gray-800/50">
      {/* Header */}
      <div className="flex items-center gap-2 px-3 py-2 border-b border-gray-700 bg-gray-800">
        <span className="text-xs text-gray-400 font-mono">#{index + 1}</span>
        <LanguageSelector
          value={value.language}
          onChange={(language) => onChange({ ...value, language })}
        />
        <input
          type="text"
          value={value.file_name}
          onChange={(e) => onChange({ ...value, file_name: e.target.value })}
          placeholder={t('post.fileName')}
          maxLength={255}
          className="flex-1 px-3 py-1.5 bg-gray-900/50 border border-gray-700 rounded text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
        <button
          type="button"
          onClick={onRemove}
          className="p-1.5 text-gray-400 hover:text-red-400 hover:bg-gray-700 rounded transition-colors"
          aria-label={t('post.removeSnippet')}
        >
          <X className="w-4 h-4" aria-hidden="true" />
        </button>
      </div>

      {/* Code textarea */}
      <textarea
        value={value.code}
        onChange={(e) => onChange({ ...value, code: e.target.value })}
        placeholder={t('post.codePlaceholder')}
        className="w-full p-4 bg-transparent text-white resize-none focus:outline-none font-mono text-sm"
        style={{ minHeight: '120px' }}
      />
    </div>
  );
}
