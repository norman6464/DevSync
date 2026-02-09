import { useCallback } from 'react';
import { useAsyncData } from './useAsyncData';
import { getSnippetComments, createSnippetComment, deleteSnippetComment } from '../api/snippets';
import type { SnippetComment } from '../types/post';

export function useSnippetComments(snippetId: number) {
  const { data: comments, loading, refetch } = useAsyncData(
    async () => {
      if (!snippetId) return [] as SnippetComment[];
      const res = await getSnippetComments(snippetId);
      return (res.data || []) as SnippetComment[];
    },
    { initialData: [] as SnippetComment[], deps: [snippetId], enabled: !!snippetId }
  );

  const addComment = useCallback(
    async (lineNumber: number, content: string) => {
      if (!content.trim() || !snippetId) return false;
      try {
        await createSnippetComment(snippetId, { line_number: lineNumber, content });
        await refetch();
        return true;
      } catch {
        return false;
      }
    },
    [snippetId, refetch]
  );

  const removeComment = useCallback(
    async (commentId: number) => {
      try {
        await deleteSnippetComment(snippetId, commentId);
        await refetch();
        return true;
      } catch {
        return false;
      }
    },
    [snippetId, refetch]
  );

  return { comments, loading, addComment, removeComment, refetch };
}
