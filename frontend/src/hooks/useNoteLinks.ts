import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getNoteLinks,
  getNoteBacklinks,
  createNoteLink,
  deleteNoteLink,
  type NoteLink,
  type CreateNoteLinkRequest,
} from '../api/noteLinks';
import { useAsyncData } from './useAsyncData';

export function useNoteLinks(noteId: number) {
  const { t } = useTranslation();

  const { data: links, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getNoteLinks(noteId);
      return data || [];
    },
    { initialData: [], deps: [noteId] }
  );

  const handleCreateLink = useCallback(async (data: CreateNoteLinkRequest) => {
    try {
      await createNoteLink(noteId, data);
      await refetch();
      toast.success(t('noteLinks.created', { defaultValue: 'リンクを作成しました' }));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [noteId, refetch, t]);

  const handleDeleteLink = useCallback(async (targetNoteId: number) => {
    if (!confirm(t('noteLinks.confirmDelete', { defaultValue: 'このリンクを削除しますか？' }))) return false;
    try {
      await deleteNoteLink(noteId, targetNoteId);
      await refetch();
      toast.success(t('noteLinks.deleted', { defaultValue: 'リンクを削除しました' }));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [noteId, refetch, t]);

  return {
    links,
    loading,
    createLink: handleCreateLink,
    deleteLink: handleDeleteLink,
    refetch,
  };
}

export function useNoteBacklinks(noteId: number) {
  const { data: backlinks, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getNoteBacklinks(noteId);
      return data || [];
    },
    { initialData: [], deps: [noteId] }
  );

  return {
    backlinks,
    loading,
    refetch,
  };
}
