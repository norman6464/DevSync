import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Star, Plus, Edit, Trash2 } from 'lucide-react';
import { useNoteForm } from '../hooks';
import { PageLoader, SearchInput } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import { buttonPrimaryClass, buttonSecondaryClass } from '../constants/styles';

export default function NotesPage() {
  const { t } = useTranslation();
  const {
    filteredNotes, loading, saving,
    showForm, setShowForm, editingNote, searchQuery, setSearchQuery,
    title, setTitle, content, setContent, tags, setTags,
    resetForm, handleSubmit, handleEdit, deleteNote, toggleFavorite,
  } = useNoteForm();

  const handleToggleForm = useCallback(() => setShowForm((prev) => !prev), [setShowForm]);
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTitle(e.target.value), [setTitle]);
  const handleContentChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setContent(e.target.value), [setContent]);
  const handleTagsChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTags(e.target.value), [setTags]);

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <BookOpen className="w-8 h-8 text-blue-500" />
          <h1 className="text-3xl font-bold">{t('notes.title')}</h1>
        </div>
        <button
          onClick={handleToggleForm}
          className={`${buttonPrimaryClass} flex items-center gap-2`}
        >
          <Plus className="w-5 h-5" />
          {t('notes.createNote')}
        </button>
      </div>

      {/* Search Bar */}
      <div className="mb-6">
        <SearchInput
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder={t('notes.searchPlaceholder')}
        />
      </div>

      {/* Create/Edit Form */}
      {showForm && (
        <div className="mb-8 p-6 bg-gray-800 border border-gray-700 rounded-md">
          <h2 className="text-xl font-semibold mb-4">
            {editingNote ? t('notes.editNote') : t('notes.createNote')}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="note-title" className="block text-sm font-medium mb-2">{t('notes.noteTitle')}</label>
              <input
                id="note-title"
                type="text"
                value={title}
                onChange={handleTitleChange}
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label htmlFor="note-content" className="block text-sm font-medium mb-2">{t('notes.noteContent')}</label>
              <textarea
                id="note-content"
                value={content}
                onChange={handleContentChange}
                rows={8}
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
            <div>
              <label htmlFor="note-tags" className="block text-sm font-medium mb-2">{t('notes.tags')}</label>
              <input
                id="note-tags"
                type="text"
                value={tags}
                onChange={handleTagsChange}
                placeholder={t('notes.tagsPlaceholder')}
                className="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div className="flex gap-3">
              <button
                type="submit"
                disabled={saving}
                className={`${buttonPrimaryClass} disabled:opacity-50`}
              >
                {saving ? t('common.saving') : t('common.save')}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className={buttonSecondaryClass}
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
                    onClick={() => toggleFavorite(note.id)}
                    className="p-2 text-gray-400 hover:text-yellow-500 transition-colors"
                    aria-label={t('notes.favorites')}
                  >
                    <Star aria-hidden="true" className={`w-5 h-5 ${note.is_favorite ? 'fill-yellow-500 text-yellow-500' : ''}`} />
                  </button>
                  <button
                    onClick={() => handleEdit(note)}
                    className="p-2 text-gray-400 hover:text-blue-500 transition-colors"
                    aria-label={t('common.edit')}
                  >
                    <Edit aria-hidden="true" className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => deleteNote(note.id)}
                    className="p-2 text-gray-400 hover:text-red-500 transition-colors"
                    aria-label={t('common.delete')}
                  >
                    <Trash2 aria-hidden="true" className="w-5 h-5" />
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
