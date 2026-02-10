import { getMyLevelInfo, getLevelInfo, getXPBreakdown } from '../api/level';
import type { LevelInfo, XPBreakdown } from '../types/level';
import { useAsyncData } from './useAsyncData';

/** 自分のレベル情報を取得するフック */
export function useMyLevel() {
  const { data, loading, error, refetch } = useAsyncData<LevelInfo>(
    async () => {
      const res = await getMyLevelInfo();
      return res.data;
    }
  );

  return { levelInfo: data, loading, error, refetch };
}

/** 指定ユーザーのレベル情報を取得するフック */
export function useLevel(userId: number | undefined) {
  const { data, loading, error, refetch } = useAsyncData<LevelInfo>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getLevelInfo(userId);
      return res.data;
    },
    { deps: [userId], enabled: !!userId }
  );

  return { levelInfo: data, loading, error, refetch };
}

/** 指定ユーザーのXP内訳を取得するフック */
export function useLevelBreakdown(userId: number | undefined) {
  const { data, loading, error, refetch } = useAsyncData<XPBreakdown>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getXPBreakdown(userId);
      return res.data;
    },
    { deps: [userId], enabled: !!userId }
  );

  return { breakdown: data, loading, error, refetch };
}
