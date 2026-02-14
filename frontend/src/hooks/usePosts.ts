import { useState, useCallback } from 'react';
import { getTimeline, getPosts, createPost, updatePost, deletePost as apiDeletePost } from '../api/posts';
import type { Post } from '../types/post';
import { useAsyncData } from './useAsyncData';
import toast from 'react-hot-toast';

export function usePosts() {
  const [tab, setTab] = useState<'timeline' | 'all'>('timeline');

  const { data: posts, loading, refetch } = useAsyncData(
    async () => {
      const { data } = tab === 'timeline' ? await getTimeline() : await getPosts();
      return data || [];
    },
    { initialData: [] as Post[], deps: [tab] }
  );

  const handleCreatePost = useCallback(async (
    title: string,
    content: string,
    imageUrls?: string,
    codeSnippets?: { language: string; file_name?: string; code: string }[],
    isDraft?: boolean
  ) => {
    await createPost({ title, content, image_urls: imageUrls, code_snippets: codeSnippets, is_draft: isDraft });
    refetch();
  }, [refetch]);

  const handleUpdatePost = useCallback(async (
    id: number,
    title: string,
    content: string,
    imageUrls?: string
  ) => {
    try {
      await updatePost(id, { title, content, image_urls: imageUrls });
      await refetch();
      toast.success('投稿を更新しました');
      return true;
    } catch {
      toast.error('投稿の更新に失敗しました');
      return false;
    }
  }, [refetch]);

  const handleDeletePost = useCallback(async (post: Post) => {
    if (!confirm(`「${post.title}」を削除してもよろしいですか？`)) return false;

    try {
      await apiDeletePost(post.id);
      await refetch();
      toast.success('投稿を削除しました');
      return true;
    } catch {
      toast.error('投稿の削除に失敗しました');
      return false;
    }
  }, [refetch]);

  return {
    posts,
    loading,
    tab,
    setTab,
    createPost: handleCreatePost,
    updatePost: handleUpdatePost,
    deletePost: handleDeletePost,
    refetch,
  };
}
