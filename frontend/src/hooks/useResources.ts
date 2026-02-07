import { useState, useCallback } from 'react';
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
import { useAsyncData } from './useAsyncData';

type TabType = 'explore' | 'saved' | 'mine';

export function useResources() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [tab, setTab] = useState<TabType>('explore');
  const [saving, setSaving] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [categoryFilter, setCategoryFilter] = useState<ResourceCategory | ''>('');
  const [difficultyFilter, setDifficultyFilter] = useState<ResourceDifficulty | ''>('');
  const [page, setPage] = useState(0);
  const limit = 20;

  const { data, loading, refetch } = useAsyncData(
    async () => {
      let result: { resources: LearningResource[]; total: number };

      if (tab === 'saved') {
        result = await getSavedResources(limit, page * limit);
      } else if (tab === 'mine' && user) {
        const res = await import('../api/resources');
        const myResources = await res.getResourcesByUserId(user.id);
        result = { resources: myResources, total: myResources.length };
      } else {
        if (searchQuery.trim()) {
          result = await searchResources(searchQuery, limit, page * limit);
        } else {
          result = await getPublicResources(limit, page * limit, categoryFilter, difficultyFilter);
        }
      }
      return result;
    },
    { deps: [tab, categoryFilter, difficultyFilter, page] }
  );

  const resources = data?.resources ?? [];
  const total = data?.total ?? 0;

  const [localResources, setLocalResources] = useState<LearningResource[] | null>(null);
  const currentResources = localResources ?? resources;

  const handleSearch = useCallback(() => {
    setPage(0);
    refetch();
  }, [refetch]);

  const handleCreate = useCallback(async (reqData: CreateResourceRequest) => {
    setSaving(true);
    try {
      const newResource = await createResource(reqData);
      setLocalResources(prev => [newResource, ...(prev ?? resources)]);
      toast.success(t('resources.createSuccess'));
      return newResource;
    } catch {
      toast.error(t('resources.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, resources]);

  const handleUpdate = useCallback(async (resourceId: number, reqData: CreateResourceRequest) => {
    setSaving(true);
    try {
      const updated = await updateResource(resourceId, reqData);
      setLocalResources(prev => (prev ?? resources).map(r => r.id === updated.id ? updated : r));
      toast.success(t('resources.updateSuccess'));
      return updated;
    } catch {
      toast.error(t('resources.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, resources]);

  const handleDelete = useCallback(async (resource: LearningResource) => {
    if (!confirm(t('resources.confirmDelete'))) return false;
    try {
      await deleteResource(resource.id);
      setLocalResources(prev => (prev ?? resources).filter(r => r.id !== resource.id));
      toast.success(t('resources.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('resources.deleteFailed'));
      return false;
    }
  }, [t, resources]);

  const handleLike = useCallback(async (id: number) => {
    try { await likeResource(id); } catch { /* noop */ }
  }, []);

  const handleUnlike = useCallback(async (id: number) => {
    try { await unlikeResource(id); } catch { /* noop */ }
  }, []);

  const handleSave = useCallback(async (id: number) => {
    try { await saveResource(id); } catch { /* noop */ }
  }, []);

  const handleUnsave = useCallback(async (id: number) => {
    try { await unsaveResource(id); } catch { /* noop */ }
  }, []);

  const changeTab = useCallback((newTab: TabType) => {
    setTab(newTab);
    setPage(0);
    setLocalResources(null);
  }, []);

  return {
    resources: currentResources,
    total,
    loading,
    saving,
    tab,
    setTab: changeTab,
    searchQuery,
    setSearchQuery,
    categoryFilter,
    setCategoryFilter: (v: ResourceCategory | '') => { setCategoryFilter(v); setPage(0); },
    difficultyFilter,
    setDifficultyFilter: (v: ResourceDifficulty | '') => { setDifficultyFilter(v); setPage(0); },
    page,
    setPage,
    limit,
    handleSearch,
    createResource: handleCreate,
    updateResource: handleUpdate,
    deleteResource: handleDelete,
    likeResource: handleLike,
    unlikeResource: handleUnlike,
    saveResource: handleSave,
    unsaveResource: handleUnsave,
    refetch,
  };
}
