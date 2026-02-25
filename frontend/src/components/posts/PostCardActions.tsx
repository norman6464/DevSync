import { useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Heart, MessageCircle, Smile, Eye, Link2, Bookmark, FolderPlus } from 'lucide-react';
import toast from 'react-hot-toast';
import { likePost, unlikePost, bookmarkPost, unbookmarkPost } from '../../api/posts';
import { useReactions } from '../../hooks/useReactions';
import { AVAILABLE_REACTION_EMOJIS } from '../../constants/reactions';
import BookmarkCollectionPicker from './BookmarkCollectionPicker';

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
  const { reactions, userReactions, toggleReaction } = useReactions(postId);
  const [showReactionPicker, setShowReactionPicker] = useState(false);
  const [linkCopied, setLinkCopied] = useState(false);
  const [showCollectionPicker, setShowCollectionPicker] = useState(false);

  const handleCopyLink = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(`${window.location.origin}/posts/${postId}`);
      setLinkCopied(true);
      toast.success(t('post.linkCopied'));
      setTimeout(() => setLinkCopied(false), 2000);
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  }, [postId, t]);

  const handleReaction = useCallback(async (emoji: string) => {
    await toggleReaction(emoji);
    setShowReactionPicker(false);
  }, [toggleReaction]);

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
          <Heart className="w-4 h-4" fill={liked ? 'currentColor' : 'none'} aria-hidden="true" />
          {likeCount}
        </button>
        <Link
          to={`/posts/${postId}`}
          aria-label={`${t('post.comment')} ${commentCount}`}
          className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-blue-400 transition-colors"
        >
          <MessageCircle className="w-4 h-4" aria-hidden="true" />
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
            <Smile className="w-4 h-4" aria-hidden="true" />
          </button>
          {showReactionPicker && (
            <div className="absolute bottom-8 left-0 z-10 flex gap-1 p-1.5 bg-gray-800 border border-gray-700 rounded-lg shadow-lg" role="menu" aria-label={t('post.addReaction')}>
              {AVAILABLE_REACTION_EMOJIS.map((emoji) => (
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
            <Eye className="w-4 h-4" aria-hidden="true" />
            {viewCount}
          </span>
        )}
        <button
          onClick={handleCopyLink}
          aria-label={t('post.copyLink')}
          className={`flex items-center gap-1.5 text-sm transition-colors ml-auto ${
            linkCopied ? 'text-green-400' : 'text-gray-500 hover:text-gray-300'
          }`}
        >
          <Link2 className="w-4 h-4" aria-hidden="true" />
        </button>
        <button
          onClick={handleBookmark}
          aria-label={bookmarked ? t('post.unbookmark') : t('post.bookmark')}
          aria-pressed={bookmarked}
          className={`flex items-center gap-1.5 text-sm transition-colors ${
            bookmarked ? 'text-yellow-400' : 'text-gray-500 hover:text-yellow-400'
          }`}
        >
          <Bookmark className="w-4 h-4" fill={bookmarked ? 'currentColor' : 'none'} aria-hidden="true" />
          {bookmarkCount > 0 && <span>{bookmarkCount}</span>}
        </button>
        <div className="relative">
          <button
            onClick={() => setShowCollectionPicker(!showCollectionPicker)}
            aria-label={t('bookmarkCollections.saveToCollection')}
            aria-expanded={showCollectionPicker}
            className="text-gray-500 hover:text-blue-400 transition-colors"
          >
            <FolderPlus className="w-4 h-4" aria-hidden="true" />
          </button>
          {showCollectionPicker && (
            <BookmarkCollectionPicker
              postId={postId}
              onClose={() => setShowCollectionPicker(false)}
            />
          )}
        </div>
      </div>
    </>
  );
}
