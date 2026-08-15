import { useTranslation } from 'react-i18next';
import { Filter } from 'lucide-react';
import type { NotificationType } from '../../types/notification';

const FILTER_TYPES: { key: NotificationType | ''; labelKey: string }[] = [
  { key: '', labelKey: 'notifications.filterAll' },
  { key: 'post', labelKey: 'notifications.filterPost' },
  { key: 'like', labelKey: 'notifications.filterLike' },
  { key: 'comment', labelKey: 'notifications.filterComment' },
  { key: 'mention', labelKey: 'notifications.filterMention' },
  { key: 'follow', labelKey: 'notifications.filterFollow' },
  { key: 'message', labelKey: 'notifications.filterMessage' },
  { key: 'answer', labelKey: 'notifications.filterAnswer' },
  { key: 'badge', labelKey: 'notifications.filterBadge' },
];

interface NotificationFiltersProps {
  filterType: NotificationType | '';
  setFilterType: (type: NotificationType | '') => void;
  showUnreadOnly: boolean;
  onToggleUnreadOnly: () => void;
}

export default function NotificationFilters({
  filterType, setFilterType, showUnreadOnly, onToggleUnreadOnly,
}: NotificationFiltersProps) {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col gap-4 mb-6">
      <div className="flex flex-wrap gap-2" role="group" aria-label={t('notifications.filterGroup')}>
        {FILTER_TYPES.map(({ key, labelKey }) => (
          <button
            key={key}
            onClick={() => setFilterType(key)}
            aria-pressed={filterType === key}
            className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              filterType === key
                ? 'bg-gray-700 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            {t(labelKey)}
          </button>
        ))}
      </div>

      <button
        onClick={onToggleUnreadOnly}
        aria-pressed={showUnreadOnly}
        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors self-start ${
          showUnreadOnly
            ? 'bg-blue-600 text-white'
            : 'bg-gray-800 text-gray-400 hover:text-white'
        }`}
      >
        <Filter className="w-4 h-4" aria-hidden="true" />
        {t('notifications.showUnreadOnly')}
      </button>
    </div>
  );
}
