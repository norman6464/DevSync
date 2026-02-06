import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Project, CreateProjectRequest } from '../../types/project';
import type { GitHubRepository } from '../../types/github';

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
  const [techStackInput, setTechStackInput] = useState('');
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

  const addTech = () => {
    if (techStackInput.trim() && !techStack.includes(techStackInput.trim())) {
      setTechStack([...techStack, techStackInput.trim()]);
      setTechStackInput('');
    }
  };

  const removeTech = (tech: string) => {
    setTechStack(techStack.filter(t => t !== tech));
  };

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
          <label className="block text-sm font-medium text-gray-300 mb-1">
            {t('projects.linkGitHubRepo')}
          </label>
          <select
            value={githubRepoId || ''}
            onChange={(e) => e.target.value && handleRepoSelect(parseInt(e.target.value))}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
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
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('projects.title')} *
        </label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={200}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder={t('projects.titlePlaceholder')}
        />
      </div>

      {/* Description */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('projects.description')}
        </label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent resize-none"
          placeholder={t('projects.descriptionPlaceholder')}
        />
      </div>

      {/* Tech Stack */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('projects.techStack')}
        </label>
        <div className="flex gap-2">
          <input
            type="text"
            value={techStackInput}
            onChange={(e) => setTechStackInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTech())}
            className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
            placeholder={t('projects.techStackPlaceholder')}
          />
          <button
            type="button"
            onClick={addTech}
            className="px-4 py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg transition-colors"
          >
            {t('common.add')}
          </button>
        </div>
        {techStack.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {techStack.map((tech, idx) => (
              <span
                key={idx}
                className="inline-flex items-center gap-1 px-2 py-1 bg-gray-700 text-gray-300 text-sm rounded"
              >
                {tech}
                <button
                  type="button"
                  onClick={() => removeTech(tech)}
                  className="text-gray-400 hover:text-white"
                >
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            ))}
          </div>
        )}
      </div>

      {/* URLs */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            {t('projects.demoUrl')}
          </label>
          <input
            type="url"
            value={demoUrl}
            onChange={(e) => setDemoUrl(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
            placeholder="https://..."
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            {t('projects.githubUrl')}
          </label>
          <input
            type="url"
            value={githubUrl}
            onChange={(e) => setGithubUrl(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
            placeholder="https://github.com/..."
          />
        </div>
      </div>

      {/* Image URL */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('projects.imageUrl')}
        </label>
        <input
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder="https://..."
        />
      </div>

      {/* Role */}
      <div>
        <label className="block text-sm font-medium text-gray-300 mb-1">
          {t('projects.role')}
        </label>
        <input
          type="text"
          value={role}
          onChange={(e) => setRole(e.target.value)}
          maxLength={100}
          className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-green-500 focus:border-transparent"
          placeholder={t('projects.rolePlaceholder')}
        />
      </div>

      {/* Dates */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            {t('projects.startDate')}
          </label>
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">
            {t('projects.endDate')}
          </label>
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
            className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-green-500 focus:border-transparent"
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
          className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-green-500 focus:ring-green-500"
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
          className="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className="flex-1 px-4 py-2 bg-green-600 hover:bg-green-500 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
