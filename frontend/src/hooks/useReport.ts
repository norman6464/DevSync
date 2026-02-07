import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyWeeklyReport,
  getMyMonthlyReport,
  getComparison,
  type ActivityReport,
  type ReportComparison,
} from '../api/reports';
import { useAsyncData } from './useAsyncData';

type Period = 'weekly' | 'monthly';

export function useReport() {
  const { t } = useTranslation();
  const [period, setPeriod] = useState<Period>('weekly');

  const { data, loading } = useAsyncData(
    async () => {
      const [reportRes, comparisonRes] = await Promise.all([
        period === 'weekly' ? getMyWeeklyReport() : getMyMonthlyReport(),
        getComparison(period),
      ]);
      return {
        report: reportRes.data as ActivityReport,
        comparison: comparisonRes.data as ReportComparison,
      };
    },
    { deps: [period] }
  );

  if (data === null && !loading) {
    toast.error(t('errors.somethingWrong'));
  }

  return {
    report: data?.report ?? null,
    comparison: data?.comparison ?? null,
    loading,
    period,
    setPeriod,
  };
}
