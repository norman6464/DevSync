import client from './client';
import type { LevelInfo, XPBreakdown } from '../types/level';

/** 自分のレベル情報を取得 */
export const getMyLevelInfo = () =>
  client.get<LevelInfo>('/level/me');

/** 指定ユーザーのレベル情報を取得 */
export const getLevelInfo = (userId: number) =>
  client.get<LevelInfo>(`/level/${userId}`);

/** 指定ユーザーのXP内訳を取得 */
export const getXPBreakdown = (userId: number) =>
  client.get<XPBreakdown>(`/level/${userId}/breakdown`);
