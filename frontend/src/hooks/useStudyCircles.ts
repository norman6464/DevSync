import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  getMyCircles,
  getCircle,
  createCircle as createCircleApi,
  updateCircle as updateCircleApi,
  deleteCircle as deleteCircleApi,
  addMember as addMemberApi,
  removeMember as removeMemberApi,
  createStep as createStepApi,
  updateStep as updateStepApi,
  deleteStep as deleteStepApi,
  reorderSteps as reorderStepsApi,
  updateProgress as updateProgressApi,
  getProgress as getProgressApi,
  createCheckin as createCheckinApi,
  getCheckins as getCheckinsApi,
  getStreakRanking as getStreakRankingApi,
} from '../api/studyCircles';
import type {
  StudyCircle,
  CreateStudyCircleRequest,
  UpdateStudyCircleRequest,
  CreateStepRequest,
  UpdateStepRequest,
  StudyCircleMemberProgress,
  StudyCircleCheckin,
  CircleMemberStreak,
} from '../types/studyCircle';
import { useAsyncData } from './useAsyncData';

/** サークル一覧 + CRUD */
export function useStudyCircles() {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: circles, loading, refetch } = useAsyncData(
    async () => {
      const { data } = await getMyCircles();
      return data || [];
    },
    { initialData: [] as StudyCircle[] }
  );

  const handleCreate = useCallback(async (data: CreateStudyCircleRequest) => {
    setSaving(true);
    try {
      const { data: newCircle } = await createCircleApi(data);
      toast.success(t('studyCircle.created'));
      await refetch();
      return newCircle;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, refetch]);

  const handleUpdate = useCallback(async (id: number, data: UpdateStudyCircleRequest) => {
    try {
      const { data: updated } = await updateCircleApi(id, data);
      toast.success(t('studyCircle.updated'));
      await refetch();
      return updated;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [t, refetch]);

  const handleDelete = useCallback(async (id: number) => {
    if (!confirm(t('studyCircle.confirmDelete'))) return false;
    try {
      await deleteCircleApi(id);
      toast.success(t('studyCircle.deleted'));
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [t, refetch]);

  return {
    circles,
    loading,
    saving,
    createCircle: handleCreate,
    updateCircle: handleUpdate,
    deleteCircle: handleDelete,
    refetch,
  };
}

/** サークル詳細 + ステップ/メンバー/進捗/チェックイン操作 */
export function useStudyCircleDetail(circleId: number | null) {
  const { t } = useTranslation();
  const [saving, setSaving] = useState(false);

  const { data: circle, loading, refetch } = useAsyncData(
    async () => {
      if (!circleId) return null;
      const { data } = await getCircle(circleId);
      return data;
    },
    { deps: [circleId], enabled: circleId !== null, initialData: null as StudyCircle | null }
  );

  // メンバー管理
  const handleAddMember = useCallback(async (userId: number) => {
    if (!circleId) return false;
    try {
      await addMemberApi(circleId, userId);
      toast.success(t('studyCircle.memberAdded'));
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [circleId, refetch, t]);

  const handleRemoveMember = useCallback(async (userId: number) => {
    if (!circleId) return false;
    try {
      await removeMemberApi(circleId, userId);
      toast.success(t('studyCircle.memberRemoved'));
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [circleId, refetch, t]);

  // ステップ管理
  const handleCreateStep = useCallback(async (data: CreateStepRequest) => {
    if (!circleId) return null;
    setSaving(true);
    try {
      await createStepApi(circleId, data);
      toast.success(t('studyCircle.stepCreated'));
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [circleId, refetch, t]);

  const handleUpdateStep = useCallback(async (stepId: number, data: UpdateStepRequest) => {
    if (!circleId) return null;
    try {
      await updateStepApi(circleId, stepId, data);
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return null;
    }
  }, [circleId, refetch, t]);

  const handleDeleteStep = useCallback(async (stepId: number) => {
    if (!circleId) return false;
    if (!confirm(t('studyCircle.confirmDeleteStep'))) return false;
    try {
      await deleteStepApi(circleId, stepId);
      toast.success(t('studyCircle.stepDeleted'));
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [circleId, refetch, t]);

  const handleReorderSteps = useCallback(async (orders: Array<{ step_id: number; order_index: number }>) => {
    if (!circleId) return false;
    try {
      await reorderStepsApi(circleId, { orders });
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [circleId, refetch, t]);

  // 進捗更新
  const handleUpdateProgress = useCallback(async (stepId: number, isCompleted: boolean) => {
    if (!circleId) return false;
    try {
      await updateProgressApi(circleId, stepId, isCompleted);
      await refetch();
      return true;
    } catch {
      toast.error(t('errors.somethingWrong'));
      return false;
    }
  }, [circleId, refetch, t]);

  return {
    circle,
    loading,
    saving,
    addMember: handleAddMember,
    removeMember: handleRemoveMember,
    createStep: handleCreateStep,
    updateStep: handleUpdateStep,
    deleteStep: handleDeleteStep,
    reorderSteps: handleReorderSteps,
    updateProgress: handleUpdateProgress,
    refetch,
  };
}

/** 進捗・チェックイン・ストリーク取得 */
export function useStudyCircleActivity(circleId: number | null) {
  const { t } = useTranslation();

  const { data: progress, loading: progressLoading, refetch: refetchProgress } = useAsyncData(
    async () => {
      if (!circleId) return [];
      const { data } = await getProgressApi(circleId);
      return data || [];
    },
    { deps: [circleId], enabled: circleId !== null, initialData: [] as StudyCircleMemberProgress[] }
  );

  const { data: checkins, loading: checkinsLoading, refetch: refetchCheckins } = useAsyncData(
    async () => {
      if (!circleId) return [];
      const { data } = await getCheckinsApi(circleId);
      return data || [];
    },
    { deps: [circleId], enabled: circleId !== null, initialData: [] as StudyCircleCheckin[] }
  );

  const { data: streaks, loading: streaksLoading, refetch: refetchStreaks } = useAsyncData(
    async () => {
      if (!circleId) return [];
      const { data } = await getStreakRankingApi(circleId);
      return data || [];
    },
    { deps: [circleId], enabled: circleId !== null, initialData: [] as CircleMemberStreak[] }
  );

  const handleCheckin = useCallback(async (content: string) => {
    if (!circleId) return null;
    try {
      const { data } = await createCheckinApi(circleId, content);
      toast.success(t('studyCircle.checkedIn'));
      await refetchCheckins();
      await refetchStreaks();
      return data;
    } catch (err: unknown) {
      const status = (err as { response?: { status?: number } })?.response?.status;
      if (status === 409) {
        toast.error(t('studyCircle.alreadyCheckedIn'));
      } else {
        toast.error(t('errors.somethingWrong'));
      }
      return null;
    }
  }, [circleId, refetchCheckins, refetchStreaks, t]);

  return {
    progress,
    progressLoading,
    checkins,
    checkinsLoading,
    streaks,
    streaksLoading,
    checkin: handleCheckin,
    refetchProgress,
    refetchCheckins,
    refetchStreaks,
  };
}
