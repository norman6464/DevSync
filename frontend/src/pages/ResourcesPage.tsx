import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { LearningResource, CreateResourceRequest, ResourceCategory, ResourceDifficulty } from '../types/resource';
import {
  getPublicResources,
  getSavedResources,
  searchResources,
  createResource,
  updateResource,
  deleteResource,
  likeResource,
  unlikeResource,
  saveResource,
  unsaveResource,
} from '../api/resources';
import ResourceCard from '../components/resources/ResourceCard';
import ResourceForm from '../components/resources/ResourceForm';
import LoadingSpinner from '../components/common/LoadingSpinner';

type TabType = 'explore' | 'saved' | 'mine';

const categories: (ResourceCategory | '')[] = ['', 'book', 'video', 'article', 'course', 'tutorial', 'podcast', 'tool', 'other'];
const difficulties: (ResourceDifficulty | '')[] = ['', 'beginner', 'intermediate', 'advanced'];

export default function ResourcesPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [tab, setTab] = useState<TabType>('explore');
  const [resources, setResources] = useState<LearningResource[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [editingResource, setEditingResource] = useState<LearningResource | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<ResourceCategory | ''>('');
  const [difficultyFilter, setDifficultyFilter] = useState<ResourceDifficulty | ''>('');
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const limit = 20;

  useEffect(() => {
    fetchResources();
  }, [tab, categoryFilter, difficultyFilter, page]);

  const fetchResources = async () => {
    setLoading(true);
    try {
      let data: { resources: LearningResource[]; total: number };

      if (tab === 'saved') {
        data = await getSavedResources(limit, page * limit);
      } else if (tab === 'mine' && user) {
        const res = await import('../api/resources');
        const myResources = await res.getResourcesByUserId(user.id);
        data = { resources: myResources, total: myResources.length };
      } else {
        if (searchQuery.trim()) {
          data = await searchResources(searchQuery, limit, page * limit);
        } else {
          data = await getPublicResources(limit, page * limit, categoryFilter, difficultyFilter);
        }
      }

      setResources(data.resources);
      setTotal(data.total);
    } catch (error) {
      console.error('Failed to fetch resources:', error);
      toast.error(t('errors.somethingWrong'));
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    setPage(0);
    fetchResources();
  };

  const handleCreate = async (data: CreateResourceRequest) => {
    setSaving(true);
    try {
      const newResource = await createResource(data);
      setResources([newResource, ...resources]);
      setShowForm(false);
      toast.success(t('resources.createSuccess'));
    } catch (error) {
      console.error('Failed to create resource:', error);
      toast.error(t('resources.createFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleUpdate = async (data: CreateResourceRequest) => {
    if (!editingResource) return;
    setSaving(true);
    try {
      const updated = await updateResource(editingResource.id, data);
      setResources(resources.map(r => r.id === updated.id ? updated : r));
      setEditingResource(null);
      toast.success(t('resources.updateSuccess'));
    } catch (error) {
      console.error('Failed to update resource:', error);
      toast.error(t('resources.updateFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (resource: LearningResource) => {
    if (!confirm(t('resources.confirmDelete'))) return;
    try {
      await deleteResource(resource.id);
      setResources(resources.filter(r => r.id !== resource.id));
      toast.success(t('resources.deleteSuccess'));
    } catch (error) {
      console.error('Failed to delete resource:', error);
      toast.error(t('resources.deleteFailed'));
    }
  };

  const handleLike = async (id: number) => {
    try {
      await likeResource(id);
    } catch (error) {
      console.error('Failed to like resource:', error);
    }
  };

  const handleUnlike = async (id: number) => {
    try {
      await unlikeResource(id);
    } catch (error) {
      console.error('Failed to unlike resource:', error);
    }
  };

  const handleSave = async (id: number) => {
    try {
      await saveResource(id);
    } catch (error) {
      console.error('Failed to save resource:', error);
    }
  };

  const handleUnsave = async (id: number) => {
    try {
      await unsaveResource(id);
    } catch (error) {
      console.error('Failed to unsave resource:', error);
    }
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('resources.pageTitle')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('resources.pageSubtitle')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('resources.addResource')}
        </button>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 mb-6">
        <button
          onClick={() => { setTab('explore'); setPage(0); }}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'explore' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.explore')}
        </button>
        <button
          onClick={() => { setTab('saved'); setPage(0); }}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'saved' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.savedTab')}
        </button>
        <button
          onClick={() => { setTab('mine'); setPage(0); }}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            tab === 'mine' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
          }`}
        >
          {t('resources.myResources')}
        </button>
      </div>

      {/* Filters (only for explore tab) */}
      {tab === 'explore' && (
        <div className="flex flex-wrap gap-4 mb-6">
          <div className="flex-1 min-w-[200px]">
            <div className="flex gap-2">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder={t('resources.searchPlaceholder')}
                className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
              />
              <button
                onClick={handleSearch}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
              >
                {t('common.search')}
              </button>
            </div>
          </div>
          <select
            value={categoryFilter}
            onChange={(e) => { setCategoryFilter(e.target.value as ResourceCategory | ''); setPage(0); }}
            className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
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
            onChange={(e) => { setDifficultyFilter(e.target.value as ResourceDifficulty | ''); setPage(0); }}
            className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
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
      {(showForm || editingResource) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingResource ? t('resources.editResource') : t('resources.newResource')}
            </h2>
            <ResourceForm
              resource={editingResource || undefined}
              onSubmit={editingResource ? handleUpdate : handleCreate}
              onCancel={() => {
                setShowForm(false);
                setEditingResource(null);
              }}
              loading={saving}
            />
          </div>
        </div>
      )}

      {/* Content */}
      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <LoadingSpinner />
        </div>
      ) : resources.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 0 0 6 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 0 1 6 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 0 1 6-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0 0 18 18a8.967 8.967 0 0 0-6 2.292m0-14.25v14.25" />
            </svg>
          </div>
          <p className="text-gray-400">
            {tab === 'saved' ? t('resources.noSavedResources') :
             tab === 'mine' ? t('resources.noMyResources') :
             t('resources.noResources')}
          </p>
          {tab !== 'saved' && (
            <button
              onClick={() => setShowForm(true)}
              className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
            >
              {t('resources.addFirstResource')}
            </button>
          )}
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {resources.map(resource => (
              <ResourceCard
                key={resource.id}
                resource={resource}
                isOwner={user?.id === resource.user_id}
                showUser={tab !== 'mine'}
                onLike={() => handleLike(resource.id)}
                onUnlike={() => handleUnlike(resource.id)}
                onSave={() => handleSave(resource.id)}
                onUnsave={() => handleUnsave(resource.id)}
                onEdit={() => setEditingResource(resource)}
                onDelete={() => handleDelete(resource)}
              />
            ))}
          </div>

          {/* Pagination */}
          {total > limit && (
            <div className="flex justify-center gap-2 mt-8">
              <button
                onClick={() => setPage(p => Math.max(0, p - 1))}
                disabled={page === 0}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.previous')}
              </button>
              <span className="px-4 py-2 text-gray-400">
                {page + 1} / {Math.ceil(total / limit)}
              </span>
              <button
                onClick={() => setPage(p => p + 1)}
                disabled={(page + 1) * limit >= total}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.next')}
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
