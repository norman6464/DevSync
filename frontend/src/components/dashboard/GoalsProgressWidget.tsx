import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Target, ChevronRight } from 'lucide-react';
import type { LearningGoal } from '../../api/goals';

interface GoalsProgressWidgetProps {
  activeGoals: LearningGoal[];
  completedGoals: LearningGoal[];
  avgProgress: number;
  loading: boolean;
}

export default function GoalsProgressWidget({ activeGoals, completedGoals, avgProgress, loading }: GoalsProgressWidgetProps) {
  const { t } = useTranslation();

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="flex items-center gap-2 text-sm font-medium text-white">
          <Target className="w-4 h-4 text-blue-400" aria-hidden="true" />
          {t('dashboard.goalsProgress')}
        </h3>
        <Link to="/goals" className="text-xs text-gray-400 hover:text-blue-400 transition-colors">
          {t('dashboard.viewAll')}
        </Link>
      </div>

      {loading ? (
        <div className="space-y-3">
          <div className="h-4 bg-gray-800 rounded animate-pulse" />
          <div className="h-4 bg-gray-800 rounded animate-pulse w-2/3" />
        </div>
      ) : activeGoals.length === 0 ? (
        <div className="text-center py-4">
          <p className="text-xs text-gray-500 mb-2">{t('dashboard.noActiveGoals')}</p>
          <Link
            to="/goals"
            className="text-xs text-blue-400 hover:text-blue-300 transition-colors"
          >
            {t('dashboard.createGoal')}
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {/* Stats Row */}
          <div className="grid grid-cols-3 gap-2">
            <div className="bg-gray-800/50 rounded-lg p-2 text-center">
              <div className="text-lg font-bold text-blue-400">{activeGoals.length}</div>
              <div className="text-[10px] text-gray-500">{t('dashboard.active')}</div>
            </div>
            <div className="bg-gray-800/50 rounded-lg p-2 text-center">
              <div className="text-lg font-bold text-green-400">{completedGoals.length}</div>
              <div className="text-[10px] text-gray-500">{t('dashboard.completed')}</div>
            </div>
            <div className="bg-gray-800/50 rounded-lg p-2 text-center">
              <div className="text-lg font-bold text-orange-400">{avgProgress}%</div>
              <div className="text-[10px] text-gray-500">{t('dashboard.avgProgress')}</div>
            </div>
          </div>

          {/* Active Goals List */}
          <div className="space-y-2">
            {activeGoals.slice(0, 3).map((goal) => (
              <div key={goal.id} className="group">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-xs text-gray-300 truncate flex-1 mr-2">{goal.title}</span>
                  <span className="text-xs text-gray-500 shrink-0">{goal.progress}%</span>
                </div>
                <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden" role="progressbar" aria-valuenow={goal.progress} aria-valuemin={0} aria-valuemax={100} aria-label={`${goal.title}: ${goal.progress}%`}>
                  <div
                    className={`h-full rounded-full transition-all ${
                      goal.progress >= 80 ? 'bg-green-500' : goal.progress >= 50 ? 'bg-blue-500' : 'bg-orange-500'
                    }`}
                    style={{ width: `${goal.progress}%` }}
                  />
                </div>
              </div>
            ))}
            {activeGoals.length > 3 && (
              <Link
                to="/goals"
                className="flex items-center justify-center gap-1 text-xs text-gray-400 hover:text-blue-400 pt-1 transition-colors"
              >
                {t('dashboard.moreGoals', { count: activeGoals.length - 3 })}
                <ChevronRight className="w-3 h-3" aria-hidden="true" />
              </Link>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
