import { useState, useCallback, useEffect } from 'react';
import { getReactions, addReaction, removeReaction } from '../api/posts';
import type { ReactionCount } from '../types/post';

export function useReactions(postId: number) {
  const [reactions, setReactions] = useState<ReactionCount[]>([]);
  const [userReactions, setUserReactions] = useState<string[]>([]);

  useEffect(() => {
    getReactions(postId).then((res) => {
      setReactions(res.data.reactions || []);
      setUserReactions(res.data.user_reactions || []);
    }).catch(() => {});
  }, [postId]);

  const toggleReaction = useCallback(async (emoji: string) => {
    try {
      if (userReactions.includes(emoji)) {
        await removeReaction(postId, emoji);
        setUserReactions((prev) => prev.filter((e) => e !== emoji));
        setReactions((prev) =>
          prev.map((r) => r.emoji === emoji ? { ...r, count: r.count - 1 } : r).filter((r) => r.count > 0)
        );
      } else {
        await addReaction(postId, emoji);
        setUserReactions((prev) => [...prev, emoji]);
        setReactions((prev) => {
          const existing = prev.find((r) => r.emoji === emoji);
          if (existing) {
            return prev.map((r) => r.emoji === emoji ? { ...r, count: r.count + 1 } : r);
          }
          return [...prev, { emoji, count: 1 }];
        });
      }
    } catch (e) {
      console.warn('Failed to toggle reaction:', e);
    }
  }, [postId, userReactions]);

  return { reactions, userReactions, toggleReaction };
}
