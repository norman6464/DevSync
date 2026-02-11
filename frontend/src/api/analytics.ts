import client from './client';
import type {
  HeatmapEntry,
  CategoryBreakdown,
  WeeklyTrend,
  ProductivityScore,
  AIInsight,
} from '../types/analytics';

/** 学習時間ヒートマップを取得 */
export const getHeatmap = (userId: number) =>
  client.get<HeatmapEntry[]>(`/analytics/heatmap/${userId}`);

/** カテゴリ別学習時間を取得 */
export const getCategoryBreakdown = (userId: number) =>
  client.get<CategoryBreakdown[]>(`/analytics/categories/${userId}`);

/** 生産性スコアを取得 */
export const getProductivityScore = (userId: number) =>
  client.get<ProductivityScore>(`/analytics/productivity/${userId}`);

/** 週間トレンドを取得 */
export const getWeeklyTrends = (userId: number, weeks = 12) =>
  client.get<WeeklyTrend[]>(`/analytics/trends/${userId}`, { params: { weeks } });

/** AIインサイトを取得 */
export const getInsights = () =>
  client.get<AIInsight[]>('/analytics/insights');
