import client from './client';

export interface ZennArticle {
  id: number;
  user_id: number;
  zenn_id: number;
  title: string;
  slug: string;
  emoji: string;
  article_type: string;
  liked_count: number;
  comments_count: number;
  published_at: string;
  updated_at: string;
}

export interface ZennStats {
  total_articles: number;
  total_likes: number;
  total_comments: number;
}

export const connectZenn = (username: string) =>
  client.post('/zenn/connect', { username });

export const disconnectZenn = () =>
  client.delete('/zenn/disconnect');

export const syncZenn = () =>
  client.post('/zenn/sync');

export const getZennArticles = (userId: number) =>
  client.get<ZennArticle[]>(`/zenn/articles/${userId}`);

export const getZennStats = (userId: number) =>
  client.get<ZennStats>(`/zenn/stats/${userId}`);
