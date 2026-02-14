import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Star, Search, Plus, Edit, Trash2 } from 'lucide-react';
import { useNotes } from '../hooks';
import { PageLoader } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import type { Note } from '../api/notes';

export default function NotesPage() {
  const { t } = useTranslation();
  const { notes, loading, saving, favoriteNotes, createNote, updateNote, deleteNote, toggleFavorite } = useNotes();

  const [showForm, setShowForm] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  // Form state
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState('');

  const resetForm = () => {
    setTitle('');
    setContent('');
    setTags('');
    setEditingNote(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;

    if (editingNote) {
      const result = await updateNote(editingNote.id, { title, content, tags });
      if (result) resetForm();
    } else {
      const result = await createNote({ title, content, tags });
      if (result) resetForm();
    }
  };

  const handleEdit = (note: Note) => {
    setEditingNote(note);
    setTitle(note.title);
    setContent(note.content);
    setTags(note.tags);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    await deleteNote(id);
  };

  const handleToggleFavorite = async (id: number) => {
    await toggleFavorite(id);
  };

  const filteredNotes = searchQuery
    ? notes.filter(note =>
        note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        note.content.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : notes;

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <BookOpen className="w-8 h-8 text-blue-500" />
          <h1 className="text-3xl font-bold">{t('notes.title')}</h1>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-5 h-5" />
          {t('notes.createNote')}
        </button>
      </div>

      {/* Search Bar */}
      <div className="mb-6">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder={t('notes.searchPlaceholder')}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Create/Edit Form */}
      {showForm && (
        <div className="mb-8 p-6 bg-gray-800 border border-gray-700 rounded-md">
          <h2 className="text-xl font-semibold mb-4">
            {editingNote ? t('notes.editNote') : t('notes.createNote')}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">{t('notes.noteTitle')}</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">{t('notes.noteContent')}</label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={8}
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">{t('notes.tags')}</label>
              <input
                type="text"
                value={tags}
                onChange={(e) => setTags(e.target.value)}
                placeholder="React, TypeScript, Go"
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="flex gap-3">
              <button
                type="submit"
                disabled={saving}
                className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                {saving ? t('common.saving') : t('common.save')}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="px-6 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600 transition-colors"
              >
                {t('common.cancel')}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Notes List */}
      {filteredNotes.length === 0 ? (
        <EmptyState
          icon={BookOpen}
          title={searchQuery ? t('notes.noSearchResults') : t('notes.noNotes')}
          description=""
        />
      ) : (
        <div className="grid gap-4">
          {filteredNotes.map((note) => (
            <div
              key={note.id}
              className="p-6 bg-gray-800 border border-gray-700 rounded-md hover:border-gray-600 transition-colors"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="text-xl font-semibold">{note.title}</h3>
                    {note.is_favorite && (
                      <Star className="w-5 h-5 fill-yellow-500 text-yellow-500" />
                    )}
                  </div>
                  <p className="text-gray-400 mb-3 line-clamp-2">{note.content}</p>
                  {note.tags && (
                    <div className="flex flex-wrap gap-2 mb-3">
                      {note.tags.split(',').map((tag, idx) => (
                        <span
                          key={idx}
                          className="px-2 py-1 text-xs bg-gray-700 text-gray-300 rounded"
                        >
                          {tag.trim()}
                        </span>
                      ))}
                    </div>
                  )}
                  <p className="text-sm text-gray-500">
                    {t('notes.lastUpdated')}: {new Date(note.updated_at).toLocaleString()}
                  </p>
                </div>
                <div className="flex gap-2 ml-4">
                  <button
                    onClick={() => handleToggleFavorite(note.id)}
                    className="p-2 text-gray-400 hover:text-yellow-500 transition-colors"
                    title={t('notes.favorites')}
                  >
                    <Star className={`w-5 h-5 ${note.is_favorite ? 'fill-yellow-500 text-yellow-500' : ''}`} />
                  </button>
                  <button
                    onClick={() => handleEdit(note)}
                    className="p-2 text-gray-400 hover:text-blue-500 transition-colors"
                    title={t('common.edit')}
                  >
                    <Edit className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(note.id)}
                    className="p-2 text-gray-400 hover:text-red-500 transition-colors"
                    title={t('common.delete')}
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
