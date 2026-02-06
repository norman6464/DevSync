import type { User } from './user';
import type { Post } from './post';

export type NotificationType = 'post' | 'message';

export interface Notification {
  id: number;
  user_id: number;
  type: NotificationType;
  actor_id: number;
  actor: User;
  post_id?: number;
  post?: Post;
  read: boolean;
  created_at: string;
}
