import client from './client';

export type GoalStatus = 'active' | 'completed' | 'paused';
export type GoalCategory = 'language' | 'framework' | 'skill' | 'project' | 'other';

export interface LearningGoal {
  id: number;
  user_id: number;
  title: string;
  description: string;
  category: GoalCategory;
  target_date: string | null;
  progress: number;
  status: GoalStatus;
  created_at: string;
  updated_at: string;
  completed_at: string | null;
}

export interface LearningGoalStats {
  total_goals: number;
  active_goals: number;
  completed_goals: number;
  average_progress: number;
}

export interface CreateGoalRequest {
  title: string;
  description?: string;
  category?: GoalCategory;
  target_date?: string;
}

export interface UpdateGoalRequest {
  title?: string;
  description?: string;
  category?: GoalCategory;
  target_date?: string;
  progress?: number;
  status?: GoalStatus;
}

export const createGoal = (data: CreateGoalRequest) =>
  client.post<LearningGoal>('/goals', data);

export const updateGoal = (id: number, data: UpdateGoalRequest) =>
  client.put<LearningGoal>(`/goals/${id}`, data);

export const deleteGoal = (id: number) =>
  client.delete(`/goals/${id}`);

export const getGoal = (id: number) =>
  client.get<LearningGoal>(`/goals/${id}`);

export const getMyGoals = () =>
  client.get<LearningGoal[]>('/goals');

export const getUserGoals = (userId: number) =>
  client.get<LearningGoal[]>(`/goals/user/${userId}`);

export const getGoalStats = (userId: number) =>
  client.get<LearningGoalStats>(`/goals/stats/${userId}`);
