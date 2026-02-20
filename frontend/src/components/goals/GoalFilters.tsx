import { useTranslation } from 'react-i18next';
import { Filter, ArrowUpDown } from 'lucide-react';
import { type GoalCategory } from '../../api/goals';
import { CATEGORIES } from './GoalCard';

interface GoalFiltersProps {
  filterStatus: 'all' | 'active' | 'paused' | 'completed';
  setFilterStatus: (s: 'all' | 'active' | 'paused' | 'completed') => void;
  filterCategory: GoalCategory | 'all';
  setFilterCategory: (c: GoalCategory | 'all') => void;
  sortBy: 'newest' | 'oldest' | 'deadline' | 'progress';
  setSortBy: (s: 'newest' | 'oldest' | 'deadline' | 'progress') => void;
}

const STATUSES = ['all', 'active', 'paused', 'completed'] as const;
const SORTS = ['newest', 'oldest', 'deadline', 'progress'] as const;

export default function GoalFilters({
  filterStatus, setFilterStatus,
  filterCategory, setFilterCategory,
  sortBy, setSortBy,
}: GoalFiltersProps) {
  const { t } = useTranslation();

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4 space-y-3">
      <div className="flex items-center gap-2 text-sm text-gray-400">
        <Filter className="w-4 h-4" />
        <span>{t('goals.filter')}</span>
      </div>
      <div className="flex flex-wrap gap-2">
        <span className="text-xs text-gray-500 self-center mr-1">{t('goals.status')}:</span>
        {STATUSES.map((s) => (
          <button
            key={s}
            onClick={() => setFilterStatus(s)}
            className={`px-3 py-1 text-xs rounded-full border transition-colors ${
              filterStatus === s
                ? 'border-blue-500 bg-blue-500/10 text-blue-400'
                : 'border-gray-700 text-gray-400 hover:border-gray-600'
            }`}
          >
            {s === 'all' ? t('common.all') : t(`goals.status${s.charAt(0).toUpperCase() + s.slice(1)}`)}
          </button>
        ))}
      </div>
      <div className="flex flex-wrap gap-2">
        <span className="text-xs text-gray-500 self-center mr-1">{t('goals.category')}:</span>
        {(['all', ...CATEGORIES.map(c => c.value)] as const).map((c) => (
          <button
            key={c}
            onClick={() => setFilterCategory(c)}
            className={`px-3 py-1 text-xs rounded-full border transition-colors ${
              filterCategory === c
                ? 'border-purple-500 bg-purple-500/10 text-purple-400'
                : 'border-gray-700 text-gray-400 hover:border-gray-600'
            }`}
          >
            {c === 'all' ? t('common.all') : t(`goals.category${c.charAt(0).toUpperCase() + c.slice(1)}`)}
          </button>
        ))}
      </div>
      <div className="flex flex-wrap gap-2">
        <span className="text-xs text-gray-500 self-center mr-1">
          <ArrowUpDown className="w-3 h-3 inline mr-1" />
          {t('goals.sort')}:
        </span>
        {SORTS.map((s) => (
          <button
            key={s}
            onClick={() => setSortBy(s)}
            className={`px-3 py-1 text-xs rounded-full border transition-colors ${
              sortBy === s
                ? 'border-green-500 bg-green-500/10 text-green-400'
                : 'border-gray-700 text-gray-400 hover:border-gray-600'
            }`}
          >
            {t(`goals.sort${s.charAt(0).toUpperCase() + s.slice(1)}`)}
          </button>
        ))}
      </div>
    </div>
  );
}
