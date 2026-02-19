import { useState, useCallback } from 'react';
import { likeComment, unlikeComment } from '../api/posts';
import toast from 'react-hot-toast';

export function useCommentLike(
  commentId: number,
  initialLiked: boolean,
  initialCount: number,
) {
  const [liked, setLiked] = useState(initialLiked);
  const [likeCount, setLikeCount] = useState(initialCount);
  const [loading, setLoading] = useState(false);

  const toggleLike = useCallback(async () => {
    if (loading) return;
    setLoading(true);
    const prevLiked = liked;
    const prevCount = likeCount;

    // 楽観的UI更新
    setLiked(!prevLiked);
    setLikeCount(prevLiked ? prevCount - 1 : prevCount + 1);

    try {
      if (prevLiked) {
        await unlikeComment(commentId);
      } else {
        await likeComment(commentId);
      }
    } catch {
      // 失敗時はロールバック
      setLiked(prevLiked);
      setLikeCount(prevCount);
      toast.error('操作に失敗しました');
    } finally {
      setLoading(false);
    }
  }, [commentId, liked, likeCount, loading]);

  return { liked, likeCount, toggleLike, loading };
}
