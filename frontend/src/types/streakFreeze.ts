export interface StreakFreezeStatus {
  max_freezes: number;
  used_freezes: number;
  remaining: number;
  used_dates: string[];
  today_used: boolean;
  can_use_today: boolean;
}
