import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { likePost, unlikePost, bookmarkPost, unbookmarkPost, getReactions, addReaction, removeReaction } from '../../api/posts';
import type { ReactionCount } from '../../types/post';
import { useEffect } from 'react';

interface PostCardActionsProps {
  postId: number;
  initialLiked: boolean;
  initialLikeCount: number;
  initialBookmarked: boolean;
  bookmarkCount: number;
  commentCount: number;
  viewCount: number;
  onUpdate?: () => void;
}

const availableEmojis = ['👍', '🎉', '❤️', '🔥', '👀'];

export default function PostCardActions({
  postId,
  initialLiked,
  initialLikeCount,
  initialBookmarked,
  bookmarkCount,
  commentCount,
  viewCount,
  onUpdate,
}: PostCardActionsProps) {
  const { t } = useTranslation();
  const [liked, setLiked] = useState(initialLiked);
  const [likeCount, setLikeCount] = useState(initialLikeCount);
  const [bookmarked, setBookmarked] = useState(initialBookmarked);
  const [reactions, setReactions] = useState<ReactionCount[]>([]);
  const [userReactions, setUserReactions] = useState<string[]>([]);
  const [showReactionPicker, setShowReactionPicker] = useState(false);

  useEffect(() => {
    getReactions(postId).then((res) => {
      setReactions(res.data.reactions || []);
      setUserReactions(res.data.user_reactions || []);
    }).catch(() => {});
  }, [postId]);

  const handleReaction = async (emoji: string) => {
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
      setShowReactionPicker(false);
    } catch (e) {
      console.warn('Failed to toggle reaction:', e);
    }
  };

  const handleLike = async () => {
    try {
      if (liked) {
        await unlikePost(postId);
        setLiked(false);
        setLikeCount((c) => c - 1);
      } else {
        await likePost(postId);
        setLiked(true);
        setLikeCount((c) => c + 1);
      }
      onUpdate?.();
    } catch (e) {
      console.warn('Failed to toggle like:', e);
    }
  };

  const handleBookmark = async () => {
    try {
      if (bookmarked) {
        await unbookmarkPost(postId);
        setBookmarked(false);
      } else {
        await bookmarkPost(postId);
        setBookmarked(true);
      }
      onUpdate?.();
    } catch (e) {
      console.warn('Failed to toggle bookmark:', e);
    }
  };

  return (
    <>
      {/* Reactions display */}
      {reactions.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mt-3" role="group" aria-label={t('post.addReaction')}>
          {reactions.map((r) => (
            <button
              key={r.emoji}
              onClick={() => handleReaction(r.emoji)}
              aria-label={`${r.emoji} ${r.count}`}
              aria-pressed={userReactions.includes(r.emoji)}
              className={`flex items-center gap-1 px-2 py-0.5 rounded-full text-xs border transition-colors ${
                userReactions.includes(r.emoji)
                  ? 'border-blue-500/50 bg-blue-500/10 text-blue-400'
                  : 'border-gray-700 bg-gray-800/50 text-gray-400 hover:border-gray-600'
              }`}
            >
              <span>{r.emoji}</span>
              <span>{r.count}</span>
            </button>
          ))}
        </div>
      )}

      <div className="flex items-center gap-4 mt-4 pt-3 border-t border-gray-800">
        <button
          onClick={handleLike}
          aria-label={liked ? t('post.unlike') : t('post.like')}
          aria-pressed={liked}
          className={`flex items-center gap-1.5 text-sm transition-colors ${
            liked ? 'text-red-400' : 'text-gray-500 hover:text-red-400'
          }`}
        >
          <svg className="w-4 h-4" fill={liked ? 'currentColor' : 'none'} stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
          </svg>
          {likeCount}
        </button>
        <Link
          to={`/posts/${postId}`}
          aria-label={`${t('post.comment')} ${commentCount}`}
          className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-blue-400 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
          </svg>
          {commentCount}
        </Link>
        <div className="relative">
          <button
            onClick={() => setShowReactionPicker(!showReactionPicker)}
            aria-label={t('post.addReaction')}
            aria-expanded={showReactionPicker}
            aria-haspopup="true"
            className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-gray-300 transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.182 15.182a4.5 4.5 0 0 1-6.364 0M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z" />
            </svg>
          </button>
          {showReactionPicker && (
            <div className="absolute bottom-8 left-0 z-10 flex gap-1 p-1.5 bg-gray-800 border border-gray-700 rounded-lg shadow-lg" role="menu" aria-label={t('post.addReaction')}>
              {availableEmojis.map((emoji) => (
                <button
                  key={emoji}
                  onClick={() => handleReaction(emoji)}
                  role="menuitem"
                  aria-label={emoji}
                  className={`text-lg p-1 rounded hover:bg-gray-700 transition-colors ${
                    userReactions.includes(emoji) ? 'bg-blue-500/20' : ''
                  }`}
                >
                  {emoji}
                </button>
              ))}
            </div>
          )}
        </div>
        {viewCount > 0 && (
          <span className="flex items-center gap-1.5 text-sm text-gray-500" aria-label={t('postViews.views', { count: viewCount })}>
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
            </svg>
            {viewCount}
          </span>
        )}
        <button
          onClick={handleBookmark}
          aria-label={bookmarked ? t('post.unbookmark') : t('post.bookmark')}
          aria-pressed={bookmarked}
          className={`flex items-center gap-1.5 text-sm transition-colors ml-auto ${
            bookmarked ? 'text-yellow-400' : 'text-gray-500 hover:text-yellow-400'
          }`}
        >
          <svg className="w-4 h-4" fill={bookmarked ? 'currentColor' : 'none'} stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
          </svg>
          {bookmarkCount > 0 && <span>{bookmarkCount}</span>}
        </button>
      </div>
    </>
  );
}
