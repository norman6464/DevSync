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
