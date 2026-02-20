import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, Copy, Pencil, Trash2, type LucideIcon } from 'lucide-react';
import { type GoalCategory, type GoalStatus, type LearningGoal } from '../../api/goals';
import { formatDate } from '../../utils/timeFormat';
import { badgeBaseClass } from '../../constants/styles';

export const CATEGORIES: { value: GoalCategory; label: string; icon: string; Icon: LucideIcon }[] = [
  { value: 'language', label: 'goals.categoryLanguage', icon: '💻', Icon: Monitor },
  { value: 'framework', label: 'goals.categoryFramework', icon: '🚀', Icon: Rocket },
  { value: 'skill', label: 'goals.categorySkill', icon: '🎯', Icon: Target },
  { value: 'project', label: 'goals.categoryProject', icon: '📁', Icon: FolderOpen },
  { value: 'other', label: 'goals.categoryOther', icon: '📝', Icon: FileText },
];

const getCategoryInfo = (cat: GoalCategory) => {
  return CATEGORIES.find((c) => c.value === cat) || CATEGORIES[4];
};

const getStatusColor = (status: GoalStatus) => {
  switch (status) {
    case 'active':
      return 'text-blue-400 bg-blue-400/10';
    case 'completed':
      return 'text-green-400 bg-green-400/10';
    case 'paused':
      return 'text-yellow-400 bg-yellow-400/10';
  }
};

const getDeadlineInfo = (goal: LearningGoal): { status: 'overdue' | 'approaching' | ''; daysLeft: number } => {
  if (!goal.target_date || goal.status !== 'active') return { status: '', daysLeft: -1 };
  const now = new Date();
  const target = new Date(goal.target_date);
  const diffMs = target.getTime() - now.getTime();
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
  if (days < 0) return { status: 'overdue', daysLeft: days };
  if (days <= 3) return { status: 'approaching', daysLeft: days };
  return { status: '', daysLeft: days };
};

interface GoalCardProps {
  goal: LearningGoal;
  onEdit: (goal: LearningGoal) => void;
  onDelete: (id: number) => void;
  onDuplicate: (id: number) => void;
  onProgressChange: (goal: LearningGoal, progress: number) => void;
  onStatusChange: (goal: LearningGoal, status: GoalStatus) => void;
}

const GoalCard = memo(function GoalCard({
  goal,
  onEdit,
  onDelete,
  onDuplicate,
  onProgressChange,
  onStatusChange,
}: GoalCardProps) {
  const { t } = useTranslation();
  const categoryInfo = getCategoryInfo(goal.category);
  const CategoryIcon = categoryInfo.Icon;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3 min-w-0 flex-1">
          <CategoryIcon className="w-6 h-6 text-purple-400 flex-shrink-0 mt-0.5" />
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="font-medium">{goal.title}</h3>
              <span className={`${badgeBaseClass} ${getStatusColor(goal.status)}`}>
                {t(`goals.status${goal.status.charAt(0).toUpperCase() + goal.status.slice(1)}`)}
              </span>
            </div>
            {goal.description && (
              <p className="text-sm text-gray-400 mt-1">{goal.description}</p>
            )}
            <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
              <span>{t(categoryInfo.label)}</span>
              {goal.target_date && (
                <span>
                  {t('goals.targetDateLabel')}: {formatDate(goal.target_date)}
                </span>
              )}
              {(() => {
                const deadline = getDeadlineInfo(goal);
                if (deadline.status === 'overdue') {
                  return (
                    <span className="px-2 py-0.5 rounded-full text-red-400 bg-red-400/10 font-medium">
                      {t('goals.deadlineOverdue')}
                    </span>
                  );
                }
                if (deadline.status === 'approaching') {
                  return (
                    <span className="px-2 py-0.5 rounded-full text-orange-400 bg-orange-400/10 font-medium">
                      {t('goals.deadlineApproaching', { days: deadline.daysLeft })}
                    </span>
                  );
                }
                return null;
              })()}
            </div>

            {/* Progress Bar */}
            <div className="mt-3">
              <div className="flex items-center justify-between text-xs mb-1">
                <span className="text-gray-400">{t('goals.progress')}</span>
                <span className="text-gray-300">{goal.progress}%</span>
              </div>
              <div className="h-2 bg-gray-700 rounded-full overflow-hidden">
                <div
                  className={`h-full transition-all ${
                    goal.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'
                  }`}
                  style={{ width: `${goal.progress}%` }}
                />
              </div>
              {goal.status === 'active' && (
                <input
                  type="range"
                  min="0"
                  max="100"
                  value={goal.progress}
                  onChange={(e) => onProgressChange(goal, parseInt(e.target.value))}
                  className="w-full mt-2 accent-blue-500"
                />
              )}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2">
          {goal.status === 'active' && (
            <button
              onClick={() => onStatusChange(goal, 'paused')}
              className="p-2 text-gray-400 hover:text-yellow-400 transition-colors"
              aria-label={t('goals.pause')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 5.25v13.5m-7.5-13.5v13.5" />
              </svg>
            </button>
          )}
          {goal.status === 'paused' && (
            <button
              onClick={() => onStatusChange(goal, 'active')}
              className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
              aria-label={t('goals.resume')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z" />
              </svg>
            </button>
          )}
          <button
            onClick={() => onDuplicate(goal.id)}
            className="p-2 text-gray-400 hover:text-purple-400 transition-colors"
            aria-label={t('goals.duplicate')}
          >
            <Copy className="w-4 h-4" />
          </button>
          <button
            onClick={() => onEdit(goal)}
            className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
            aria-label={t('common.edit')}
          >
            <Pencil className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(goal.id)}
            className="p-2 text-gray-400 hover:text-red-400 transition-colors"
            aria-label={t('common.delete')}
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
});

export default GoalCard;
