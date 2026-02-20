import { useTranslation } from 'react-i18next';
import { Pencil, Trash2, CircleDot, CheckCircle2, ExternalLink } from 'lucide-react';
import type { Project } from '../../types/project';
import { cardClass, iconButtonClass, deleteIconButtonClass, badgeBaseClass } from '../../constants/styles';
import { parseJsonArray } from '../../utils/json';
import { sanitizeUrl } from '../../utils/url';
import { formatDate } from '../../utils/timeFormat';

interface ProjectCardProps {
  project: Project;
  onEdit?: () => void;
  onDelete?: () => void;
  isOwner?: boolean;
}

export default function ProjectCard({ project, onEdit, onDelete, isOwner }: ProjectCardProps) {
  const { t } = useTranslation();

  const techStack = parseJsonArray(project.tech_stack);

  return (
    <div className={cardClass}>
      {sanitizeUrl(project.image_url) && (
        <div className="aspect-video bg-gray-700 overflow-hidden">
          <img
            src={sanitizeUrl(project.image_url)!}
            alt={project.title}
            referrerPolicy="no-referrer"
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
                <span className={`${badgeBaseClass} bg-yellow-500/20 text-yellow-400`}>
                  {t('projects.featured')}
                </span>
              )}
              {project.end_date ? (
                <span className={`${badgeBaseClass} inline-flex items-center gap-0.5 bg-blue-400/10 text-blue-400`}>
                  <CheckCircle2 className="w-3 h-3" />
                  {t('projects.statusCompleted')}
                </span>
              ) : (
                <span className={`${badgeBaseClass} inline-flex items-center gap-0.5 bg-green-400/10 text-green-400`}>
                  <CircleDot className="w-3 h-3" />
                  {t('projects.statusInProgress')}
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
                className={iconButtonClass}
                aria-label={t('common.edit')}
              >
                <Pencil className="w-4 h-4" />
              </button>
              <button
                onClick={onDelete}
                className={deleteIconButtonClass}
                aria-label={t('common.delete')}
              >
                <Trash2 className="w-4 h-4" />
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
          {sanitizeUrl(project.demo_url) && (
            <a
              href={sanitizeUrl(project.demo_url)!}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 flex items-center justify-center gap-1.5 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-lg transition-colors"
            >
              <ExternalLink className="w-4 h-4" />
              {t('projects.liveDemo')}
            </a>
          )}
          {sanitizeUrl(project.github_url) && (
            <a
              href={sanitizeUrl(project.github_url)!}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 flex items-center justify-center gap-1.5 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-lg transition-colors"
            >
              <svg aria-hidden="true" className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
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
