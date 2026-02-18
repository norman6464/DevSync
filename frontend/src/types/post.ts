import type { User } from './user';

export interface Post {
  id: number;
  user_id: number;
  user: User;
  title: string;
  content: string;
  image_urls: string;
  is_draft: boolean;
  like_count: number;
  comment_count: number;
  code_snippets?: CodeSnippet[];
  liked?: boolean;
  bookmarked?: boolean;
  tags?: string[];
  created_at: string;
  updated_at: string;
}

export interface TagCount {
  tag: string;
  count: number;
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

export interface CodeSnippet {
  id: number;
  post_id: number;
  user_id: number;
  language: string;
  file_name: string;
  code: string;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

export interface ReactionCount {
  emoji: string;
  count: number;
}

export interface ReactionResponse {
  reactions: ReactionCount[];
  user_reactions: string[];
}

export interface PostSeries {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface PostSeriesItem {
  id: number;
  series_id: number;
  post_id: number;
  post?: Post;
  order_index: number;
}

export interface PostCollection {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  description: string;
  is_public: boolean;
  created_at: string;
  updated_at: string;
}

export interface PostCollectionItem {
  id: number;
  collection_id: number;
  post_id: number;
  post?: Post;
  note: string;
  order_index: number;
  created_at: string;
}

export interface SnippetComment {
  id: number;
  snippet_id: number;
  user_id: number;
  user: User;
  line_number: number;
  content: string;
  created_at: string;
  updated_at: string;
}
