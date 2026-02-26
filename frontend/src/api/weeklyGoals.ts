import client from './client';

export type LogCategory = 'coding' | 'reading' | 'video' | 'exercise' | 'other';

export interface WeeklyGoal {
  id: number;
  user_id: number;
  category: LogCategory;
  target_minutes: number;
  created_at: string;
  updated_at: string;
}

export interface WeeklyGoalProgress {
  category: LogCategory;
  target_minutes: number;
  actual_minutes: number;
  progress_percent: number;
}

export interface SetWeeklyGoalRequest {
  category: LogCategory;
  target_minutes: number;
}

export const setWeeklyGoal = (data: SetWeeklyGoalRequest) =>
  client.put<WeeklyGoal>('/weekly-goals', data);

export const getWeeklyGoals = () =>
  client.get<WeeklyGoal[]>('/weekly-goals');

export const getWeeklyGoalProgress = () =>
  client.get<WeeklyGoalProgress[]>('/weekly-goals/progress');
