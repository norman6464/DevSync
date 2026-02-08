export type LogCategory = 'coding' | 'reading' | 'course' | 'meetup' | 'other';

export interface LearningLog {
  id: number;
  user_id: number;
  title: string;
  content: string;
  category: LogCategory;
  duration: number;
  created_at: string;
  updated_at: string;
}

export interface CalendarEntry {
  date: string;
  count: number;
}

export interface CreateLogRequest {
  title: string;
  content: string;
  category?: LogCategory;
  duration?: number;
}

export interface UpdateLogRequest {
  title?: string;
  content?: string;
  category?: LogCategory;
  duration?: number;
}

export interface StreakInfo {
  current_streak: number;
  longest_streak: number;
  total_days: number;
  last_log_date: string;
}
