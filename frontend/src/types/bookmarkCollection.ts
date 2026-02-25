import type { Post } from './post';

export interface BookmarkCollection {
  id: number;
  user_id: number;
  name: string;
  description: string;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface BookmarkCollectionWithPosts {
  posts: Post[];
  total: number;
}
