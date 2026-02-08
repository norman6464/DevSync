import client from './client';
import type { User } from '../types/user';

export interface AtCoderRatingInfo {
  username: string;
  rating: number;
  color: string;
  rank: string;
}

export const connectAtCoder = (username: string) =>
  client.post<User>('/atcoder/connect', { username });

export const disconnectAtCoder = () =>
  client.delete<User>('/atcoder/disconnect');

export const getAtCoderRating = (username: string) =>
  client.get<AtCoderRatingInfo>(`/atcoder/rating/${username}`);
