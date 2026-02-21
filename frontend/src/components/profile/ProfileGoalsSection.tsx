import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, CheckCircle, ChevronRight, type LucideIcon } from 'lucide-react';

interface Goal {
  id: number;
  title: string;
  description: string;
  category: string;
  status: string;
  progress: number;
}

interface GoalStats {
  active_goals: number;
  completed_goals: number;
}

const categoryIcons: Record<string, LucideIcon> = { language: Monitor, framework: Rocket, skill: Target, project: FolderOpen, other: FileText };
const statusColors: Record<string, string> = { active: 'text-green-400 bg-green-400/10', completed: 'text-blue-400 bg-blue-400/10', paused: 'text-yellow-400 bg-yellow-400/10' };

interface ProfileGoalsSectionProps {
  goals: Goal[];
  goalStats: GoalStats | null;
  isOwnProfile: boolean;
}

export default function ProfileGoalsSection({ goals, goalStats, isOwnProfile }: ProfileGoalsSectionProps) {
  const { t } = useTranslation();

  if (goals.length === 0) return null;

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide flex items-center gap-2">
          <CheckCircle className="w-5 h-5 text-purple-400" aria-hidden="true" />
          {t('goals.title')}
          {goalStats && <span className="text-xs text-gray-500 font-normal ml-2">{goalStats.active_goals} {t('goals.active')} · {goalStats.completed_goals} {t('goals.completed')}</span>}
        </h2>
        {isOwnProfile && (
          <Link to="/goals" className="text-xs text-gray-500 hover:text-purple-400 transition-colors flex items-center gap-1">{t('goals.manageGoals')}<ChevronRight className="w-3 h-3" aria-hidden="true" /></Link>
        )}
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {goals.slice(0, 4).map((goal) => {
          const CategoryIcon = categoryIcons[goal.category] || FileText;
          return (
            <div key={goal.id} className="bg-gray-900 border border-gray-800 rounded-md p-4">
              <div className="flex items-start gap-3">
                <CategoryIcon className="w-6 h-6 text-purple-400 flex-shrink-0 mt-0.5" />
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <div className="font-medium text-sm text-white truncate flex-1">{goal.title}</div>
                    <span className={`px-2 py-0.5 text-xs rounded ${statusColors[goal.status]}`}>{t(`goals.${goal.status}`)}</span>
                  </div>
                  {goal.description && <p className="text-xs text-gray-500 mt-1 line-clamp-2">{goal.description}</p>}
                  <div className="mt-2">
                    <div className="flex items-center justify-between text-xs text-gray-400 mb-1"><span>{t('goals.progress')}</span><span>{goal.progress}%</span></div>
                    <div className="h-1.5 bg-gray-800 rounded-full overflow-hidden" role="progressbar" aria-valuenow={goal.progress} aria-valuemin={0} aria-valuemax={100} aria-label={`${goal.title}: ${goal.progress}%`}><div className="h-full bg-purple-500 rounded-full transition-all" style={{ width: `${goal.progress}%` }} /></div>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
