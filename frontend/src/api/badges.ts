import client from './client';
import type { BadgeResult } from '../types/badge';

export const getUserBadges = (userId: number) =>
  client.get<{ badges: BadgeResult[] }>(`/badges/${userId}`);

export const notifyBadgeEarned = (badgeId: string) =>
  client.post('/badges/notify', { badge_id: badgeId });
