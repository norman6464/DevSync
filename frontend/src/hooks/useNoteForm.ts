import { useState, useCallback, useMemo } from 'react';
import { useNotes } from './useNotes';
import type { Note } from '../api/notes';

export function useNoteForm() {
  const { notes, loading, saving, favoriteNotes, createNote, updateNote, deleteNote, toggleFavorite } = useNotes();

  const [showForm, setShowForm] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');

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

  const filteredNotes = useMemo(() =>
    searchQuery
      ? notes.filter(note =>
          note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
          note.content.toLowerCase().includes(searchQuery.toLowerCase())
        )
      : notes,
    [notes, searchQuery]
  );

  return {
    // Data
    notes, favoriteNotes, filteredNotes, loading, saving,
    // Form state
    showForm, setShowForm, editingNote, searchQuery, setSearchQuery,
    title, setTitle, content, setContent, tags, setTags,
    // Actions
    resetForm, handleSubmit, handleEdit, deleteNote, toggleFavorite,
  };
}
