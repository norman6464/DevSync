import client from './client';

export interface QiitaArticle {
  id: number;
  user_id: number;
  qiita_id: string;
  title: string;
  url: string;
  likes_count: number;
  comments_count: number;
  tags: string;
  published_at: string;
  updated_at: string;
}

export interface QiitaStats {
  total_articles: number;
  total_likes: number;
  total_comments: number;
}

export const connectQiita = (username: string) =>
  client.post('/qiita/connect', { username });

export const disconnectQiita = () =>
  client.delete('/qiita/disconnect');

export const syncQiita = () =>
  client.post('/qiita/sync');

export const getQiitaArticles = (userId: number) =>
  client.get<QiitaArticle[]>(`/qiita/articles/${userId}`);

export const getQiitaStats = (userId: number) =>
  client.get<QiitaStats>(`/qiita/stats/${userId}`);
