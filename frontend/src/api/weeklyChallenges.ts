import client from './client';
import type { WeeklyChallenge } from '../types/weeklyChallenge';

export const getCurrentChallenge = () =>
  client.get<WeeklyChallenge>('/weekly-challenges/current');

export const updateChallengeProgress = (value: number) =>
  client.put<WeeklyChallenge>('/weekly-challenges/progress', { value });
