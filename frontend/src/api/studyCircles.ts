import client from './client';
import type {
  StudyCircle,
  StudyCircleMember,
  StudyCircleStep,
  StudyCircleMemberProgress,
  StudyCircleCheckin,
  CircleMemberStreak,
  CreateStudyCircleRequest,
  UpdateStudyCircleRequest,
  CreateStepRequest,
  UpdateStepRequest,
} from '../types/studyCircle';

// サークルCRUD
export const getMyCircles = () =>
  client.get<StudyCircle[]>('/study-circles');

export const getCircle = (id: number) =>
  client.get<StudyCircle>(`/study-circles/${id}`);

export const createCircle = (data: CreateStudyCircleRequest) =>
  client.post<StudyCircle>('/study-circles', data);

export const updateCircle = (id: number, data: UpdateStudyCircleRequest) =>
  client.put<StudyCircle>(`/study-circles/${id}`, data);

export const deleteCircle = (id: number) =>
  client.delete(`/study-circles/${id}`);

// メンバー管理
export const getMembers = (circleId: number) =>
  client.get<StudyCircleMember[]>(`/study-circles/${circleId}/members`);

export const addMember = (circleId: number, userId: number) =>
  client.post(`/study-circles/${circleId}/members`, { user_id: userId });

export const removeMember = (circleId: number, userId: number) =>
  client.delete(`/study-circles/${circleId}/members/${userId}`);

// ステップCRUD
export const createStep = (circleId: number, data: CreateStepRequest) =>
  client.post<StudyCircleStep>(`/study-circles/${circleId}/steps`, data);

export const updateStep = (circleId: number, stepId: number, data: UpdateStepRequest) =>
  client.put<StudyCircleStep>(`/study-circles/${circleId}/steps/${stepId}`, data);

export const deleteStep = (circleId: number, stepId: number) =>
  client.delete(`/study-circles/${circleId}/steps/${stepId}`);

export const reorderSteps = (circleId: number, orders: { orders: Array<{ step_id: number; order_index: number }> }) =>
  client.put(`/study-circles/${circleId}/steps/reorder`, orders);

// 進捗
export const updateProgress = (circleId: number, stepId: number, isCompleted: boolean) =>
  client.put(`/study-circles/${circleId}/steps/${stepId}/progress`, { is_completed: isCompleted });

export const getProgress = (circleId: number) =>
  client.get<StudyCircleMemberProgress[]>(`/study-circles/${circleId}/progress`);

// チェックイン
export const createCheckin = (circleId: number, content: string) =>
  client.post<StudyCircleCheckin>(`/study-circles/${circleId}/checkins`, { content });

export const getCheckins = (circleId: number) =>
  client.get<StudyCircleCheckin[]>(`/study-circles/${circleId}/checkins`);

// ストリークランキング
export const getStreakRanking = (circleId: number) =>
  client.get<CircleMemberStreak[]>(`/study-circles/${circleId}/streak-ranking`);

// 検索
export const searchCircles = (query: string, limit = 20, offset = 0) =>
  client.get<StudyCircle[]>('/search/circles', { params: { q: query, limit, offset } });

export const getCirclesByStatus = (status: string) =>
  client.get<StudyCircle[]>(`/study-circles/status/${status}`);

export const updateMemberRole = (circleId: number, userId: number, role: 'owner' | 'member') =>
  client.put(`/study-circles/${circleId}/members/${userId}/role`, { role });
