import client from './client';
import type { RankingEntry } from '../types/ranking';

export const getContributionRanking = (period: 'weekly' | 'monthly' = 'weekly') =>
  client.get<RankingEntry[]>('/rankings/contributions', { params: { period } });

export const getLanguageRanking = (language: string, period: 'weekly' | 'monthly' = 'weekly') =>
  client.get<RankingEntry[]>(`/rankings/languages/${language}`, { params: { period } });

export const getAvailableLanguages = () =>
  client.get<string[]>('/rankings/languages');
