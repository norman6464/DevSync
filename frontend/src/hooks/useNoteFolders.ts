import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyNoteFolders,
  getRootNoteFolders,
  getNoteFolderChildren,
  createNoteFolder,
  updateNoteFolder,
  deleteNoteFolder,
} from '../api/noteFolders';
import type { NoteFolder } from '../types/noteFolder';
import { useAsyncData } from './useAsyncData';

export function useNoteFolders() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: folders, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyNoteFolders();
      return data || [];
    },
    { initialData: [] as NoteFolder[] }
  );

  const [localFolders, setLocalFolders] = useState<NoteFolder[] | null>(null);
  const currentFolders = localFolders ?? folders;

  const setFolders = useCallback((updater: NoteFolder[] | ((prev: NoteFolder[]) => NoteFolder[])) => {
    setLocalFolders(prev => {
      const current = prev ?? folders;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [folders]);

  const handleCreate = useCallback(async (data: {
    name: string;
    parent_id?: number;
  }) => {
    setSaving(true);
    try {
      const { data: newFolder } = await createNoteFolder(data);
      setFolders(prev => [newFolder, ...prev]);
      toast.success(t('notes.folderCreated', { defaultValue: 'フォルダを作成しました' }));
      return newFolder;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setFolders]);

  const handleUpdate = useCallback(async (folderId: number, data: {
    name?: string;
    parent_id?: number;
  }) => {
    try {
      const { data: updated } = await updateNoteFolder(folderId, data);
      setFolders(prev => prev.map(f => f.id === updated.id ? updated : f));
      toast.success(t('notes.folderUpdated', { defaultValue: 'フォルダを更新しました' }));
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, setFolders]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('notes.confirmDeleteFolder', { defaultValue: 'このフォルダを削除しますか？' }))) return false;
    try {
      await deleteNoteFolder(id);
      setFolders(prev => prev.filter(f => f.id !== id));
      toast.success(t('notes.folderDeleted', { defaultValue: 'フォルダを削除しました' }));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setFolders]);

  return {
    folders: currentFolders,
    loading,
    saving,
    createFolder: handleCreate,
    updateFolder: handleUpdate,
    deleteFolder: handleDelete,
    refetch,
  };
}

export function useRootFolders() {
  const { data: rootFolders, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getRootNoteFolders();
      return data || [];
    },
    { initialData: [] as NoteFolder[] }
  );

  return { rootFolders, loading, refetchRootFolders: refetch };
}

export function useFolderChildren(folderId: number | undefined) {
  const { data: children, loading, refetch } = useAsyncData(
    async () => {
      if (!folderId) return [];
      const { data } = await getNoteFolderChildren(folderId);
      return data || [];
    },
    { initialData: [] as NoteFolder[], deps: [folderId], enabled: !!folderId }
  );

  return { children, loading, refetchChildren: refetch };
}
