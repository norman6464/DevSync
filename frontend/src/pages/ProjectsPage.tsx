import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { useAuthStore } from '../store/authStore';
import type { Project, CreateProjectRequest } from '../types/project';
import type { GitHubRepository } from '../types/github';
import { getProjectsByUserId, createProject, updateProject, deleteProject } from '../api/projects';
import { getRepos } from '../api/github';
import ProjectCard from '../components/projects/ProjectCard';
import ProjectForm from '../components/projects/ProjectForm';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function ProjectsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [projects, setProjects] = useState<Project[]>([]);
  const [repos, setRepos] = useState<GitHubRepository[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  useEffect(() => {
    if (user) {
      fetchData();
    }
  }, [user]);

  const fetchData = async () => {
    if (!user) return;
    setLoading(true);
    try {
      const [projectsData, reposResponse] = await Promise.all([
        getProjectsByUserId(user.id),
        user.github_connected ? getRepos(user.id) : Promise.resolve({ data: [] }),
      ]);
      setProjects(projectsData);
      setRepos(reposResponse.data);
    } catch (error) {
      console.error('Failed to fetch data:', error);
      toast.error(t('errors.somethingWrong'));
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (data: CreateProjectRequest) => {
    setSaving(true);
    try {
      const newProject = await createProject(data);
      setProjects([newProject, ...projects]);
      setShowForm(false);
      toast.success(t('projects.createSuccess'));
    } catch (error) {
      console.error('Failed to create project:', error);
      toast.error(t('projects.createFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleUpdate = async (data: CreateProjectRequest) => {
    if (!editingProject) return;
    setSaving(true);
    try {
      const updated = await updateProject(editingProject.id, data);
      setProjects(projects.map(p => p.id === updated.id ? updated : p));
      setEditingProject(null);
      toast.success(t('projects.updateSuccess'));
    } catch (error) {
      console.error('Failed to update project:', error);
      toast.error(t('projects.updateFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (project: Project) => {
    if (!confirm(t('projects.confirmDelete'))) return;
    try {
      await deleteProject(project.id);
      setProjects(projects.filter(p => p.id !== project.id));
      toast.success(t('projects.deleteSuccess'));
    } catch (error) {
      console.error('Failed to delete project:', error);
      toast.error(t('projects.deleteFailed'));
    }
  };

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
          className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('projects.addProject')}
        </button>
      </div>

      {/* Form Modal */}
      {(showForm || editingProject) && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingProject ? t('projects.editProject') : t('projects.newProject')}
            </h2>
            <ProjectForm
              project={editingProject || undefined}
              repos={repos}
              onSubmit={editingProject ? handleUpdate : handleCreate}
              onCancel={() => {
                setShowForm(false);
                setEditingProject(null);
              }}
              loading={saving}
            />
          </div>
        </div>
      )}

      {/* Projects Grid */}
      {projects.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" />
            </svg>
          </div>
          <p className="text-gray-400">{t('projects.noProjects')}</p>
          <button
            onClick={() => setShowForm(true)}
            className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
          >
            {t('projects.addFirstProject')}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map(project => (
            <ProjectCard
              key={project.id}
              project={project}
              isOwner
              onEdit={() => setEditingProject(project)}
              onDelete={() => handleDelete(project)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
