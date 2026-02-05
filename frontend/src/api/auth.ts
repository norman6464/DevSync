import client from './client';
import type { AuthResponse, User } from '../types/user';

export const register = (name: string, email: string, password: string) =>
  client.post<AuthResponse>('/auth/register', { name, email, password });

export const login = (email: string, password: string) =>
  client.post<AuthResponse>('/auth/login', { email, password });

export const getMe = () =>
  client.get<User>('/auth/me');

export const getGitHubLoginURL = () =>
  client.get<{ url: string }>('/auth/github');

export const gitHubLoginCallback = (code: string, state: string) =>
  client.get<AuthResponse>('/auth/github/callback', { params: { code, state } });
