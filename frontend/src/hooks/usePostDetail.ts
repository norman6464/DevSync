import { useState, useCallback } from 'react';
import { getPost, getComments, createComment } from '../api/posts';
import type { Post, Comment } from '../types/post';
import { useAsyncData } from './useAsyncData';

export function usePostDetail(id: string | undefined) {
  const postId = id ? parseInt(id) : 0;
  const [submitting, setSubmitting] = useState(false);

  const { data, loading, refetch } = useAsyncData(
    async () => {
      const [postRes, commentsRes] = await Promise.all([
        getPost(postId),
        getComments(postId),
      ]);
      return {
        post: postRes.data as Post,
        comments: (commentsRes.data || []) as Comment[],
      };
    },
    { deps: [postId], enabled: !!postId }
  );

  const handleSubmitComment = useCallback(async (content: string) => {
    if (!content.trim() || !postId) return false;
    setSubmitting(true);
    try {
      await createComment(postId, content);
      await refetch();
      return true;
    } catch {
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [postId, refetch]);

  return {
    post: data?.post ?? null,
    comments: data?.comments ?? [],
    loading,
    submitting,
    submitComment: handleSubmitComment,
    refetch,
  };
}
