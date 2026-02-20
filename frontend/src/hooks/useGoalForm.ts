import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { type GoalCategory, type GoalStatus, type LearningGoal } from '../api/goals';
import { useGoals, useConfirm } from './index';

export function useGoalForm() {
  const { t } = useTranslation();
  const goalsData = useGoals();
  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingGoal, setEditingGoal] = useState<LearningGoal | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState<GoalCategory>('other');
  const [targetDate, setTargetDate] = useState('');
  const [filterStatus, setFilterStatus] = useState<GoalStatus | 'all'>('all');
  const [filterCategory, setFilterCategory] = useState<GoalCategory | 'all'>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const resetForm = () => {
    setTitle('');
    setDescription('');
    setCategory('other');
    setTargetDate('');
    setEditingGoal(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    if (editingGoal) {
      const result = await goalsData.updateGoal(editingGoal.id, {
        title,
        description,
        category,
        target_date: targetDate || undefined,
      });
      if (result) resetForm();
    } else {
      const result = await goalsData.createGoal({
        title,
        description,
        category,
        target_date: targetDate || undefined,
      });
      if (result) resetForm();
    }
  };

  const handleEdit = (goal: LearningGoal) => {
    setEditingGoal(goal);
    setTitle(goal.title);
    setDescription(goal.description);
    setCategory(goal.category);
    setTargetDate(goal.target_date ? goal.target_date.split('T')[0] : '');
    setShowForm(true);
  };

  const handleProgressChange = async (goal: LearningGoal, newProgress: number) => {
    await goalsData.updateGoal(goal.id, { progress: newProgress });
  };

  const handleStatusChange = async (goal: LearningGoal, newStatus: GoalStatus) => {
    await goalsData.updateGoal(goal.id, { status: newStatus });
  };

  const handleDeleteGoal = async (id: number) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('goals.confirmDelete'), variant: 'danger' });
    if (ok) goalsData.deleteGoal(id);
  };

  const handleDuplicateGoal = async (id: number) => {
    await goalsData.duplicateGoal(id);
  };

  const filteredGoals = goalsData.goals.filter(g => {
    if (filterStatus !== 'all' && g.status !== filterStatus) return false;
    if (filterCategory !== 'all' && g.category !== filterCategory) return false;
    const q = searchQuery.toLowerCase().trim();
    if (q && !g.title.toLowerCase().includes(q) && !g.description.toLowerCase().includes(q)) return false;
    return true;
  });

  const filteredActiveGoals = filteredGoals.filter(g => g.status === 'active');
  const filteredPausedGoals = filteredGoals.filter(g => g.status === 'paused');
  const filteredCompletedGoals = filteredGoals.filter(g => g.status === 'completed');

  return {
    ...goalsData,
    filteredGoals,
    filteredActiveGoals,
    filteredPausedGoals,
    filteredCompletedGoals,
    filterStatus,
    setFilterStatus,
    filterCategory,
    setFilterCategory,
    searchQuery,
    setSearchQuery,
    showForm,
    setShowForm,
    editingGoal,
    title,
    setTitle,
    description,
    setDescription,
    category,
    setCategory,
    targetDate,
    setTargetDate,
    resetForm,
    handleSubmit,
    handleEdit,
    handleProgressChange,
    handleStatusChange,
    handleDeleteGoal,
    handleDuplicateGoal,
    dialogProps,
  };
}
