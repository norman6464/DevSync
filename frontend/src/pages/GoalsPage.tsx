import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, Copy, type LucideIcon } from 'lucide-react';
import { type GoalCategory, type GoalStatus, type LearningGoal } from '../api/goals';
import { useGoalForm } from '../hooks';
import { Modal, PageLoader } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import ConfirmDialog from '../components/common/ConfirmDialog';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';

const CATEGORIES: { value: GoalCategory; label: string; icon: string; Icon: LucideIcon }[] = [
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

export default function GoalsPage() {
  const { t } = useTranslation();
  const {
    goals, loading, saving, activeGoals, completedGoals, pausedGoals,
    showForm, setShowForm, editingGoal,
    title, setTitle, description, setDescription,
    category, setCategory, targetDate, setTargetDate,
    resetForm, handleSubmit, handleEdit,
    handleProgressChange, handleStatusChange, handleDeleteGoal,
    handleDuplicateGoal,
    dialogProps,
  } = useGoalForm();

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t('goals.title')}</h1>
        <button
          onClick={() => setShowForm(true)}
          className={`${buttonSecondaryClass} font-medium text-sm`}
        >
          {t('goals.addGoal')}
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <p className="text-2xl font-bold">{goals.length}</p>
          <p className="text-sm text-gray-400">{t('goals.totalGoals')}</p>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <p className="text-2xl font-bold text-blue-400">{activeGoals.length}</p>
          <p className="text-sm text-gray-400">{t('goals.activeGoals')}</p>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <p className="text-2xl font-bold text-green-400">{completedGoals.length}</p>
          <p className="text-sm text-gray-400">{t('goals.completedGoals')}</p>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <p className="text-2xl font-bold text-yellow-400">{pausedGoals.length}</p>
          <p className="text-sm text-gray-400">{t('goals.pausedGoals')}</p>
        </div>
      </div>

      {/* Create/Edit Form Modal */}
      <Modal
        isOpen={showForm}
        onClose={resetForm}
        title={editingGoal ? t('goals.editGoal') : t('goals.addGoal')}
        maxWidth="max-w-md"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="goal-title" className={labelClass}>
              {t('goals.goalTitle')}
            </label>
            <input
              id="goal-title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder={t('goals.titlePlaceholder')}
              className={inputClass}
              required
            />
          </div>
          <div>
            <label htmlFor="goal-description" className={labelClass}>
              {t('goals.description')}
            </label>
            <textarea
              id="goal-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t('goals.descriptionPlaceholder')}
              rows={3}
              className={`${inputClass} resize-none`}
            />
          </div>
          <div>
            <label htmlFor="goal-category" className={labelClass}>
              {t('goals.category')}
            </label>
            <select
              id="goal-category"
              value={category}
              onChange={(e) => setCategory(e.target.value as GoalCategory)}
              className={inputClass}
            >
              {CATEGORIES.map((cat) => (
                <option key={cat.value} value={cat.value}>
                  {cat.icon} {t(cat.label)}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="goal-target-date" className={labelClass}>
              {t('goals.targetDate')}
            </label>
            <input
              id="goal-target-date"
              type="date"
              value={targetDate}
              onChange={(e) => setTargetDate(e.target.value)}
              className={inputClass}
            />
          </div>
          <div className="flex gap-3 justify-end pt-2">
            <button
              type="button"
              onClick={resetForm}
              className={`${buttonSecondaryClass} text-sm font-medium`}
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={saving || !title.trim()}
              className={`${buttonSecondaryClass} disabled:opacity-50 text-sm font-medium`}
            >
              {saving ? t('common.loading') : editingGoal ? t('common.save') : t('goals.create')}
            </button>
          </div>
        </form>
      </Modal>

      {/* Goals List */}
      {goals.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-md">
          <EmptyState
            icon={Target}
            message={t('goals.noGoals')}
            actionLabel={t('goals.createFirst')}
            onAction={() => setShowForm(true)}
          />
        </div>
      ) : (
        <div className="space-y-4">
          {activeGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.activeGoals')}
              </h2>
              <div className="space-y-3">
                {activeGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
                    getCategoryInfo={getCategoryInfo}
                    getStatusColor={getStatusColor}
                    t={t}
                  />
                ))}
              </div>
            </div>
          )}

          {pausedGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.pausedGoals')}
              </h2>
              <div className="space-y-3">
                {pausedGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
                    getCategoryInfo={getCategoryInfo}
                    getStatusColor={getStatusColor}
                    t={t}
                  />
                ))}
              </div>
            </div>
          )}

          {completedGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.completedGoals')}
              </h2>
              <div className="space-y-3">
                {completedGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
                    getCategoryInfo={getCategoryInfo}
                    getStatusColor={getStatusColor}
                    t={t}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      <ConfirmDialog {...dialogProps} />
    </div>
  );
}

interface GoalCardProps {
  goal: LearningGoal;
  onEdit: (goal: LearningGoal) => void;
  onDelete: (id: number) => void;
  onDuplicate: (id: number) => void;
  onProgressChange: (goal: LearningGoal, progress: number) => void;
  onStatusChange: (goal: LearningGoal, status: GoalStatus) => void;
  getCategoryInfo: (cat: GoalCategory) => { value: GoalCategory; label: string; icon: string; Icon: LucideIcon };
  getStatusColor: (status: GoalStatus) => string;
  t: (key: string) => string;
}

const GoalCard = memo(function GoalCard({
  goal,
  onEdit,
  onDelete,
  onDuplicate,
  onProgressChange,
  onStatusChange,
  getCategoryInfo,
  getStatusColor,
  t,
}: GoalCardProps) {
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
              <span className={`px-2 py-0.5 text-xs rounded-full ${getStatusColor(goal.status)}`}>
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
                  {t('goals.targetDateLabel')}: {new Date(goal.target_date).toLocaleDateString()}
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
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
            </svg>
          </button>
          <button
            onClick={() => onDelete(goal.id)}
            className="p-2 text-gray-400 hover:text-red-400 transition-colors"
            aria-label={t('common.delete')}
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
});
