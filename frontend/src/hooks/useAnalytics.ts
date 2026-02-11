import {
  getHeatmap,
  getCategoryBreakdown,
  getProductivityScore,
  getWeeklyTrends,
  getInsights,
} from '../api/analytics';
import type {
  HeatmapEntry,
  CategoryBreakdown,
  ProductivityScore,
  WeeklyTrend,
  AIInsight,
} from '../types/analytics';
import { useAsyncData } from './useAsyncData';

/** ヒートマップデータを取得するフック */
export function useHeatmap(userId: number | undefined) {
  const { data, loading, error, refetch } = useAsyncData<HeatmapEntry[]>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getHeatmap(userId);
      return res.data;
    },
    { deps: [userId], enabled: !!userId }
  );
  return { heatmap: data, loading, error, refetch };
}

/** カテゴリ別学習時間を取得するフック */
export function useCategoryBreakdown(userId: number | undefined) {
  const { data, loading, error, refetch } = useAsyncData<CategoryBreakdown[]>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getCategoryBreakdown(userId);
      return res.data;
    },
    { deps: [userId], enabled: !!userId }
  );
  return { categories: data, loading, error, refetch };
}

/** 生産性スコアを取得するフック */
export function useProductivityScore(userId: number | undefined) {
  const { data, loading, error, refetch } = useAsyncData<ProductivityScore>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getProductivityScore(userId);
      return res.data;
    },
    { deps: [userId], enabled: !!userId }
  );
  return { score: data, loading, error, refetch };
}

/** 週間トレンドを取得するフック */
export function useWeeklyTrends(userId: number | undefined, weeks = 12) {
  const { data, loading, error, refetch } = useAsyncData<WeeklyTrend[]>(
    async () => {
      if (!userId) throw new Error('userId is required');
      const res = await getWeeklyTrends(userId, weeks);
      return res.data;
    },
    { deps: [userId, weeks], enabled: !!userId }
  );
  return { trends: data, loading, error, refetch };
}

/** AIインサイトを取得するフック */
export function useInsights() {
  const { data, loading, error, refetch } = useAsyncData<AIInsight[]>(
    async () => {
      const res = await getInsights();
      return res.data;
    }
  );
  return { insights: data, loading, error, refetch };
}
