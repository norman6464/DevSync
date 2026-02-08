import { useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useToast } from '../contexts/ToastContext';
import { notifyBadgeEarned } from '../api/badges';
import type { BadgeResult } from '../types/badge';

const STORAGE_KEY = 'devsync_earned_badges';

export function useBadgeNotifier(badges: BadgeResult[]) {
  const { t } = useTranslation();
  const { showToast } = useToast();
  const initialLoadRef = useRef(true);

  useEffect(() => {
    if (badges.length === 0) return;

    const earnedBadges = badges.filter((b) => b.earned);
    const earnedIds = earnedBadges.map((b) => b.id);

    const stored = localStorage.getItem(STORAGE_KEY);
    const previousIds: string[] = stored ? JSON.parse(stored) : [];

    // Skip toasts on initial load when localStorage is empty
    if (initialLoadRef.current) {
      initialLoadRef.current = false;
      if (previousIds.length === 0) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(earnedIds));
        return;
      }
    }

    // Detect newly earned badges
    const newBadges = earnedBadges.filter((b) => !previousIds.includes(b.id));

    if (newBadges.length > 0) {
      // Show toast for each new badge
      newBadges.forEach((badge) => {
        showToast(t('badges.badgeEarned', { name: t(badge.name) }), 'success');
        // Notify backend
        notifyBadgeEarned(badge.id).catch(() => {});
      });

      // Update localStorage
      localStorage.setItem(STORAGE_KEY, JSON.stringify(earnedIds));
    } else if (earnedIds.length !== previousIds.length) {
      // Update if badge count changed without new ones (edge case)
      localStorage.setItem(STORAGE_KEY, JSON.stringify(earnedIds));
    }
  }, [badges, showToast, t]);
}
