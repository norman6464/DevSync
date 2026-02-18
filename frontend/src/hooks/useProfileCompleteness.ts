import { getProfileCompleteness } from '../api/users';
import { useAsyncData } from './useAsyncData';

export function useProfileCompleteness() {
  const { data, loading, refetch } = useAsyncData(
    async () => {
      const res = await getProfileCompleteness();
      return res.data;
    },
    { deps: [] }
  );

  return {
    percentage: data?.percentage ?? 0,
    missingFields: data?.missing_fields ?? [],
    loading,
    refetch,
  };
}
