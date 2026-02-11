/** スタディサークルのステータス */
export type StudyCircleStatus = 'active' | 'completed' | 'archived';

/** メンバーの役割 */
export type StudyCircleMemberRole = 'owner' | 'member';

/** スタディサークル */
export interface StudyCircle {
  id: number;
  name: string;
  topic: string;
  description: string;
  owner_id: number;
  owner?: {
    id: number;
    name: string;
    avatar_url: string;
  };
  max_members: number;
  status: StudyCircleStatus;
  steps?: StudyCircleStep[];
  members?: StudyCircleMember[];
  created_at: string;
  updated_at: string;
}

/** サークルメンバー */
export interface StudyCircleMember {
  id: number;
  circle_id: number;
  user_id: number;
  user?: {
    id: number;
    name: string;
    avatar_url: string;
  };
  role: StudyCircleMemberRole;
  joined_at: string;
}

/** 共有ロードマップのステップ */
export interface StudyCircleStep {
  id: number;
  circle_id: number;
  title: string;
  description: string;
  order_index: number;
  resource_url: string;
  created_at: string;
  updated_at: string;
}

/** メンバー別ステップ進捗 */
export interface StudyCircleMemberProgress {
  id: number;
  circle_id: number;
  step_id: number;
  user_id: number;
  user?: {
    id: number;
    name: string;
    avatar_url: string;
  };
  is_completed: boolean;
  completed_at: string | null;
}

/** 日次チェックイン */
export interface StudyCircleCheckin {
  id: number;
  circle_id: number;
  user_id: number;
  user?: {
    id: number;
    name: string;
    avatar_url: string;
  };
  date: string;
  content: string;
  created_at: string;
}

/** ストリークランキング */
export interface CircleMemberStreak {
  user_id: number;
  user_name: string;
  avatar_url: string;
  current_streak: number;
  total_checkins: number;
}

/** サークル作成リクエスト */
export interface CreateStudyCircleRequest {
  name: string;
  topic: string;
  description?: string;
  max_members?: number;
  member_ids?: number[];
}

/** サークル更新リクエスト */
export interface UpdateStudyCircleRequest {
  name?: string;
  topic?: string;
  description?: string;
}

/** ステップ作成リクエスト */
export interface CreateStepRequest {
  title: string;
  description?: string;
  resource_url?: string;
  order_index?: number;
}

/** ステップ更新リクエスト */
export interface UpdateStepRequest {
  title?: string;
  description?: string;
}
