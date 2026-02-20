import { useTranslation } from 'react-i18next';
import { Calendar, List, Star, ArrowDownWideNarrow } from 'lucide-react';
import type { LogCategory } from '../../types/learningLog';
import { CATEGORIES } from './LogCard';

type SortBy = 'latest' | 'oldest' | 'duration_desc' | 'duration_asc';

interface LogFiltersBarProps {
  view: 'list' | 'calendar';
  filterCategory: LogCategory | 'all';
  showFavoritesOnly: boolean;
  sortBy: SortBy;
  filterDate: string | null;
  onViewList: () => void;
  onViewCalendar: () => void;
  onToggleFavorites: () => void;
  onFilterCategory: (cat: LogCategory | 'all') => void;
  onSortBy: (sort: SortBy) => void;
  onClearFilterDate: () => void;
}

const SORT_OPTIONS: { value: SortBy; label: string }[] = [
  { value: 'latest', label: 'learningLogs.sortLatest' },
  { value: 'oldest', label: 'learningLogs.sortOldest' },
  { value: 'duration_desc', label: 'learningLogs.sortDurationDesc' },
  { value: 'duration_asc', label: 'learningLogs.sortDurationAsc' },
];

export default function LogFiltersBar({
  view, filterCategory, showFavoritesOnly, sortBy, filterDate,
  onViewList, onViewCalendar, onToggleFavorites, onFilterCategory, onSortBy, onClearFilterDate,
}: LogFiltersBarProps) {
  const { t } = useTranslation();

  return (
    <>
      {/* View Toggle */}
      <div className="flex items-center gap-2">
        <button
          onClick={onViewList}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            view === 'list' ? 'bg-purple-500/20 text-purple-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <List className="w-4 h-4" />
          {t('learningLogs.list')}
        </button>
        <button
          onClick={onViewCalendar}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            view === 'calendar' ? 'bg-purple-500/20 text-purple-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <Calendar className="w-4 h-4" />
          {t('learningLogs.calendar')}
        </button>
        <button
          onClick={onToggleFavorites}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            showFavoritesOnly ? 'bg-yellow-500/20 text-yellow-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <Star className={`w-4 h-4 ${showFavoritesOnly ? 'fill-yellow-400' : ''}`} />
          {t('learningLogs.favorites')}
        </button>
        {filterDate && (
          <div className="flex items-center gap-2 ml-2">
            <span className="text-sm text-purple-400">
              {t('learningLogs.logsOnDate', { date: filterDate })}
            </span>
            <button
              onClick={onClearFilterDate}
              className="text-xs text-gray-500 hover:text-white"
            >
              &times;
            </button>
          </div>
        )}
      </div>

      {/* Category Filter */}
      <div className="flex flex-col gap-2">
        <span className="text-sm text-gray-400">{t('learningLogs.category')}</span>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => onFilterCategory('all')}
            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              filterCategory === 'all'
                ? 'bg-purple-500/20 text-purple-400'
                : 'bg-gray-800/50 text-gray-400 hover:text-white'
            }`}
          >
            {t('learningLogs.filterAll')}
          </button>
          {CATEGORIES.map(({ value, label, Icon }) => (
            <button
              key={value}
              onClick={() => onFilterCategory(value)}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                filterCategory === value
                  ? 'bg-purple-500/20 text-purple-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              <Icon className="w-4 h-4" />
              {t(label)}
            </button>
          ))}
        </div>
      </div>

      {/* Sort */}
      <div className="flex flex-col gap-2">
        <span className="text-sm text-gray-400 flex items-center gap-1.5">
          <ArrowDownWideNarrow className="w-4 h-4" />
          {t('learningLogs.sort')}
        </span>
        <div className="flex flex-wrap gap-2">
          {SORT_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              onClick={() => onSortBy(opt.value)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                sortBy === opt.value
                  ? 'bg-purple-500/20 text-purple-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              {t(opt.label)}
            </button>
          ))}
        </div>
      </div>
    </>
  );
}
