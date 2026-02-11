/** 曜日×時間帯ごとの学習時間（ヒートマップ用） */
export interface HeatmapEntry {
  day_of_week: number; // 0=日曜 〜 6=土曜
  hour: number;        // 0〜23
  total_minutes: number;
}

/** カテゴリ別学習時間の内訳 */
export interface CategoryBreakdown {
  category: string;
  total_minutes: number;
  log_count: number;
  percentage: number;
}

/** 週ごとの学習時間推移 */
export interface WeeklyTrend {
  week_start: string;
  total_minutes: number;
  log_count: number;
}

/** 生産性スコア */
export interface ProductivityScore {
  pomodoro_rate: number;
  goal_rate: number;
  streak_consistency: number;
  overall_score: number;
}

/** AIインサイト */
export interface AIInsight {
  type: string;
  message: string;
}
