import { useTranslation } from 'react-i18next';
import { Pencil, Trash2, CircleDot, CheckCircle2, ExternalLink, Calendar, Github, Archive, ArchiveRestore } from 'lucide-react';
import type { Project } from '../../types/project';
import { cardClass, iconButtonClass, deleteIconButtonClass, badgeBaseClass } from '../../constants/styles';
import { parseJsonArray } from '../../utils/json';
import { sanitizeUrl } from '../../utils/url';
import { formatDate } from '../../utils/timeFormat';

interface ProjectCardProps {
  project: Project;
  onEdit?: () => void;
  onDelete?: () => void;
  onArchive?: () => void;
  onUnarchive?: () => void;
  isOwner?: boolean;
}

export default function ProjectCard({ project, onEdit, onDelete, onArchive, onUnarchive, isOwner }: ProjectCardProps) {
  const { t } = useTranslation();

  const techStack = parseJsonArray(project.tech_stack);

  const durationDays = (() => {
    if (!project.start_date) return null;
    const start = new Date(project.start_date).getTime();
    const end = project.end_date ? new Date(project.end_date).getTime() : Date.now();
    return Math.max(1, Math.ceil((end - start) / (1000 * 60 * 60 * 24)));
  })();

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
              {project.is_archived ? (
                <button
                  onClick={onUnarchive}
                  className={iconButtonClass}
                  aria-label={t('projects.unarchive')}
                  title={t('projects.unarchive')}
                >
                  <ArchiveRestore className="w-4 h-4" />
                </button>
              ) : (
                <>
                  <button
                    onClick={onArchive}
                    className={iconButtonClass}
                    aria-label={t('projects.archive')}
                    title={t('projects.archive')}
                  >
                    <Archive className="w-4 h-4" />
                  </button>
                  <button
                    onClick={onEdit}
                    className={iconButtonClass}
                    aria-label={t('common.edit')}
                  >
                    <Pencil className="w-4 h-4" />
                  </button>
                </>
              )}
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
          <div className="flex items-center gap-2 mt-2">
            <p className="text-xs text-gray-500">
              {formatDate(project.start_date)} - {project.end_date ? formatDate(project.end_date) : t('projects.present')}
            </p>
            {durationDays && (
              <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-cyan-400/10 text-cyan-400">
                <Calendar className="w-3 h-3" />
                {project.end_date
                  ? t('projects.durationCompleted', { days: durationDays })
                  : t('projects.durationOngoing', { days: durationDays })}
              </span>
            )}
          </div>
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
              <Github className="w-4 h-4" />
              GitHub
            </a>
          )}
        </div>
      </div>
    </div>
  );
}
