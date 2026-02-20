import { useTranslation } from 'react-i18next';
import { Target, Filter, Search } from 'lucide-react';
import { type GoalCategory } from '../api/goals';
import { useGoalForm } from '../hooks';
import { Modal, PageLoader } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import ConfirmDialog from '../components/common/ConfirmDialog';
import GoalCard, { CATEGORIES } from '../components/goals/GoalCard';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';

export default function GoalsPage() {
  const { t } = useTranslation();
  const {
    goals, loading, saving, activeGoals, completedGoals, pausedGoals,
    filteredGoals, filteredActiveGoals, filteredPausedGoals, filteredCompletedGoals,
    filterStatus, setFilterStatus, filterCategory, setFilterCategory,
    searchQuery, setSearchQuery,
    showForm, setShowForm, editingGoal,
    title, setTitle, description, setDescription,
    category, setCategory, targetDate, setTargetDate,
    resetForm, handleSubmit, handleEdit,
    handleProgressChange, handleStatusChange, handleDeleteGoal,
    handleDuplicateGoal,
    dialogProps,
  } = useGoalForm();

  const isFiltered = filterStatus !== 'all' || filterCategory !== 'all' || searchQuery.trim() !== '';

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

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder={t('goals.searchPlaceholder')}
          className={`${inputClass} pl-9`}
        />
      </div>

      {/* Filters */}
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4 space-y-3">
        <div className="flex items-center gap-2 text-sm text-gray-400">
          <Filter className="w-4 h-4" />
          <span>{t('goals.filter')}</span>
        </div>
        <div className="flex flex-wrap gap-2">
          <span className="text-xs text-gray-500 self-center mr-1">{t('goals.status')}:</span>
          {(['all', 'active', 'paused', 'completed'] as const).map((s) => (
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
              maxLength={200}
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
              maxLength={2000}
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
      ) : isFiltered && filteredGoals.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-8 text-center text-gray-400">
          {t('goals.noFilterResults')}
        </div>
      ) : (
        <div className="space-y-4">
          {filteredActiveGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.activeGoals')}
              </h2>
              <div className="space-y-3">
                {filteredActiveGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
                  />
                ))}
              </div>
            </div>
          )}

          {filteredPausedGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.pausedGoals')}
              </h2>
              <div className="space-y-3">
                {filteredPausedGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
                  />
                ))}
              </div>
            </div>
          )}

          {filteredCompletedGoals.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">
                {t('goals.completedGoals')}
              </h2>
              <div className="space-y-3">
                {filteredCompletedGoals.map((goal) => (
                  <GoalCard
                    key={goal.id}
                    goal={goal}
                    onEdit={handleEdit}
                    onDelete={handleDeleteGoal}
                    onDuplicate={handleDuplicateGoal}
                    onProgressChange={handleProgressChange}
                    onStatusChange={handleStatusChange}
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
