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

export const getWeeklyDuration = (userId: number) =>
  client.get<{ duration: number }>(`/learning-logs/weekly-duration/${userId}`);

export const getLogsByCategory = (category: string) =>
  client.get<LearningLog[]>(`/learning-logs/category/${category}`);

export const getLogsBySource = (source: string) =>
  client.get<LearningLog[]>(`/learning-logs/source/${source}`);

export const favoriteLog = (id: number) =>
  client.put<LearningLog>(`/learning-logs/${id}/favorite`);

export const unfavoriteLog = (id: number) =>
  client.put<LearningLog>(`/learning-logs/${id}/unfavorite`);

export type ExportPeriod = '7' | '30' | '90' | 'all';

export const exportLogsCSV = (period: ExportPeriod = '30') =>
  client.get<Blob>(`/learning-logs/export?period=${period}`, { responseType: 'blob' });
