import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { Project, CreateProjectRequest } from '../types/project';
import type { GitHubRepository } from '../types/github';
import { getProjectsByUserId, createProject, updateProject, deleteProject } from '../api/projects';
import { getRepos } from '../api/github';
import { useAsyncData } from './useAsyncData';

export function useProjects() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [saving, setSaving] = useState(false);

  const { data, loading, refetch } = useAsyncData(
    async () => {
      if (!user) return { projects: [] as Project[], repos: [] as GitHubRepository[] };
      const [projectsData, reposResponse] = await Promise.all([
        getProjectsByUserId(user.id),
        user.github_connected ? getRepos(user.id) : Promise.resolve({ data: [] as GitHubRepository[] }),
      ]);
      return {
        projects: projectsData,
        repos: reposResponse.data,
      };
    },
    { deps: [user?.id], enabled: !!user }
  );

  const projects = data?.projects ?? [];
  const repos = data?.repos ?? [];

  const [localProjects, setLocalProjects] = useState<Project[] | null>(null);
  const currentProjects = localProjects ?? projects;

  const handleCreate = useCallback(async (reqData: CreateProjectRequest) => {
    setSaving(true);
    try {
      const newProject = await createProject(reqData);
      setLocalProjects(prev => [newProject, ...(prev ?? projects)]);
      toast.success(t('projects.createSuccess'));
      return newProject;
    } catch {
      toast.error(t('projects.createFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, projects]);

  const handleUpdate = useCallback(async (projectId: number, reqData: CreateProjectRequest) => {
    setSaving(true);
    try {
      const updated = await updateProject(projectId, reqData);
      setLocalProjects(prev => (prev ?? projects).map(p => p.id === updated.id ? updated : p));
      toast.success(t('projects.updateSuccess'));
      return updated;
    } catch {
      toast.error(t('projects.updateFailed'));
      return null;
    } finally {
      setSaving(false);
    }
  }, [t, projects]);

  const handleDelete = useCallback(async (project: Project) => {
    if (!confirm(t('projects.confirmDelete'))) return false;
    try {
      await deleteProject(project.id);
      setLocalProjects(prev => (prev ?? projects).filter(p => p.id !== project.id));
      toast.success(t('projects.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('projects.deleteFailed'));
      return false;
    }
  }, [t, projects]);

  return {
    projects: currentProjects,
    repos,
    loading,
    saving,
    createProject: handleCreate,
    updateProject: handleUpdate,
    deleteProject: handleDelete,
    refetch,
  };
}
