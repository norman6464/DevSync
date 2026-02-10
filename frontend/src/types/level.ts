/** ユーザーのレベル情報 */
export interface LevelInfo {
  level: number;
  total_xp: number;
  current_level_xp: number;
  next_level_xp: number;
  progress_xp: number;
  progress_percent: number;
}

/** XPの内訳 */
export interface XPBreakdown {
  learning_logs: number;
  posts: number;
  github: number;
  goals: number;
  comments: number;
  likes: number;
  streak_bonus: number;
  total: number;
}
