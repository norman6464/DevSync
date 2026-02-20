import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { LearningResource } from '../types/resource';
import { useResources, useDebounce } from '../hooks';
import ResourceCard from '../components/resources/ResourceCard';
import ResourceForm from '../components/resources/ResourceForm';
import ResourceFilters from '../components/resources/ResourceFilters';
import { ResourceCardSkeleton } from '../components/common/Skeleton';
import EmptyState from '../components/common/EmptyState';
import { Pagination, PageHeader, Modal } from '../components/common';

export default function ResourcesPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const {
    resources, total, loading, saving,
    tab, setTab,
    searchQuery, setSearchQuery,
    categoryFilter, setCategoryFilter,
    difficultyFilter, setDifficultyFilter,
    page, setPage, limit,
    handleSearch,
    createResource, updateResource, deleteResource,
    likeResource, unlikeResource, saveResource, unsaveResource,
  } = useResources();

  const [showForm, setShowForm] = useState(false);
  const [editingResource, setEditingResource] = useState<LearningResource | null>(null);

  const handleFormClose = useCallback(() => {
    setShowForm(false);
    setEditingResource(null);
  }, []);

  const handleFormSubmit = useCallback(async (data: Parameters<typeof createResource>[0]) => {
    if (editingResource) {
      const result = await updateResource(editingResource.id, data);
      if (result) setEditingResource(null);
    } else {
      const result = await createResource(data);
      if (result) setShowForm(false);
    }
  }, [editingResource, updateResource, createResource]);

  // デバウンス処理（300ms）
  const debouncedQuery = useDebounce(searchQuery, 300);

  // デバウンスされたクエリで自動検索
  useEffect(() => {
    if (debouncedQuery !== undefined) {
      handleSearch();
    }
  }, [debouncedQuery, handleSearch]);

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <PageHeader
        title={t('resources.pageTitle')}
        subtitle={t('resources.pageSubtitle')}
        actionLabel={t('resources.addResource')}
        onAction={() => setShowForm(true)}
      />

      {/* Tabs */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => setTab('explore')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'explore' ? 'bg-gray-700 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.explore')}
        </button>
        <button
          onClick={() => setTab('saved')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'saved' ? 'bg-gray-700 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.savedTab')}
        </button>
        <button
          onClick={() => setTab('mine')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'mine' ? 'bg-gray-700 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.myResources')}
        </button>
      </div>

      {/* Filters (only for explore tab) */}
      {tab === 'explore' && (
        <ResourceFilters
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          onSearch={handleSearch}
          categoryFilter={categoryFilter}
          onCategoryChange={setCategoryFilter}
          difficultyFilter={difficultyFilter}
          onDifficultyChange={setDifficultyFilter}
        />
      )}

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingResource}
        onClose={handleFormClose}
        title={editingResource ? t('resources.editResource') : t('resources.newResource')}
      >
        <ResourceForm
          resource={editingResource || undefined}
          onSubmit={handleFormSubmit}
          onCancel={handleFormClose}
          loading={saving}
        />
      </Modal>

      {/* Content */}
      {loading ? (
        <div className="space-y-4">
          <ResourceCardSkeleton />
          <ResourceCardSkeleton />
          <ResourceCardSkeleton />
        </div>
      ) : resources.length === 0 ? (
        <EmptyState
          icon={BookOpen}
          message={
            tab === 'saved' ? t('resources.noSavedResources') :
            tab === 'mine' ? t('resources.noMyResources') :
            t('resources.noResources')
          }
          actionLabel={tab !== 'saved' ? t('resources.addFirstResource') : undefined}
          onAction={tab !== 'saved' ? () => setShowForm(true) : undefined}
        />
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {resources.map(resource => (
              <ResourceCard
                key={resource.id}
                resource={resource}
                isOwner={user?.id === resource.user_id}
                showUser={tab !== 'mine'}
                onLike={() => likeResource(resource.id)}
                onUnlike={() => unlikeResource(resource.id)}
                onSave={() => saveResource(resource.id)}
                onUnsave={() => unsaveResource(resource.id)}
                onEdit={() => setEditingResource(resource)}
                onDelete={() => deleteResource(resource)}
              />
            ))}
          </div>

          <Pagination
            currentPage={page}
            totalItems={total}
            itemsPerPage={limit}
            onPageChange={setPage}
          />
        </>
      )}
    </div>
  );
}
