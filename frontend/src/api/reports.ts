import client from './client';

export interface DailyActivity {
  date: string;
  contributions: number;
  posts: number;
  comments: number;
}

export interface LanguageActivity {
  language: string;
  bytes: number;
  repos: number;
}

export interface ActivityReport {
  period: 'weekly' | 'monthly';
  start_date: string;
  end_date: string;
  user_id: number;
  total_contributions: number;
  posts_created: number;
  comments_created: number;
  likes_received: number;
  goals_completed: number;
  goals_progress: number;
  new_followers: number;
  messages_exchanged: number;
  daily_contributions: DailyActivity[];
  top_languages: LanguageActivity[];
}

export interface ReportComparison {
  contributions_diff: number;
  posts_diff: number;
  followers_diff: number;
  goals_diff: number;
  trend_percentage: number;
}

export const getMyWeeklyReport = () =>
  client.get<ActivityReport>('/reports/weekly');

export const getMyMonthlyReport = () =>
  client.get<ActivityReport>('/reports/monthly');

export const getWeeklyReport = (userId: number) =>
  client.get<ActivityReport>(`/reports/weekly/${userId}`);

export const getMonthlyReport = (userId: number) =>
  client.get<ActivityReport>(`/reports/monthly/${userId}`);

export const getComparison = (period: 'weekly' | 'monthly') =>
  client.get<ReportComparison>(`/reports/comparison?period=${period}`);
