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

export const getMonthlySummary = (userId: number, months = 12) =>
  client.get<{ month: string; total_minutes: number; log_count: number }[]>(`/learning-logs/monthly-summary/${userId}?months=${months}`);

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

export const getFavoriteLogs = (limit = 20, offset = 0) =>
  client.get<{ logs: LearningLog[]; total: number; limit: number; offset: number }>(`/learning-logs/favorites?limit=${limit}&offset=${offset}`);

export type ExportPeriod = '7' | '30' | '90' | 'all';
export type ExportFormat = 'csv' | 'json';

export const exportLogsCSV = (period: ExportPeriod = '30') =>
  client.get<Blob>(`/learning-logs/export?period=${period}`, { responseType: 'blob' });

export const exportLogsJSON = (period: ExportPeriod = '30') =>
  client.get<Blob>(`/learning-logs/export?format=json&period=${period}`, { responseType: 'blob' });

export const exportLogs = (period: ExportPeriod = '30', format: ExportFormat = 'csv') =>
  client.get<Blob>(`/learning-logs/export?format=${format}&period=${period}`, { responseType: 'blob' });

export interface ImportCSVResponse {
  imported: number;
  logs: LearningLog[];
}

export const importLogsCSV = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  return client.post<ImportCSVResponse>('/learning-logs/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};
