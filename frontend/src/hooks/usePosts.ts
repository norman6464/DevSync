import { useState, useCallback } from 'react';
import { getTimeline, getPosts, createPost } from '../api/posts';
import type { Post } from '../types/post';
import { useAsyncData } from './useAsyncData';

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
    codeSnippets?: { language: string; file_name?: string; code: string }[]
  ) => {
    await createPost({ title, content, image_urls: imageUrls, code_snippets: codeSnippets });
    refetch();
  }, [refetch]);

  return {
    posts,
    loading,
    tab,
    setTab,
    createPost: handleCreatePost,
    refetch,
  };
}
