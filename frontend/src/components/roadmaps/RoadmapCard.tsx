import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, Pencil, Trash2, type LucideIcon } from 'lucide-react';
import { type Roadmap, type RoadmapCategory } from '../../api/roadmaps';
import { badgeBaseClass, editIconButtonClass, deleteIconButtonLargeClass } from '../../constants/styles';

const CATEGORIES: { value: RoadmapCategory; label: string; icon: string; Icon: LucideIcon }[] = [
  { value: 'language', label: 'roadmaps.categoryLanguage', icon: '💻', Icon: Monitor },
  { value: 'framework', label: 'roadmaps.categoryFramework', icon: '🚀', Icon: Rocket },
  { value: 'skill', label: 'roadmaps.categorySkill', icon: '🎯', Icon: Target },
  { value: 'project', label: 'roadmaps.categoryProject', icon: '📁', Icon: FolderOpen },
  { value: 'other', label: 'roadmaps.categoryOther', icon: '📝', Icon: FileText },
];

const getCategoryInfo = (cat: RoadmapCategory) =>
  CATEGORIES.find(c => c.value === cat) || CATEGORIES[4];

interface RoadmapCardProps {
  roadmap: Roadmap;
  onView: () => void;
  onEdit: (r: Roadmap) => void;
  onDelete: (id: number) => void;
}

export default function RoadmapCard({ roadmap, onView, onEdit, onDelete }: RoadmapCardProps) {
  const { t } = useTranslation();
  const catInfo = getCategoryInfo(roadmap.category);
  const CategoryIcon = catInfo.Icon;

  return (
    <div
      onClick={onView}
      className="bg-gray-800 border border-gray-700 rounded-md p-4 hover:border-gray-600 transition-colors cursor-pointer"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3 min-w-0 flex-1">
          <CategoryIcon className="w-6 h-6 text-purple-400 flex-shrink-0 mt-0.5" />
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <h3 className="font-medium text-white">{roadmap.title}</h3>
              {roadmap.is_public && (
                <span className={`${badgeBaseClass} bg-blue-500/10 text-blue-400`}>
                  {t('roadmaps.public')}
                </span>
              )}
              {roadmap.status === 'completed' && (
                <span className={`${badgeBaseClass} bg-green-500/10 text-green-400`}>
                  {t('roadmaps.completed')}
                </span>
              )}
            </div>
            {roadmap.description && (
              <p className="text-sm text-gray-400 mt-1 line-clamp-1">{roadmap.description}</p>
            )}
            <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
              <span>{t(catInfo.label)}</span>
              <span>{roadmap.completed_step_count} / {roadmap.step_count} {t('roadmaps.stepsCompleted')}</span>
            </div>

            {/* Progress Bar */}
            <div className="mt-3">
              <div className="flex items-center justify-between text-xs mb-1">
                <span className="text-gray-400">{t('roadmaps.progress')}</span>
                <span className="text-gray-300">{roadmap.progress}%</span>
              </div>
              <div className="h-2 bg-gray-700 rounded-full overflow-hidden">
                <div
                  className={`h-full transition-all ${roadmap.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'}`}
                  style={{ width: `${roadmap.progress}%` }}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-1" onClick={e => e.stopPropagation()}>
          <button
            onClick={() => onEdit(roadmap)}
            className={editIconButtonClass}
          >
            <Pencil className="w-4 h-4" />
          </button>
          <button
            onClick={() => onDelete(roadmap.id)}
            className={deleteIconButtonLargeClass}
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
