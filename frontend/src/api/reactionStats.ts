import client from './client';
import type { ReactionSummary } from '../types/post';

export const getReactionSummary = (userId: number) =>
  client.get<ReactionSummary>(`/users/${userId}/reaction-summary`);
