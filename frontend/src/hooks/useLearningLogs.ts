import { useState, useCallback } from 'react';
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
} from '../api/learningLogs';
import type { LearningLog, CalendarEntry, LogCategory, StreakInfo } from '../types/learningLog';
import { useAsyncData } from './useAsyncData';

export function useLearningLogs() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: logs, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyLogs();
      return data || [];
    },
    { initialData: [] as LearningLog[] }
  );

  const [localLogs, setLocalLogs] = useState<LearningLog[] | null>(null);
  const currentLogs = localLogs ?? logs;

  const setLogs = useCallback((updater: LearningLog[] | ((prev: LearningLog[]) => LearningLog[])) => {
    setLocalLogs(prev => {
      const current = prev ?? logs;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [logs]);

  const handleCreate = useCallback(async (data: {
    title: string;
    content: string;
    category?: LogCategory;
    duration?: number;
  }) => {
    setSaving(true);
    try {
      const { data: newLog } = await createLog(data);
      setLogs(prev => [newLog, ...prev]);
      toast.success(t('learningLogs.created'));
      return newLog;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setLogs]);

  const handleUpdate = useCallback(async (logId: number, data: {
    title?: string;
    content?: string;
    category?: LogCategory;
    duration?: number;
  }) => {
    try {
      const { data: updated } = await updateLog(logId, data);
      setLogs(prev => prev.map(l => l.id === updated.id ? updated : l));
      toast.success(t('learningLogs.updated'));
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, setLogs]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('learningLogs.confirmDelete'))) return false;
    try {
      await deleteLog(id);
      setLogs(prev => prev.filter(l => l.id !== id));
      toast.success(t('learningLogs.deleted'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setLogs]);

  return {
    logs: currentLogs,
    loading,
    saving,
    createLog: handleCreate,
    updateLog: handleUpdate,
    deleteLog: handleDelete,
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
