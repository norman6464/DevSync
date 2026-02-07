import { useTranslation } from 'react-i18next';
import type { Project } from '../../types/project';

interface ProjectCardProps {
  project: Project;
  onEdit?: () => void;
  onDelete?: () => void;
  isOwner?: boolean;
}

export default function ProjectCard({ project, onEdit, onDelete, isOwner }: ProjectCardProps) {
  const { t } = useTranslation();

  const techStack = project.tech_stack ? JSON.parse(project.tech_stack) : [];

  const formatDate = (dateStr: string | null) => {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleDateString();
  };

  return (
    <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700 hover:border-gray-600 transition-colors">
      {project.image_url && (
        <div className="aspect-video bg-gray-700 overflow-hidden">
          <img
            src={project.image_url}
            alt={project.title}
            className="w-full h-full object-cover"
          />
        </div>
      )}

      <div className="p-4">
        <div className="flex items-start justify-between gap-2">
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <h3 className="text-lg font-semibold text-white">{project.title}</h3>
              {project.featured && (
                <span className="px-2 py-0.5 bg-yellow-500/20 text-yellow-400 text-xs rounded-full">
                  {t('projects.featured')}
                </span>
              )}
            </div>
            {project.role && (
              <p className="text-sm text-gray-400 mt-0.5">{project.role}</p>
            )}
          </div>

          {isOwner && (
            <div className="flex gap-1">
              <button
                onClick={onEdit}
                className="p-1.5 text-gray-400 hover:text-white transition-colors"
                title={t('common.edit')}
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                </svg>
              </button>
              <button
                onClick={onDelete}
                className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
                title={t('common.delete')}
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                </svg>
              </button>
            </div>
          )}
        </div>

        {project.description && (
          <p className="text-gray-300 text-sm mt-2 line-clamp-2">{project.description}</p>
        )}

        {techStack.length > 0 && (
          <div className="flex flex-wrap gap-1.5 mt-3">
            {techStack.slice(0, 5).map((tech: string, idx: number) => (
              <span
                key={idx}
                className="px-2 py-0.5 bg-gray-700 text-gray-300 text-xs rounded"
              >
                {tech}
              </span>
            ))}
            {techStack.length > 5 && (
              <span className="px-2 py-0.5 text-gray-400 text-xs">
                +{techStack.length - 5}
              </span>
            )}
          </div>
        )}

        {(project.start_date || project.end_date) && (
          <p className="text-xs text-gray-500 mt-2">
            {formatDate(project.start_date)} - {project.end_date ? formatDate(project.end_date) : t('projects.present')}
          </p>
        )}

        <div className="flex gap-2 mt-4">
          {project.demo_url && (
            <a
              href={project.demo_url}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 flex items-center justify-center gap-1.5 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-lg transition-colors"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
              </svg>
              {t('projects.liveDemo')}
            </a>
          )}
          {project.github_url && (
            <a
              href={project.github_url}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 flex items-center justify-center gap-1.5 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-lg transition-colors"
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
              </svg>
              GitHub
            </a>
          )}
        </div>
      </div>
    </div>
  );
}
