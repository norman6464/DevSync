export type ChallengeType = 'duration_total' | 'streak_days' | 'category_count' | 'log_count';

export interface WeeklyChallenge {
  id: number;
  user_id: number;
  year: number;
  week: number;
  challenge_type: ChallengeType;
  description: string;
  target_value: number;
  current_value: number;
  is_completed: boolean;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}
