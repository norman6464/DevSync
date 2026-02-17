import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import type { LearningResource, ResourceCategory, ResourceDifficulty } from '../types/resource';
import { useResources, useDebounce } from '../hooks';
import ResourceCard from '../components/resources/ResourceCard';
import ResourceForm from '../components/resources/ResourceForm';
import { ResourceCardSkeleton } from '../components/common/Skeleton';
import EmptyState from '../components/common/EmptyState';
import { Pagination, SearchInput, PageHeader, Modal } from '../components/common';

const categories: (ResourceCategory | '')[] = ['', 'book', 'video', 'article', 'course', 'tutorial', 'podcast', 'tool', 'other'];
const difficulties: (ResourceDifficulty | '')[] = ['', 'beginner', 'intermediate', 'advanced'];

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
        <div className="flex flex-wrap gap-4 mb-6">
          <div className="flex-1 min-w-[200px]">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              onSearch={handleSearch}
              placeholder={t('resources.searchPlaceholder')}
              showButton
            />
          </div>
          <select
            value={categoryFilter}
            onChange={(e) => setCategoryFilter(e.target.value as ResourceCategory | '')}
            className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:ring-2 focus:ring-gray-500 focus:border-transparent"
          >
            <option value="">{t('resources.allCategories')}</option>
            {categories.slice(1).map(cat => (
              <option key={cat} value={cat}>
                {t(`resources.categories.${cat}`)}
              </option>
            ))}
          </select>
          <select
            value={difficultyFilter}
            onChange={(e) => setDifficultyFilter(e.target.value as ResourceDifficulty | '')}
            className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:ring-2 focus:ring-gray-500 focus:border-transparent"
          >
            <option value="">{t('resources.allDifficulties')}</option>
            {difficulties.slice(1).map(diff => (
              <option key={diff} value={diff}>
                {t(`resources.difficulty.${diff}`)}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingResource}
        onClose={() => { setShowForm(false); setEditingResource(null); }}
        title={editingResource ? t('resources.editResource') : t('resources.newResource')}
      >
        <ResourceForm
          resource={editingResource || undefined}
          onSubmit={async (data) => {
            if (editingResource) {
              const result = await updateResource(editingResource.id, data);
              if (result) setEditingResource(null);
            } else {
              const result = await createResource(data);
              if (result) setShowForm(false);
            }
          }}
          onCancel={() => {
            setShowForm(false);
            setEditingResource(null);
          }}
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
