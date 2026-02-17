import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { FolderOpen } from 'lucide-react';
import type { Project } from '../types/project';
import { useProjects, useConfirm } from '../hooks';
import ProjectCard from '../components/projects/ProjectCard';
import ProjectForm from '../components/projects/ProjectForm';
import EmptyState from '../components/common/EmptyState';
import ConfirmDialog from '../components/common/ConfirmDialog';
import { Modal, PageHeader, PageLoader } from '../components/common';

export default function ProjectsPage() {
  const { t } = useTranslation();
  const {
    projects, repos, loading, saving,
    createProject, updateProject, deleteProject,
  } = useProjects();

  const { confirm, dialogProps } = useConfirm();

  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  const handleDeleteProject = async (project: Project) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('projects.confirmDelete'), variant: 'danger' });
    if (ok) deleteProject(project);
  };

  if (loading) {
    return <PageLoader />;
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <PageHeader
        title={t('projects.title')}
        subtitle={t('projects.subtitle')}
        actionLabel={t('projects.addProject')}
        onAction={() => setShowForm(true)}
      />

      {/* Form Modal */}
      <Modal
        isOpen={showForm || !!editingProject}
        onClose={() => { setShowForm(false); setEditingProject(null); }}
        title={editingProject ? t('projects.editProject') : t('projects.newProject')}
      >
        <ProjectForm
          project={editingProject || undefined}
          repos={repos}
          onSubmit={async (data) => {
            if (editingProject) {
              const result = await updateProject(editingProject.id, data);
              if (result) setEditingProject(null);
            } else {
              const result = await createProject(data);
              if (result) setShowForm(false);
            }
          }}
          onCancel={() => {
            setShowForm(false);
            setEditingProject(null);
          }}
          loading={saving}
        />
      </Modal>

      {/* Projects Grid */}
      {projects.length === 0 ? (
        <EmptyState
          icon={FolderOpen}
          message={t('projects.noProjects')}
          actionLabel={t('projects.addFirstProject')}
          onAction={() => setShowForm(true)}
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map(project => (
            <ProjectCard
              key={project.id}
              project={project}
              isOwner
              onEdit={() => setEditingProject(project)}
              onDelete={() => handleDeleteProject(project)}
            />
          ))}
        </div>
      )}

      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
