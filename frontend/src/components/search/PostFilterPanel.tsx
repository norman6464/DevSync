import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { PostSearchFilters } from '../../api/posts';

interface PostFilterPanelProps {
  filters: PostSearchFilters;
  onFiltersChange: (f: PostSearchFilters) => void;
}

export default function PostFilterPanel({ filters, onFiltersChange }: PostFilterPanelProps) {
  const { t } = useTranslation();
  const [tagInput, setTagInput] = useState('');

  const handleSortChange = (sortBy: PostSearchFilters['sortBy']) => {
    onFiltersChange({ ...filters, sortBy });
  };

  const handleAddTag = () => {
    const tag = tagInput.trim();
    if (tag && !filters.tags?.includes(tag)) {
      onFiltersChange({ ...filters, tags: [...(filters.tags || []), tag] });
      setTagInput('');
    }
  };

  const handleRemoveTag = (tag: string) => {
    onFiltersChange({ ...filters, tags: filters.tags?.filter((t) => t !== tag) });
  };

  const handleDateFromChange = (value: string) => {
    onFiltersChange({ ...filters, dateFrom: value || undefined });
  };

  const handleDateToChange = (value: string) => {
    onFiltersChange({ ...filters, dateTo: value || undefined });
  };

  const sortOptions: { value: PostSearchFilters['sortBy']; label: string }[] = [
    { value: 'latest', label: t('search.sortLatest') },
    { value: 'popular', label: t('search.sortPopular') },
    { value: 'views', label: t('search.sortViews') },
  ];

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg p-4 space-y-4">
      {/* ソート順 */}
      <div>
        <label className="block text-xs font-medium text-gray-400 mb-2">{t('search.sortBy')}</label>
        <div className="flex gap-2">
          {sortOptions.map((opt) => (
            <button
              key={opt.value}
              onClick={() => handleSortChange(opt.value)}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                filters.sortBy === opt.value
                  ? 'bg-blue-500/20 text-blue-400 border border-blue-500/50'
                  : 'bg-gray-800 text-gray-400 border border-gray-700 hover:border-gray-600'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      {/* タグフィルター */}
      <div>
        <label className="block text-xs font-medium text-gray-400 mb-2">{t('search.tagFilter')}</label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleAddTag()}
            placeholder={t('search.tagPlaceholder')}
            className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <button
            onClick={handleAddTag}
            className="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-300 hover:border-gray-600 transition-colors"
          >
            {t('common.add')}
          </button>
        </div>
        {filters.tags && filters.tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {filters.tags.map((tag) => (
              <span
                key={tag}
                className="inline-flex items-center gap-1 px-2 py-1 bg-blue-500/10 text-blue-400 border border-blue-500/30 rounded-full text-xs"
              >
                #{tag}
                <button
                  onClick={() => handleRemoveTag(tag)}
                  className="hover:text-blue-300 ml-0.5"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
      </div>

      {/* 日付範囲フィルター */}
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-1">{t('search.dateFrom')}</label>
          <input
            type="date"
            value={filters.dateFrom || ''}
            onChange={(e) => handleDateFromChange(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 focus:outline-none focus:border-blue-500"
          />
        </div>
        <div>
          <label className="block text-xs font-medium text-gray-400 mb-1">{t('search.dateTo')}</label>
          <input
            type="date"
            value={filters.dateTo || ''}
            onChange={(e) => handleDateToChange(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-100 focus:outline-none focus:border-blue-500"
          />
        </div>
      </div>
    </div>
  );
}
