import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyLogs,
  createLog,
  updateLog,
  deleteLog,
  getCalendarData,
  getStreakInfo,
  getWeeklyDuration,
  favoriteLog,
  unfavoriteLog,
} from '../api/learningLogs';
import type { LearningLog, CalendarEntry, LogCategory, StreakInfo } from '../types/learningLog';
import { useAsyncData } from './useAsyncData';
import { useCRUDList } from './useCRUDList';

export function useLearningLogs() {
  const { t } = useTranslation();
  const errMsg = t('errors.somethingWrong');

  const { items: logs, loading, saving, setItems, addItem, updateItem, removeItem, refetch } = useCRUDList<LearningLog>({
    fetcher: async () => {
      const { data } = await getMyLogs();
      return data || [];
    },
  });

  const handleCreate = useCallback(async (data: {
    title: string;
    content: string;
    category?: LogCategory;
    duration?: number;
  }) => {
    return addItem(
      async () => { const { data: l } = await createLog(data); return l; },
      { successMsg: t('learningLogs.created'), errorMsg: errMsg },
    );
  }, [addItem, t, errMsg]);

  const handleUpdate = useCallback(async (logId: number, data: {
    title?: string;
    content?: string;
    category?: LogCategory;
    duration?: number;
  }) => {
    return updateItem(
      async () => { const { data: l } = await updateLog(logId, data); return l; },
      { successMsg: t('learningLogs.updated'), errorMsg: errMsg },
    );
  }, [updateItem, t, errMsg]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('learningLogs.confirmDelete'))) return false;
    return removeItem(id, () => deleteLog(id), { successMsg: t('learningLogs.deleted'), errorMsg: errMsg });
  }, [removeItem, t, errMsg]);

  const handleToggleFavorite = useCallback(async (id: number) => {
    const log = logs.find(l => l.id === id);
    if (!log) return;
    try {
      const { data: updated } = log.is_favorite
        ? await unfavoriteLog(id)
        : await favoriteLog(id);
      setItems(prev => prev.map(l => l.id === updated.id ? updated : l));
    } catch {
      toast.error(errMsg);
    }
  }, [logs, errMsg, setItems]);

  return {
    logs,
    loading,
    saving,
    createLog: handleCreate,
    updateLog: handleUpdate,
    deleteLog: handleDelete,
    toggleFavorite: handleToggleFavorite,
    refetch,
  };
}

export function useStreak(userId: number | undefined) {
  const { data: streakInfo, loading, refetch } = useAsyncData(
    async () => {
      if (!userId) return null;
      const { data } = await getStreakInfo(userId);
      return data || null;
    },
    { initialData: null as StreakInfo | null, deps: [userId] }
  );

  return { streakInfo, loading, refetchStreak: refetch };
}

export function useWeeklyDuration(userId: number | undefined) {
  const { data: weeklyDuration, loading } = useAsyncData(
    async () => {
      if (!userId) return 0;
      const { data } = await getWeeklyDuration(userId);
      return data?.duration ?? 0;
    },
    { initialData: 0, deps: [userId] }
  );

  return { weeklyDuration, loading };
}

export function useLearningLogCalendar(userId: number | undefined) {
  const { data: calendarData, loading, refetch } = useAsyncData(
    async () => {
      if (!userId) return [];
      const { data } = await getCalendarData(userId);
      return data || [];
    },
    { initialData: [] as CalendarEntry[], deps: [userId] }
  );

  return { calendarData, loading, refetchCalendar: refetch };
}
