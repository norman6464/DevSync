import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyRoadmaps,
  createRoadmap as createRoadmapApi,
  updateRoadmap as updateRoadmapApi,
  deleteRoadmap as deleteRoadmapApi,
  copyRoadmap as copyRoadmapApi,
  getRoadmap,
  createStep as createStepApi,
  updateStep as updateStepApi,
  deleteStep as deleteStepApi,
  reorderSteps as reorderStepsApi,
  type Roadmap,
  type CreateRoadmapRequest,
  type UpdateRoadmapRequest,
  type CreateStepRequest,
  type UpdateStepRequest,
} from '../api/roadmaps';
import { useAsyncData } from './useAsyncData';

export function useRoadmaps() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: roadmaps, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyRoadmaps();
      return data || [];
    },
    { initialData: [] as Roadmap[] }
  );

  const [localRoadmaps, setLocalRoadmaps] = useState<Roadmap[] | null>(null);
  const currentRoadmaps = localRoadmaps ?? roadmaps;

  const setRoadmaps = useCallback((updater: Roadmap[] | ((prev: Roadmap[]) => Roadmap[])) => {
    setLocalRoadmaps(prev => {
      const current = prev ?? roadmaps;
      return typeof updater === 'function' ? updater(current) : updater;
    });
  }, [roadmaps]);

  const handleCreate = useCallback(async (data: CreateRoadmapRequest) => {
    setSaving(true);
    try {
      const { data: newRoadmap } = await createRoadmapApi(data);
      setRoadmaps(prev => [newRoadmap, ...prev]);
      toast.success(t('roadmaps.created'));
      return newRoadmap;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, setRoadmaps]);

  const handleUpdate = useCallback(async (roadmapId: number, data: UpdateRoadmapRequest) => {
    try {
      const { data: updated } = await updateRoadmapApi(roadmapId, data);
      setRoadmaps(prev => prev.map(r => r.id === updated.id ? updated : r));
      toast.success(t('roadmaps.updated'));
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, setRoadmaps]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('roadmaps.confirmDelete'))) return false;
    try {
      await deleteRoadmapApi(id);
      setRoadmaps(prev => prev.filter(r => r.id !== id));
      toast.success(t('roadmaps.deleted'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, setRoadmaps]);

  const handleCopy = useCallback(async (id: number) => {
    try {
      const { data: copied } = await copyRoadmapApi(id);
      setRoadmaps(prev => [copied, ...prev]);
      toast.success(t('roadmaps.copied'));
      return copied;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, setRoadmaps]);

  const activeRoadmaps = currentRoadmaps.filter(r => r.status === 'active');
  const completedRoadmaps = currentRoadmaps.filter(r => r.status === 'completed');

  return {
    roadmaps: currentRoadmaps,
    loading,
    saving,
    activeRoadmaps,
    completedRoadmaps,
    createRoadmap: handleCreate,
    updateRoadmap: handleUpdate,
    deleteRoadmap: handleDelete,
    copyRoadmap: handleCopy,
    refetch,
  };
}

export function useRoadmapDetail(roadmapId: number | null) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: roadmap, loading, refetch } = useAsyncData(
    async () => {
      if (!roadmapId) return null;
      const { data } = await getRoadmap(roadmapId);
      return data;
    },
    { deps: [roadmapId], enabled: roadmapId !== null, initialData: null as Roadmap | null }
  );

  const handleCreateStep = useCallback(async (data: CreateStepRequest) => {
    if (!roadmapId) return null;
    setSaving(true);
    try {
      await createStepApi(roadmapId, data);
      await refetch();
      toast.success(t('roadmaps.stepCreated'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [roadmapId, refetch, t]);

  const handleUpdateStep = useCallback(async (stepId: number, data: UpdateStepRequest) => {
    if (!roadmapId) return null;
    try {
      await updateStepApi(roadmapId, stepId, data);
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [roadmapId, refetch, t]);

  const handleDeleteStep = useCallback(async (stepId: number) => {
    if (!roadmapId) return false;
    if (!confirm(t('roadmaps.confirmDeleteStep'))) return false;
    try {
      await deleteStepApi(roadmapId, stepId);
      await refetch();
      toast.success(t('roadmaps.stepDeleted'));
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [roadmapId, refetch, t]);

  const handleReorderSteps = useCallback(async (orders: Array<{ step_id: number; order_index: number }>) => {
    if (!roadmapId) return false;
    try {
      await reorderStepsApi(roadmapId, { orders });
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [roadmapId, refetch, t]);

  return {
    roadmap,
    loading,
    saving,
    createStep: handleCreateStep,
    updateStep: handleUpdateStep,
    deleteStep: handleDeleteStep,
    reorderSteps: handleReorderSteps,
    refetch,
  };
}
