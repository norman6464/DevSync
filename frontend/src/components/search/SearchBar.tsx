import { useTranslation } from 'react-i18next';
import { Search } from 'lucide-react';
import { searchInputClass } from '../../constants/styles';

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  onSearch: () => void;
  placeholder?: string;
}

export default function SearchBar({ value, onChange, onSearch, placeholder }: SearchBarProps) {
  const { t } = useTranslation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch();
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <div className="flex-1 relative">
        <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder || t('search.placeholder')}
          className={searchInputClass}
        />
      </div>
      <button
        type="submit"
        className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium text-sm transition-colors"
      >
        {t('common.search')}
      </button>
    </form>
  );
}
