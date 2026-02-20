import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { TrendingUp, CheckCircle2, Clock } from 'lucide-react';
import { panelClass } from '../../constants/styles';

interface QuickStatsWidgetProps {
  activeCount: number;
  completedCount: number;
}

export default function QuickStatsWidget({ activeCount, completedCount }: QuickStatsWidgetProps) {
  const { t } = useTranslation();

  return (
    <div className={panelClass}>
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
        <TrendingUp className="w-4 h-4 text-green-400" aria-hidden="true" />
        {t('dashboard.quickStats')}
      </h3>
      <div className="space-y-2">
        <Link
          to="/goals"
          className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
        >
          <CheckCircle2 className="w-4 h-4 text-green-400" aria-hidden="true" />
          <span className="text-xs text-gray-300 flex-1">{t('dashboard.goalsCompleted')}</span>
          <span className="text-xs font-medium text-white">{completedCount}</span>
        </Link>
        <Link
          to="/goals"
          className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
        >
          <Clock className="w-4 h-4 text-blue-400" aria-hidden="true" />
          <span className="text-xs text-gray-300 flex-1">{t('dashboard.goalsInProgress')}</span>
          <span className="text-xs font-medium text-white">{activeCount}</span>
        </Link>
      </div>
    </div>
  );
}
