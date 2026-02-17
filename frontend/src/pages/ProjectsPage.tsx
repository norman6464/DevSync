import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { FolderOpen } from 'lucide-react';
import type { Project } from '../types/project';
import { useProjects } from '../hooks';
import ProjectCard from '../components/projects/ProjectCard';
import ProjectForm from '../components/projects/ProjectForm';
import LoadingSpinner from '../components/common/LoadingSpinner';
import EmptyState from '../components/common/EmptyState';
import { Modal } from '../components/common';

export default function ProjectsPage() {
  const { t } = useTranslation();
  const {
    projects, repos, loading, saving,
    createProject, updateProject, deleteProject,
  } = useProjects();

  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('projects.title')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('projects.subtitle')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('projects.addProject')}
        </button>
      </div>

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
              onDelete={() => deleteProject(project)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
