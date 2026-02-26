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
  target_hours: number;
  status: GoalStatus;
  is_public: boolean;
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

export interface GoalDeadlineAlert {
  goal: LearningGoal;
  status: 'overdue' | 'approaching';
  days_left: number;
}

export interface GoalProgress {
  goal_id: number;
  target_hours: number;
  actual_minutes: number;
  percentage: number;
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

export const getDeadlineAlerts = () =>
  client.get<GoalDeadlineAlert[]>('/goals/deadline-alerts');

export const duplicateGoal = (id: number) =>
  client.post<LearningGoal>(`/goals/${id}/duplicate`);

export const getLinkedLogs = (goalId: number, limit = 20, offset = 0) =>
  client.get(`/goals/${goalId}/linked-logs`, { params: { limit, offset } });

export const toggleGoalShare = (id: number) =>
  client.put<LearningGoal>(`/goals/${id}/share`);

export const getPublicGoals = (limit = 20, offset = 0) =>
  client.get<LearningGoal[]>('/goals/public', { params: { limit, offset } });

export const getPublicGoalsByUser = (userId: number, limit = 20, offset = 0) =>
  client.get<LearningGoal[]>(`/goals/public/user/${userId}`, { params: { limit, offset } });

export const getGoalProgress = (goalId: number) =>
  client.get<GoalProgress>(`/goals/${goalId}/progress`);
