import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Plus, ArrowDownWideNarrow, Tag, List, LayoutGrid, Sparkles } from 'lucide-react';
import { useNoteForm } from '../hooks';
import { PageLoader, SearchInput } from '../components/common';
import EmptyState from '../components/common/EmptyState';
import NoteFormPanel from '../components/notes/NoteFormPanel';
import NoteCard from '../components/notes/NoteCard';
import NoteTemplatesModal from '../components/notes/NoteTemplatesModal';
import { buttonPrimaryClass, buttonSecondaryClass } from '../constants/styles';
import { NOTES_SORT_OPTIONS } from '../constants/notes';
import type { NoteTemplate } from '../constants/noteTemplates';

export default function NotesPage() {
  const { t } = useTranslation();
  const [showTemplates, setShowTemplates] = useState(false);
  const {
    filteredNotes, loading, saving,
    showForm, setShowForm, editingNote, searchQuery, setSearchQuery, sortBy, setSortBy,
    filterTag, setFilterTag, allTags, viewMode, setViewMode,
    title, setTitle, content, setContent, tags, setTags,
    resetForm, handleSubmit, handleEdit, deleteNote, toggleFavorite,
  } = useNoteForm();

  const handleToggleForm = useCallback(() => setShowForm((prev) => !prev), [setShowForm]);

  const handleTemplateSelect = useCallback((template: NoteTemplate) => {
    setTitle(template.title);
    setContent(template.content);
    setTags(template.tags);
    setShowTemplates(false);
    setShowForm(true);
  }, [setTitle, setContent, setTags, setShowForm]);

  if (loading) return <PageLoader />;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <BookOpen className="w-8 h-8 text-blue-500" />
          <div className="flex items-baseline gap-2">
            <h1 className="text-3xl font-bold">{t('notes.title')}</h1>
            <span className="text-lg text-gray-500">({filteredNotes.length})</span>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowTemplates(true)}
            className={`${buttonSecondaryClass} flex items-center gap-2`}
          >
            <Sparkles className="w-5 h-5" />
            テンプレートから作成
          </button>
          <button
            onClick={handleToggleForm}
            className={`${buttonPrimaryClass} flex items-center gap-2`}
          >
            <Plus className="w-5 h-5" />
            {t('notes.createNote')}
          </button>
        </div>
      </div>

      {/* Search Bar */}
      <div className="mb-6">
        <SearchInput
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder={t('notes.searchPlaceholder')}
        />
      </div>

      {/* Sort & View Toggle */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <ArrowDownWideNarrow className="w-4 h-4 text-gray-400" />
          <div className="flex flex-wrap gap-2">
            {NOTES_SORT_OPTIONS.map((opt) => (
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
        <div className="flex gap-1">
          <button
            onClick={() => setViewMode('list')}
            className={`p-2 rounded-lg transition-colors ${viewMode === 'list' ? 'bg-blue-500/20 text-blue-400' : 'text-gray-400 hover:text-white'}`}
            title={t('notes.viewList')}
          >
            <List className="w-4 h-4" />
          </button>
          <button
            onClick={() => setViewMode('grid')}
            className={`p-2 rounded-lg transition-colors ${viewMode === 'grid' ? 'bg-blue-500/20 text-blue-400' : 'text-gray-400 hover:text-white'}`}
            title={t('notes.viewGrid')}
          >
            <LayoutGrid className="w-4 h-4" />
          </button>
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

      {/* Templates Modal */}
      <NoteTemplatesModal
        isOpen={showTemplates}
        onSelect={handleTemplateSelect}
        onClose={() => setShowTemplates(false)}
      />

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
        <div className={viewMode === 'grid' ? 'grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3' : 'grid gap-4'}>
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
