import { useCallback } from 'react';
import { getCurrentChallenge, updateChallengeProgress } from '../api/weeklyChallenges';
import type { WeeklyChallenge } from '../types/weeklyChallenge';
import { useAsyncData } from './useAsyncData';

export function useWeeklyChallenge() {
  const { data: challenge, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getCurrentChallenge();
      return data;
    },
    { initialData: null as WeeklyChallenge | null }
  );

  const updateProgress = useCallback(async (value: number) => {
    await updateChallengeProgress(value);
    await refetch();
  }, [refetch]);

  return {
    challenge,
    loading,
    refetch,
    updateProgress,
  };
}
