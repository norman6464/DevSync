import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Plus, ArrowDownWideNarrow, Tag } from 'lucide-react';
import { useNoteForm } from '../hooks';
import { PageLoader, SearchInput } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import NoteFormPanel from '../components/notes/NoteFormPanel';
import NoteCard from '../components/notes/NoteCard';
import { buttonPrimaryClass } from '../constants/styles';

export default function NotesPage() {
  const { t } = useTranslation();
  const {
    filteredNotes, loading, saving,
    showForm, setShowForm, editingNote, searchQuery, setSearchQuery, sortBy, setSortBy,
    filterTag, setFilterTag, allTags,
    title, setTitle, content, setContent, tags, setTags,
    resetForm, handleSubmit, handleEdit, deleteNote, toggleFavorite,
  } = useNoteForm();

  const handleToggleForm = useCallback(() => setShowForm((prev) => !prev), [setShowForm]);

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

      {/* Sort */}
      <div className="flex items-center gap-3 mb-6">
        <ArrowDownWideNarrow className="w-4 h-4 text-gray-400" />
        <div className="flex flex-wrap gap-2">
          {([
            { value: 'latest', label: 'notes.sortLatest' },
            { value: 'oldest', label: 'notes.sortOldest' },
            { value: 'updated', label: 'notes.sortUpdated' },
            { value: 'favorites_first', label: 'notes.sortFavorites' },
          ] as const).map((opt) => (
            <button
              key={opt.value}
              onClick={() => setSortBy(opt.value)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                sortBy === opt.value
                  ? 'bg-blue-500/20 text-blue-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              {t(opt.label)}
            </button>
          ))}
        </div>
      </div>

      {/* Tag Filter */}
      {allTags.length > 0 && (
        <div className="flex items-center gap-3 mb-6">
          <Tag className="w-4 h-4 text-gray-400" />
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setFilterTag('')}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                filterTag === ''
                  ? 'bg-green-500/20 text-green-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              {t('notes.allTags')}
            </button>
            {allTags.map(tag => (
              <button
                key={tag}
                onClick={() => setFilterTag(filterTag === tag ? '' : tag)}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                  filterTag === tag
                    ? 'bg-green-500/20 text-green-400'
                    : 'bg-gray-800/50 text-gray-400 hover:text-white'
                }`}
              >
                {tag}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Create/Edit Form */}
      {showForm && (
        <NoteFormPanel
          editingNote={!!editingNote}
          title={title}
          setTitle={setTitle}
          content={content}
          setContent={setContent}
          tags={tags}
          setTags={setTags}
          saving={saving}
          onSubmit={handleSubmit}
          onCancel={resetForm}
        />
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
            <NoteCard
              key={note.id}
              note={note}
              onToggleFavorite={toggleFavorite}
              onEdit={handleEdit}
              onDelete={deleteNote}
            />
          ))}
        </div>
      )}
    </div>
  );
}
