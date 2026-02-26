export type LogCategory = 'coding' | 'reading' | 'course' | 'meetup' | 'other';

export type LogSource = 'manual' | 'pomodoro';

export interface LearningLog {
  id: number;
  user_id: number;
  title: string;
  content: string;
  category: LogCategory;
  duration: number;
  source: LogSource;
  is_favorite: boolean;
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
  source?: LogSource;
}

export interface BatchCreateLogRequest {
  logs: CreateLogRequest[];
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
