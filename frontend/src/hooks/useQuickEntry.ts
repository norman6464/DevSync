import { useState, useCallback } from 'react';
import toast from 'react-hot-toast';
import { useTranslation } from 'react-i18next';
import { createLog, getRecentCategories } from '../api/learningLogs';
import type { LogCategory } from '../types/learningLog';
import { useAsyncData } from './useAsyncData';

export function useQuickEntry() {
  const { t } = useTranslation();
  const [submitting, setSubmitting] = useState(false);

  const { data: recentCategories, refetch: refetchCategories } = useAsyncData(
    async () => {
      const { data } = await getRecentCategories();
      return data || [];
    },
    { initialData: [] as string[] }
  );

  const submit = useCallback(async (
    category: LogCategory,
    duration: number,
    note: string
  ) => {
    setSubmitting(true);
    try {
      await createLog({
        title: `${t(`quickEntry.categories.${category}`)} - ${duration}${t('quickEntry.minutes')}`,
        content: note || t('quickEntry.defaultNote'),
        category,
        duration,
        source: 'manual',
      });
      toast.success(t('quickEntry.success'));
      refetchCategories();
      return true;
    } catch {
      toast.error(t('quickEntry.failed'));
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [t, refetchCategories]);

  return {
    recentCategories,
    submitting,
    submit,
  };
}
