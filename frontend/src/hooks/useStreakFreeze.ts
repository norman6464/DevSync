import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { getFreezeStatus, useStreakFreezeAPI } from '../api/streakFreezes';
import type { StreakFreezeStatus } from '../types/streakFreeze';
import { useAsyncData } from './useAsyncData';

export function useStreakFreeze() {
  const { t } = useTranslation();

  const { data: freezeStatus, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getFreezeStatus();
      return data || null;
    },
    { initialData: null as StreakFreezeStatus | null }
  );

  const useFreeze = useCallback(async () => {
    try {
      await useStreakFreezeAPI();
      toast.success(t('streak.freezeUsed'));
      refetch();
      return true;
    } catch {
      toast.error(t('streak.freezeError'));
      return false;
    }
  }, [t, refetch]);

  return {
    freezeStatus,
    loading,
    useFreeze,
    refetchFreezeStatus: refetch,
  };
}
