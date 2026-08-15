import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { FolderOpen, Search, ArrowUpDown, Archive } from 'lucide-react';
import type { Project } from '../types/project';
import { useProjects, useConfirm } from '../hooks';
import type { ProjectSortBy, ProjectTab } from '../hooks/useProjects';
import ProjectCard from '../components/projects/ProjectCard';
import ProjectForm from '../components/projects/ProjectForm';
import EmptyState from '../components/common/EmptyState';
import ConfirmDialog from '../components/common/ConfirmDialog';
import { Modal, PageHeader, PageLoader } from '../components/common';
import { inputClass } from '../constants/styles';

export default function ProjectsPage() {
  const { t } = useTranslation();
  const {
    projects, allProjects, archivedProjects, repos, loading, saving,
    searchQuery, setSearchQuery, sortBy, setSortBy, tab, setTab,
    createProject, updateProject, deleteProject,
    archiveProject, unarchiveProject,
  } = useProjects();

  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  const handleDeleteProject = useCallback(async (project: Project) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('projects.confirmDelete'), variant: 'danger' });
    if (ok) deleteProject(project);
  }, [confirm, t, deleteProject]);

  const handleFormClose = useCallback(() => {
    setShowForm(false);
    setEditingProject(null);
  }, []);

  const handleFormSubmit = useCallback(async (data: Parameters<typeof createProject>[0]) => {
    if (editingProject) {
      const result = await updateProject(editingProject.id, data);
      if (result) setEditingProject(null);
    } else {
      const result = await createProject(data);
      if (result) setShowForm(false);
    }
  }, [editingProject, updateProject, createProject]);

  if (loading) {
    return <PageLoader />;
  }

  const SORT_OPTIONS: { value: ProjectSortBy; label: string }[] = [
    { value: 'newest', label: t('projects.sortNewest') },
    { value: 'oldest', label: t('projects.sortOldest') },
    { value: 'title', label: t('projects.sortTitle') },
  ];

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <PageHeader
        title={t('projects.title')}
        subtitle={t('projects.subtitle')}
        actionLabel={t('projects.addProject')}
        onAction={() => setShowForm(true)}
      />

      {/* Search */}
      <div className="relative mb-4">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder={t('projects.searchPlaceholder')}
          className={`${inputClass} pl-9`}
        />
      </div>

      {/* Tabs & Sort */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex gap-2">
          {(['active', 'archived'] as ProjectTab[]).map((tabValue) => (
            <button
              key={tabValue}
              onClick={() => setTab(tabValue)}
              className={`flex items-center gap-1.5 px-4 py-2 text-sm rounded-lg transition-colors ${
                tab === tabValue
                  ? tabValue === 'archived'
                    ? 'bg-yellow-500/10 text-yellow-400 border border-yellow-500/50'
                    : 'bg-blue-500/10 text-blue-400 border border-blue-500/50'
                  : 'text-gray-400 border border-gray-700 hover:border-gray-600'
              }`}
            >
              {tabValue === 'archived' && <Archive className="w-4 h-4" />}
              {tabValue === 'active'
                ? `${t('projects.title')} (${allProjects.length})`
                : `${t('projects.archivedTab')} (${archivedProjects.length})`}
            </button>
          ))}
        </div>
        <div className="flex items-center gap-2">
          <ArrowUpDown className="w-4 h-4 text-gray-500" />
          {SORT_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              onClick={() => setSortBy(opt.value)}
              className={`px-3 py-1 text-xs rounded-full border transition-colors ${
                sortBy === opt.value
                  ? 'border-blue-500 bg-blue-500/10 text-blue-400'
                  : 'border-gray-700 text-gray-400 hover:border-gray-600'
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingProject}
        onClose={handleFormClose}
        title={editingProject ? t('projects.editProject') : t('projects.newProject')}
      >
        <ProjectForm
          project={editingProject || undefined}
          repos={repos}
          onSubmit={handleFormSubmit}
          onCancel={handleFormClose}
          loading={saving}
        />
      </Modal>

      {/* Projects Grid */}
      {projects.length === 0 ? (
        searchQuery ? (
          <EmptyState
            icon={FolderOpen}
            title={t('projects.noSearchResults')}
          />
        ) : tab === 'archived' ? (
          <EmptyState
            icon={Archive}
            title={t('projects.noSearchResults')}
          />
        ) : (
          <EmptyState
            icon={FolderOpen}
            title={t('projects.noProjects')}
            actionLabel={t('projects.addFirstProject')}
            onAction={() => setShowForm(true)}
          />
        )
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map(project => (
            <ProjectCard
              key={project.id}
              project={project}
              isOwner
              onEdit={() => setEditingProject(project)}
              onDelete={() => handleDeleteProject(project)}
              onArchive={() => archiveProject(project)}
              onUnarchive={() => unarchiveProject(project)}
            />
          ))}
        </div>
      )}

      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
