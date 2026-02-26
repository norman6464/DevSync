import client from './client';
import type { User } from '../types/user';

export const getUsers = (q?: string) =>
  client.get<User[]>('/users', { params: q ? { q } : {} });

export const getUser = (id: number) =>
  client.get<User>(`/users/${id}`);

export const getUserByUsername = (username: string) =>
  client.get<User>(`/users/by-username/${username}`);

export const updateUser = (id: number, data: {
  name?: string;
  bio?: string;
  avatar_url?: string;
  skills_languages?: string;
  skills_frameworks?: string;
  atcoder_username?: string;
  paiza_rank?: string;
  onboarding_completed?: boolean;
}) =>
  client.put<User>(`/users/${id}`, data);

export const followUser = (id: number) =>
  client.post(`/users/${id}/follow`);

export const unfollowUser = (id: number) =>
  client.delete(`/users/${id}/follow`);

export const getFollowers = (id: number) =>
  client.get<User[]>(`/users/${id}/followers`);

export const getFollowing = (id: number) =>
  client.get<User[]>(`/users/${id}/following`);

export const getProfileCompleteness = () =>
  client.get<{ percentage: number; missing_fields: string[] }>('/users/me/profile-completeness');

export interface UserActivity {
  id: number;
  user_id: number;
  activity_type: string;
  target_type: string;
  target_id: number;
  metadata: string;
  created_at: string;
}

export interface UserActivityListResponse {
  activities: UserActivity[];
  total: number;
  limit: number;
  offset: number;
}

export const getUserActivity = (userId: number, params?: { type?: string; limit?: number; offset?: number }) =>
  client.get<UserActivityListResponse>(`/users/${userId}/activity`, { params });
