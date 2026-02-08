import client from './client';
import type { LearningLog, CalendarEntry, CreateLogRequest, UpdateLogRequest, StreakInfo } from '../types/learningLog';

export const createLog = (data: CreateLogRequest) =>
  client.post<LearningLog>('/learning-logs', data);

export const updateLog = (id: number, data: UpdateLogRequest) =>
  client.put<LearningLog>(`/learning-logs/${id}`, data);

export const deleteLog = (id: number) =>
  client.delete(`/learning-logs/${id}`);

export const getMyLogs = () =>
  client.get<LearningLog[]>('/learning-logs');

export const getLogById = (id: number) =>
  client.get<LearningLog>(`/learning-logs/${id}`);

export const getUserLogs = (userId: number) =>
  client.get<LearningLog[]>(`/learning-logs/user/${userId}`);

export const getCalendarData = (userId: number) =>
  client.get<CalendarEntry[]>(`/learning-logs/calendar/${userId}`);

export const getStreakInfo = (userId: number) =>
  client.get<StreakInfo>(`/learning-logs/streak/${userId}`);
