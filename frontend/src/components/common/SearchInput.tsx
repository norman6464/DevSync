import { useTranslation } from 'react-i18next';
import { Search } from 'lucide-react';
import { buttonSecondaryClass } from '../../constants/styles';

interface SearchInputProps {
  value: string;
  onChange: (value: string) => void;
  onSearch?: () => void;
  placeholder?: string;
  showButton?: boolean;
}

export default function SearchInput({ value, onChange, onSearch, placeholder, showButton = false }: SearchInputProps) {
  const { t } = useTranslation();

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && onSearch) {
      onSearch();
    }
  };

  return (
    <div className="flex gap-2">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          className="w-full pl-10 pr-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
        />
      </div>
      {showButton && onSearch && (
        <button
          type="button"
          onClick={onSearch}
          className={buttonSecondaryClass}
        >
          {t('common.search')}
        </button>
      )}
    </div>
  );
}
