import client from './client';
import type { LearningLog, CalendarEntry, CreateLogRequest, BatchCreateLogRequest, UpdateLogRequest, StreakInfo } from '../types/learningLog';

export const createLog = (data: CreateLogRequest) =>
  client.post<LearningLog>('/learning-logs', data);

export const batchCreateLogs = (data: BatchCreateLogRequest) =>
  client.post<LearningLog[]>('/learning-logs/batch', data);

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

export const getRecentCategories = () =>
  client.get<string[]>('/learning-logs/recent-categories');

export type ExportPeriod = '7' | '30' | '90' | 'all';
export type ExportFormat = 'csv' | 'json';

export const exportLogsCSV = (period: ExportPeriod = '30') =>
  client.get<Blob>(`/learning-logs/export?period=${period}`, { responseType: 'blob' });

export const exportLogsJSON = (period: ExportPeriod = '30') =>
  client.get<Blob>(`/learning-logs/export?format=json&period=${period}`, { responseType: 'blob' });

export const exportLogs = (period: ExportPeriod = '30', format: ExportFormat = 'csv') =>
  client.get<Blob>(`/learning-logs/export?format=${format}&period=${period}`, { responseType: 'blob' });
