import { useTranslation } from 'react-i18next';
import { Search } from 'lucide-react';
import { buttonSecondaryClass, searchInputClass } from '../../constants/styles';

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
          className={searchInputClass}
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
