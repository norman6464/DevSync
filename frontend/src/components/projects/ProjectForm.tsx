import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Project, CreateProjectRequest } from '../../types/project';
import type { GitHubRepository } from '../../types/github';
import { buttonSecondaryClass, inputClass, labelClass, textareaClass } from '../../constants/styles';
import { findInvalidUrlField } from '../../utils/url';
import TagInput from '../common/TagInput';
import toast from 'react-hot-toast';

interface ProjectFormProps {
  project?: Project;
  repos?: GitHubRepository[];
  onSubmit: (data: CreateProjectRequest) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

export default function ProjectForm({ project, repos = [], onSubmit, onCancel, loading }: ProjectFormProps) {
  const { t } = useTranslation();
  const [title, setTitle] = useState(project?.title || '');
  const [description, setDescription] = useState(project?.description || '');
  const [techStack, setTechStack] = useState<string[]>(() => {
    if (project?.tech_stack) {
      try {
        return JSON.parse(project.tech_stack);
      } catch {
        return [];
      }
    }
    return [];
  });
  const [demoUrl, setDemoUrl] = useState(project?.demo_url || '');
  const [githubUrl, setGithubUrl] = useState(project?.github_url || '');
  const [imageUrl, setImageUrl] = useState(project?.image_url || '');
  const [role, setRole] = useState(project?.role || '');
  const [startDate, setStartDate] = useState(project?.start_date?.split('T')[0] || '');
  const [endDate, setEndDate] = useState(project?.end_date?.split('T')[0] || '');
  const [featured, setFeatured] = useState(project?.featured || false);
  const [githubRepoId, setGithubRepoId] = useState<number | undefined>(project?.github_repo_id || undefined);

  const handleRepoSelect = (repoId: number) => {
    const repo = repos.find(r => r.id === repoId);
    if (repo) {
      setGithubRepoId(repoId);
      setGithubUrl(`https://github.com/${repo.full_name}`);
      if (!title) setTitle(repo.name);
      if (!description && repo.description) setDescription(repo.description);
      if (repo.language && !techStack.includes(repo.language)) {
        setTechStack([...techStack, repo.language]);
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const invalidField = findInvalidUrlField([
      { value: demoUrl, label: t('projects.demoUrl') },
      { value: githubUrl, label: t('projects.githubUrl') },
      { value: imageUrl, label: t('projects.imageUrl') },
    ]);
    if (invalidField) {
      toast.error(t('common.invalidUrl', { field: invalidField }));
      return;
    }
    await onSubmit({
      title,
      description,
      tech_stack: JSON.stringify(techStack),
      demo_url: demoUrl,
      github_url: githubUrl,
      image_url: imageUrl,
      role,
      start_date: startDate || undefined,
      end_date: endDate || undefined,
      featured,
      github_repo_id: githubRepoId,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Link to GitHub Repo */}
      {repos.length > 0 && (
        <div>
          <label className={labelClass}>
            {t('projects.linkGitHubRepo')}
          </label>
          <select
            value={githubRepoId || ''}
            onChange={(e) => e.target.value && handleRepoSelect(parseInt(e.target.value))}
            className={inputClass}
          >
            <option value="">{t('projects.selectRepo')}</option>
            {repos.map(repo => (
              <option key={repo.id} value={repo.id}>
                {repo.name} {repo.language && `(${repo.language})`}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Title */}
      <div>
        <label className={labelClass}>
          {t('projects.title')} *
        </label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={200}
          className={inputClass}
          placeholder={t('projects.titlePlaceholder')}
        />
      </div>

      {/* Description */}
      <div>
        <label className={labelClass}>
          {t('projects.description')}
        </label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          maxLength={2000}
          className={textareaClass}
          placeholder={t('projects.descriptionPlaceholder')}
        />
      </div>

      {/* Tech Stack */}
      <TagInput
        tags={techStack}
        onChange={setTechStack}
        label={t('projects.techStack')}
        placeholder={t('projects.techStackPlaceholder')}
        maxLength={100}
      />

      {/* URLs */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className={labelClass}>
            {t('projects.demoUrl')}
          </label>
          <input
            type="url"
            value={demoUrl}
            onChange={(e) => setDemoUrl(e.target.value)}
            maxLength={2048}
            className={inputClass}
            placeholder="https://..."
          />
        </div>
        <div>
          <label className={labelClass}>
            {t('projects.githubUrl')}
          </label>
          <input
            type="url"
            value={githubUrl}
            onChange={(e) => setGithubUrl(e.target.value)}
            maxLength={2048}
            className={inputClass}
            placeholder="https://github.com/..."
          />
        </div>
      </div>

      {/* Image URL */}
      <div>
        <label className={labelClass}>
          {t('projects.imageUrl')}
        </label>
        <input
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          maxLength={2048}
          className={inputClass}
          placeholder="https://..."
        />
      </div>

      {/* Role */}
      <div>
        <label className={labelClass}>
          {t('projects.role')}
        </label>
        <input
          type="text"
          value={role}
          onChange={(e) => setRole(e.target.value)}
          maxLength={100}
          className={inputClass}
          placeholder={t('projects.rolePlaceholder')}
        />
      </div>

      {/* Dates */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className={labelClass}>
            {t('projects.startDate')}
          </label>
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className={inputClass}
          />
        </div>
        <div>
          <label className={labelClass}>
            {t('projects.endDate')}
          </label>
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className={inputClass}
          />
        </div>
      </div>

      {/* Featured */}
      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="featured"
          checked={featured}
          onChange={(e) => setFeatured(e.target.checked)}
          className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-green-500 focus:ring-gray-500"
        />
        <label htmlFor="featured" className="text-sm text-gray-300">
          {t('projects.markFeatured')}
        </label>
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className={`flex-1 ${buttonSecondaryClass}`}
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className={`flex-1 ${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
