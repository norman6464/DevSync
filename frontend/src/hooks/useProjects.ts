import { useState, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { Project, CreateProjectRequest } from '../types/project';
import type { GitHubRepository } from '../types/github';
import { getProjectsByUserId, createProject, updateProject, deleteProject, archiveProject, unarchiveProject, getArchivedProjects } from '../api/projects';
import { getRepos } from '../api/github';
import { useAsyncData } from './useAsyncData';
import { useLocalList } from './useCRUDList';

export type ProjectSortBy = 'newest' | 'oldest' | 'title';
export type ProjectTab = 'active' | 'archived';

export function useProjects() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<ProjectSortBy>('newest');
  const [tab, setTab] = useState<ProjectTab>('active');

  const { data, loading, refetch } = useAsyncData(
    async () => {
      if (!user) return { projects: [] as Project[], archivedProjects: [] as Project[], repos: [] as GitHubRepository[] };
      const [projectsData, archivedData, reposResponse] = await Promise.all([
        getProjectsByUserId(user.id),
        getArchivedProjects().then(r => r.projects ?? []),
        user.github_connected ? getRepos(user.id) : Promise.resolve({ data: [] as GitHubRepository[] }),
      ]);
      return {
        projects: projectsData,
        archivedProjects: archivedData,
        repos: reposResponse.data,
      };
    },
    { deps: [user?.id], enabled: !!user }
  );

  const projects = data?.projects ?? [];
  const archivedProjects = data?.archivedProjects ?? [];
  const repos = data?.repos ?? [];

  const { items: currentProjects, setItems: setLocalProjects } = useLocalList(projects);
  const { items: currentArchived, setItems: setLocalArchived } = useLocalList(archivedProjects);

  const filteredProjects = useMemo(() => {
    const source = tab === 'active' ? currentProjects : currentArchived;
    let filtered = source;
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      filtered = filtered.filter(p =>
        p.title.toLowerCase().includes(q) ||
        p.description?.toLowerCase().includes(q) ||
        p.tech_stack?.toLowerCase().includes(q)
      );
    }
    return [...filtered].sort((a, b) => {
      switch (sortBy) {
        case 'oldest':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        case 'title':
          return a.title.localeCompare(b.title);
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
      }
    });
  }, [currentProjects, currentArchived, tab, searchQuery, sortBy]);

  const handleCreate = useCallback(async (reqData: CreateProjectRequest) => {
    try {
      const newProject = await createProject(reqData);
      setLocalProjects(prev => [newProject, ...prev]);
      toast.success(t('projects.createSuccess'));
      return newProject;
    } catch {
      toast.error(t('projects.createFailed'));
      return null;
    }
  }, [t, setLocalProjects]);

  const handleUpdate = useCallback(async (projectId: number, reqData: CreateProjectRequest) => {
    try {
      const updated = await updateProject(projectId, reqData);
      setLocalProjects(prev => prev.map(p => p.id === updated.id ? updated : p));
      toast.success(t('projects.updateSuccess'));
      return updated;
    } catch {
      toast.error(t('projects.updateFailed'));
      return null;
    }
  }, [t, setLocalProjects]);

  const handleDelete = useCallback(async (project: Project) => {
    try {
      await deleteProject(project.id);
      setLocalProjects(prev => prev.filter(p => p.id !== project.id));
      toast.success(t('projects.deleteSuccess'));
      return true;
    } catch {
      toast.error(t('projects.deleteFailed'));
      return false;
    }
  }, [t, setLocalProjects]);

  const handleArchive = useCallback(async (project: Project) => {
    try {
      await archiveProject(project.id);
      setLocalProjects(prev => prev.filter(p => p.id !== project.id));
      setLocalArchived(prev => [{ ...project, is_archived: true }, ...prev]);
      toast.success(t('projects.archiveSuccess'));
      return true;
    } catch {
      toast.error(t('projects.archiveFailed'));
      return false;
    }
  }, [t, setLocalProjects, setLocalArchived]);

  const handleUnarchive = useCallback(async (project: Project) => {
    try {
      await unarchiveProject(project.id);
      setLocalArchived(prev => prev.filter(p => p.id !== project.id));
      setLocalProjects(prev => [{ ...project, is_archived: false }, ...prev]);
      toast.success(t('projects.unarchiveSuccess'));
      return true;
    } catch {
      toast.error(t('projects.unarchiveFailed'));
      return false;
    }
  }, [t, setLocalProjects, setLocalArchived]);

  return {
    projects: filteredProjects,
    allProjects: currentProjects,
    archivedProjects: currentArchived,
    repos,
    loading,
    saving: false,
    searchQuery, setSearchQuery,
    sortBy, setSortBy,
    tab, setTab,
    createProject: handleCreate,
    updateProject: handleUpdate,
    deleteProject: handleDelete,
    archiveProject: handleArchive,
    unarchiveProject: handleUnarchive,
    refetch,
  };
}
