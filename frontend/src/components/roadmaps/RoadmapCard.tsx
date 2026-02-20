import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, type LucideIcon } from 'lucide-react';
import { type Roadmap, type RoadmapCategory } from '../../api/roadmaps';

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
                <span className="px-2 py-0.5 text-xs rounded-full bg-blue-500/10 text-blue-400">
                  {t('roadmaps.public')}
                </span>
              )}
              {roadmap.status === 'completed' && (
                <span className="px-2 py-0.5 text-xs rounded-full bg-green-500/10 text-green-400">
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
            className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
            </svg>
          </button>
          <button
            onClick={() => onDelete(roadmap.id)}
            className="p-2 text-gray-400 hover:text-red-400 transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
