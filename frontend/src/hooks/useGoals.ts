import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyGoals,
  createGoal,
  updateGoal,
  deleteGoal,
  type LearningGoal,
  type GoalCategory,
  type GoalStatus,
} from '../api/goals';
import { useAsyncData } from './useAsyncData';

export function useGoals() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: goals, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyGoals();
      return data || [];
    },
    { initialData: [] as LearningGoal[] }
  );

  const [localGoals, setLocalGoals] = useState<LearningGoal[] | null>(null);
  const currentGoals = localGoals ?? goals;

  // Sync localGoals when remote data changes
  const setGoals = useCallback((updater: LearningGoal[] | ((prev: LearningGoal[]) => LearningGoal[])) => {
    setLocalGoals(prev => {
      const current = prev ?? goals;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [goals]);

  const handleCreate = useCallback(async (data: {
    title: string;
    description: string;
    category: GoalCategory;
    target_date?: string;
  }) => {
    setSaving(true);
    try {
      const { data: newGoal } = await createGoal(data);
      setGoals(prev => [newGoal, ...prev]);
      toast.success(t('goals.created'));
      return newGoal;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setGoals]);

  const handleUpdate = useCallback(async (goalId: number, data: {
    title?: string;
    description?: string;
    category?: GoalCategory;
    target_date?: string;
    progress?: number;
    status?: GoalStatus;
  }) => {
    try {
      const { data: updated } = await updateGoal(goalId, data);
      setGoals(prev => prev.map(g => g.id === updated.id ? updated : g));
      if (data.progress === 100) {
        toast.success(t('goals.completed'));
      } else {
        toast.success(t('goals.updated'));
      }
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, setGoals]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('goals.confirmDelete'))) return false;
    try {
      await deleteGoal(id);
      setGoals(prev => prev.filter(g => g.id !== id));
      toast.success(t('goals.deleted'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setGoals]);

  const activeGoals = currentGoals.filter(g => g.status === 'active');
  const completedGoals = currentGoals.filter(g => g.status === 'completed');
  const pausedGoals = currentGoals.filter(g => g.status === 'paused');

  return {
    goals: currentGoals,
    loading,
    saving,
    activeGoals,
    completedGoals,
    pausedGoals,
    createGoal: handleCreate,
    updateGoal: handleUpdate,
    deleteGoal: handleDelete,
    refetch,
  };
}
