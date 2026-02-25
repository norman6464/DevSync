import client from './client';
import type { StreakFreezeStatus } from '../types/streakFreeze';

export const getFreezeStatus = () =>
  client.get<StreakFreezeStatus>('/streak-freezes/status');

export const useStreakFreezeAPI = () =>
  client.post<{ message: string }>('/streak-freezes');
