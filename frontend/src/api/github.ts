import client from './client';
import type { GitHubContribution, GitHubLanguageStat, GitHubRepository } from '../types/github';

export const getGitHubConnectURL = () =>
  client.get<{ url: string }>('/github/connect');

export const gitHubCallback = (code: string, state: string) =>
  client.get('/github/callback', { params: { code, state } });

export const getContributions = (userId: number) =>
  client.get<GitHubContribution[]>(`/github/contributions/${userId}`);

export const getLanguages = (userId: number) =>
  client.get<GitHubLanguageStat[]>(`/github/languages/${userId}`);

export const getRepos = (userId: number) =>
  client.get<GitHubRepository[]>(`/github/repos/${userId}`);

export const syncGitHub = () =>
  client.post('/github/sync');

export const disconnectGitHub = () =>
  client.delete('/github/disconnect');
