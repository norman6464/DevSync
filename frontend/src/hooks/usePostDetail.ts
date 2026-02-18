import { useState, useCallback } from 'react';
import { getPost, getComments, createComment, updatePost as apiUpdatePost } from '../api/posts';
import type { Post, Comment } from '../types/post';
import { useAsyncData } from './useAsyncData';
import toast from 'react-hot-toast';

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

  const handleSubmitComment = useCallback(async (content: string, parentId?: number) => {
    if (!content.trim() || !postId) return false;
    setSubmitting(true);
    try {
      await createComment(postId, content, parentId);
      await refetch();
      return true;
    } catch {
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [postId, refetch]);

  const handleUpdatePost = useCallback(async (
    title: string,
    content: string,
    imageUrls?: string
  ) => {
    if (!postId) return false;
    setSubmitting(true);
    try {
      await apiUpdatePost(postId, { title, content, image_urls: imageUrls });
      await refetch();
      toast.success('投稿を更新しました');
      return true;
    } catch {
      toast.error('投稿の更新に失敗しました');
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
    updatePost: handleUpdatePost,
    refetch,
  };
}
