import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Target, Search, AlertTriangle, Sparkles } from 'lucide-react';
import { useGoalForm } from '../hooks';
import { PageLoader } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import ConfirmDialog from '../components/common/ConfirmDialog';
import GoalCard from '../components/goals/GoalCard';
import GoalFilters from '../components/goals/GoalFilters';
import GoalStatsPanel from '../components/goals/GoalStatsPanel';
import GoalFormModal from '../components/goals/GoalFormModal';
import GoalTemplatesModal from '../components/goals/GoalTemplatesModal';
import { inputClass, buttonSecondaryClass } from '../constants/styles';
import type { GoalTemplate } from '../constants/goalTemplates';

export default function GoalsPage() {
  const { t } = useTranslation();
  const [showTemplates, setShowTemplates] = useState(false);
  const {
    goals, loading, saving, activeGoals, completedGoals, pausedGoals,
    filteredGoals, filteredActiveGoals, filteredPausedGoals, filteredCompletedGoals,
    filterStatus, setFilterStatus, filterCategory, setFilterCategory,
    searchQuery, setSearchQuery, sortBy, setSortBy,
    showForm, setShowForm, editingGoal,
    title, setTitle, description, setDescription,
    category, setCategory, targetDate, setTargetDate,
    resetForm, handleSubmit, handleEdit,
    handleProgressChange, handleStatusChange, handleDeleteGoal,
    handleDuplicateGoal,
    dialogProps,
  } = useGoalForm();

  const overdueGoals = activeGoals.filter((g) => {
    if (!g.target_date) return false;
    return new Date(g.target_date).getTime() < Date.now();
  });

  const isFiltered = filterStatus !== 'all' || filterCategory !== 'all' || searchQuery.trim() !== '';

  const handleTemplateSelect = (template: GoalTemplate) => {
    setTitle(template.title);
    setDescription(template.description);
    setCategory(template.category);
    // 推定日数から目標日を計算
    const targetDateValue = new Date();
    targetDateValue.setDate(targetDateValue.getDate() + template.estimatedDays);
    setTargetDate(targetDateValue.toISOString().split('T')[0]);
    setShowTemplates(false);
    setShowForm(true);
  };

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t('goals.title')}</h1>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowTemplates(true)}
            className={`${buttonSecondaryClass} font-medium text-sm flex items-center gap-1.5`}
          >
            <Sparkles className="w-4 h-4" />
            テンプレートから作成
          </button>
          <button
            onClick={() => setShowForm(true)}
            className={`${buttonSecondaryClass} font-medium text-sm`}
          >
            {t('goals.addGoal')}
          </button>
        </div>
      </div>

      {/* Stats */}
      <GoalStatsPanel
        total={goals.length}
        active={activeGoals.length}
        completed={completedGoals.length}
        paused={pausedGoals.length}
        overdue={overdueGoals.length}
      />

      {/* Overdue Warning */}
      {overdueGoals.length > 0 && (
        <div className="flex items-center gap-3 px-4 py-3 bg-red-500/10 border border-red-500/30 rounded-lg">
          <AlertTriangle className="w-5 h-5 text-red-400 shrink-0" />
          <span className="text-sm text-red-400">
            {t('goals.overdueWarning', { count: overdueGoals.length })}
          </span>
        </div>
      )}

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
      <GoalFilters
        filterStatus={filterStatus}
        setFilterStatus={setFilterStatus}
        filterCategory={filterCategory}
        setFilterCategory={setFilterCategory}
        sortBy={sortBy}
        setSortBy={setSortBy}
      />

      {/* Templates Modal */}
      <GoalTemplatesModal
        isOpen={showTemplates}
        onSelect={handleTemplateSelect}
        onClose={() => setShowTemplates(false)}
      />

      {/* Create/Edit Form Modal */}
      <GoalFormModal
        isOpen={showForm}
        isEditing={!!editingGoal}
        saving={saving}
        title={title}
        setTitle={setTitle}
        description={description}
        setDescription={setDescription}
        category={category}
        setCategory={setCategory}
        targetDate={targetDate}
        setTargetDate={setTargetDate}
        onSubmit={handleSubmit}
        onCancel={resetForm}
      />

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
