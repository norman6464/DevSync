import type { User } from './user';

export interface Post {
  id: number;
  user_id: number;
  user: User;
  title: string;
  content: string;
  image_urls: string;
  like_count: number;
  comment_count: number;
  liked?: boolean;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  user_id: number;
  user: User;
  post_id: number;
  content: string;
  created_at: string;
  updated_at: string;
}
