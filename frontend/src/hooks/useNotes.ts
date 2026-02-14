import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyNotes,
  createNote,
  updateNote,
  deleteNote,
  toggleFavorite,
  searchNotes,
  type Note,
  type CreateNoteRequest,
  type UpdateNoteRequest,
} from '../api/notes';
import { useAsyncData } from './useAsyncData';

export function useNotes(page = 1, limit = 20) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: notesData, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyNotes(page, limit);
      return data || { data: [], total: 0, page: 1, limit: 20 };
    },
    { initialData: { data: [], total: 0, page: 1, limit: 20 } }
  );

  const [localNotes, setLocalNotes] = useState<Note[] | null>(null);
  const currentNotes = localNotes ?? notesData.data;

  // Sync localNotes when remote data changes
  const setNotes = useCallback((updater: Note[] | ((prev: Note[]) => Note[])) => {
    setLocalNotes(prev => {
      const current = prev ?? notesData.data;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [notesData.data]);

  const handleCreate = useCallback(async (data: CreateNoteRequest) => {
    setSaving(true);
    try {
      const { data: newNote } = await createNote(data);
      setNotes(prev => [newNote, ...prev]);
      toast.success(t('notes.created'));
      return newNote;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setNotes]);

  const handleUpdate = useCallback(async (noteId: number, data: UpdateNoteRequest) => {
    setSaving(true);
    try {
      const { data: updated } = await updateNote(noteId, data);
      setNotes(prev => prev.map(n => n.id === updated.id ? updated : n));
      toast.success(t('notes.updated'));
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setNotes]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('notes.confirmDelete'))) return false;
    try {
      await deleteNote(id);
      setNotes(prev => prev.filter(n => n.id !== id));
      toast.success(t('notes.deleted'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setNotes]);

  const handleToggleFavorite = useCallback(async (id: number) => {
    try {
      await toggleFavorite(id);
      setNotes(prev => prev.map(n =>
        n.id === id ? { ...n, is_favorite: !n.is_favorite } : n
      ));
      toast.success(t('notes.favoriteToggled'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setNotes]);

  const favoriteNotes = currentNotes.filter(n => n.is_favorite);

  return {
    notes: currentNotes,
    total: notesData.total,
    page: notesData.page,
    limit: notesData.limit,
    loading,
    saving,
    favoriteNotes,
    createNote: handleCreate,
    updateNote: handleUpdate,
    deleteNote: handleDelete,
    toggleFavorite: handleToggleFavorite,
    refetch,
  };
}

export function useNoteSearch(query: string, page = 1, limit = 20) {
  const { data: searchResults, loading } = useAsyncData(
    async () => {
      if (!query.trim()) {
        return { data: [], total: 0, page: 1, limit: 20 };
      }
      const { data } = await searchNotes(query, page, limit);
      return data || { data: [], total: 0, page: 1, limit: 20 };
    },
    { initialData: { data: [], total: 0, page: 1, limit: 20 }, deps: [query, page, limit] }
  );

  return {
    notes: searchResults.data,
    total: searchResults.total,
    page: searchResults.page,
    limit: searchResults.limit,
    loading,
  };
}
