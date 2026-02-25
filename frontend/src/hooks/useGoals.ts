import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import {
  getMyGoals,
  createGoal,
  updateGoal,
  deleteGoal,
  duplicateGoal,
  type LearningGoal,
  type GoalCategory,
  type GoalStatus,
} from '../api/goals';
import { useCRUDList } from './useCRUDList';

export function useGoals() {
  const { t } = useTranslation();

  const { items: goals, loading, saving, addItem, updateItem, removeItem, refetch } = useCRUDList<LearningGoal>({
    fetcher: async () => {
      const { data } = await getMyGoals();
      return data || [];
    },
  });

  const errMsg = t('errors.somethingWrong');

  const handleCreate = useCallback(async (data: {
    title: string;
    description: string;
    category: GoalCategory;
    target_date?: string;
  }) => {
    return addItem(
      async () => { const { data: g } = await createGoal(data); return g; },
      { successMsg: t('goals.created'), errorMsg: errMsg },
    );
  }, [addItem, t, errMsg]);

  const handleUpdate = useCallback(async (goalId: number, data: {
    title?: string;
    description?: string;
    category?: GoalCategory;
    target_date?: string;
    progress?: number;
    status?: GoalStatus;
  }) => {
    return updateItem(
      async () => { const { data: g } = await updateGoal(goalId, data); return g; },
      { successMsg: data.progress === 100 ? t('goals.completed') : t('goals.updated'), errorMsg: errMsg },
    );
  }, [updateItem, t, errMsg]);

  const handleDelete = useCallback(async (id: number) => {
    return removeItem(id, () => deleteGoal(id), { successMsg: t('goals.deleted'), errorMsg: errMsg });
  }, [removeItem, t, errMsg]);

  const handleDuplicate = useCallback(async (id: number) => {
    return addItem(
      async () => { const { data: g } = await duplicateGoal(id); return g; },
      { successMsg: t('goals.duplicated'), errorMsg: errMsg, trackSaving: false },
    );
  }, [addItem, t, errMsg]);

  const activeGoals = goals.filter(g => g.status === 'active');
  const completedGoals = goals.filter(g => g.status === 'completed');
  const pausedGoals = goals.filter(g => g.status === 'paused');

  return {
    goals,
    loading,
    saving,
    activeGoals,
    completedGoals,
    pausedGoals,
    createGoal: handleCreate,
    updateGoal: handleUpdate,
    deleteGoal: handleDelete,
    duplicateGoal: handleDuplicate,
    refetch,
  };
}
