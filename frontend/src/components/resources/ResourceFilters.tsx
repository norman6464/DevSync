import { useTranslation } from 'react-i18next';
import type { ResourceCategory, ResourceDifficulty } from '../../types/resource';
import { selectClass } from '../../constants/styles';
import { SearchInput } from '../common';

const categories: ResourceCategory[] = ['book', 'video', 'article', 'course', 'tutorial', 'podcast', 'tool', 'other'];
const difficulties: ResourceDifficulty[] = ['beginner', 'intermediate', 'advanced'];

interface ResourceFiltersProps {
  searchQuery: string;
  onSearchChange: (value: string) => void;
  onSearch: () => void;
  categoryFilter: ResourceCategory | '';
  onCategoryChange: (value: ResourceCategory | '') => void;
  difficultyFilter: ResourceDifficulty | '';
  onDifficultyChange: (value: ResourceDifficulty | '') => void;
}

export default function ResourceFilters({
  searchQuery,
  onSearchChange,
  onSearch,
  categoryFilter,
  onCategoryChange,
  difficultyFilter,
  onDifficultyChange,
}: ResourceFiltersProps) {
  const { t } = useTranslation();

  return (
    <div className="flex flex-wrap gap-4 mb-6">
      <div className="flex-1 min-w-[200px]">
        <SearchInput
          value={searchQuery}
          onChange={onSearchChange}
          onSearch={onSearch}
          placeholder={t('resources.searchPlaceholder')}
          showButton
        />
      </div>
      <select
        value={categoryFilter}
        onChange={(e) => onCategoryChange(e.target.value as ResourceCategory | '')}
        className={selectClass}
      >
        <option value="">{t('resources.allCategories')}</option>
        {categories.map(cat => (
          <option key={cat} value={cat}>
            {t(`resources.categories.${cat}`)}
          </option>
        ))}
      </select>
      <select
        value={difficultyFilter}
        onChange={(e) => onDifficultyChange(e.target.value as ResourceDifficulty | '')}
        className={selectClass}
      >
        <option value="">{t('resources.allDifficulties')}</option>
        {difficulties.map(diff => (
          <option key={diff} value={diff}>
            {t(`resources.difficulty.${diff}`)}
          </option>
        ))}
      </select>
    </div>
  );
}
