import client from './client';
import type { RecommendedUser, TrendingPost, TrendingResource } from '../types/recommendation';

/** おすすめユーザーを取得 */
export const getRecommendedUsers = () =>
  client.get<RecommendedUser[]>('/recommendations/users');

/** トレンド投稿を取得 */
export const getTrendingPosts = () =>
  client.get<TrendingPost[]>('/recommendations/posts');

/** トレンド学習リソースを取得 */
export const getTrendingResources = () =>
  client.get<TrendingResource[]>('/recommendations/resources');
