import { useState, useCallback, useMemo } from 'react';
import { useNotes } from './useNotes';
import type { Note } from '../api/notes';

export function useNoteForm() {
  const { notes, loading, saving, favoriteNotes, createNote, updateNote, deleteNote, toggleFavorite } = useNotes();

  const [showForm, setShowForm] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'latest' | 'oldest' | 'updated' | 'favorites_first'>('latest');
  const [filterTag, setFilterTag] = useState('');
  const [showFavoritesOnly, setShowFavoritesOnly] = useState(false);
  const [viewMode, setViewMode] = useState<'list' | 'grid'>('list');
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');

  const allTags = useMemo(() => {
    const tagSet = new Set<string>();
    notes.forEach(note => {
      if (note.tags) {
        note.tags.split(',').forEach(t => {
          const trimmed = t.trim();
          if (trimmed) tagSet.add(trimmed);
        });
      }
    });
    return Array.from(tagSet).sort();
  }, [notes]);

  const resetForm = useCallback(() => {
    setTitle('');
    setContent('');
    setTags('');
    setEditingNote(null);
    setShowForm(false);
  }, []);

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;

    if (editingNote) {
      const result = await updateNote(editingNote.id, { title, content, tags });
      if (result) resetForm();
    } else {
      const result = await createNote({ title, content, tags });
      if (result) resetForm();
    }
  }, [title, content, tags, editingNote, updateNote, createNote, resetForm]);

  const handleEdit = useCallback((note: Note) => {
    setEditingNote(note);
    setTitle(note.title);
    setContent(note.content);
    setTags(note.tags);
    setShowForm(true);
  }, []);

  const filteredNotes = useMemo(() => {
    let filtered = searchQuery
      ? notes.filter(note =>
          note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          note.content.toLowerCase().includes(searchQuery.toLowerCase())
        )
      : notes;
    if (showFavoritesOnly) {
      filtered = filtered.filter(note => note.is_favorite);
    }
    if (filterTag) {
      filtered = filtered.filter(note =>
        note.tags?.split(',').some(t => t.trim() === filterTag)
      );
    }
    return [...filtered].sort((a, b) => {
      switch (sortBy) {
        case 'oldest':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        case 'updated':
          return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime();
        case 'favorites_first':
          if (a.is_favorite !== b.is_favorite) return a.is_favorite ? -1 : 1;
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
      }
    });
  }, [notes, searchQuery, sortBy, filterTag, showFavoritesOnly]);

  return {
    // Data
    notes, favoriteNotes, filteredNotes, loading, saving,
    // Form state
    showForm, setShowForm, editingNote, searchQuery, setSearchQuery, sortBy, setSortBy, filterTag, setFilterTag, showFavoritesOnly, setShowFavoritesOnly, allTags, viewMode, setViewMode,
    title, setTitle, content, setContent, tags, setTags,
    // Actions
    resetForm, handleSubmit, handleEdit, deleteNote, toggleFavorite,
  };
}
