import { useTranslation } from 'react-i18next';
import { List, Pencil, Trash2 } from 'lucide-react';
import { type RoadmapStep } from '../../api/roadmaps';
import EmptyState from '../common/EmptyState';
import { sanitizeUrl } from '../../utils/url';
import { deleteIconButtonClass, linkSmallClass } from '../../constants/styles';
import { formatDate } from '../../utils/timeFormat';

interface RoadmapStepListProps {
  steps: RoadmapStep[] | undefined;
  isOwner: boolean;
  onToggleComplete: (step: RoadmapStep) => void;
  onEdit: (step: RoadmapStep) => void;
  onDelete: (stepId: number) => void;
  onAddStep: () => void;
}

export default function RoadmapStepList({
  steps,
  isOwner,
  onToggleComplete,
  onEdit,
  onDelete,
  onAddStep,
}: RoadmapStepListProps) {
  const { t } = useTranslation();

  if (steps && steps.length === 0) {
    return (
      <EmptyState
        icon={List}
        message={t('roadmaps.noSteps')}
        actionLabel={isOwner ? t('roadmaps.addFirstStep') : undefined}
        onAction={isOwner ? onAddStep : undefined}
      />
    );
  }

  return (
    <div className="space-y-3">
      {steps?.map((step, index) => (
        <div
          key={step.id}
          className={`bg-gray-800 border rounded-md p-4 transition-colors ${
            step.is_completed ? 'border-green-500/30' : 'border-gray-700'
          }`}
        >
          <div className="flex items-start gap-3">
            {/* Checkbox */}
            {isOwner && (
              <button
                onClick={() => onToggleComplete(step)}
                className={`mt-0.5 w-5 h-5 rounded border-2 flex items-center justify-center flex-shrink-0 transition-colors ${
                  step.is_completed
                    ? 'bg-green-500 border-green-500'
                    : 'border-gray-600 hover:border-gray-500'
                }`}
              >
                {step.is_completed && (
                  <svg className="w-3 h-3 text-white" fill="none" stroke="currentColor" strokeWidth="3" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                  </svg>
                )}
              </button>
            )}

            <div className="flex-1 min-w-0">
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-gray-500 font-medium">#{index + 1}</span>
                    <h3 className={`font-medium ${step.is_completed ? 'line-through text-gray-500' : 'text-white'}`}>
                      {step.title}
                    </h3>
                  </div>
                  {step.is_completed && step.completed_at && (
                    <p className="text-xs text-green-400 mt-0.5">
                      {t('roadmaps.completedOn', { date: formatDate(step.completed_at) })}
                    </p>
                  )}
                  {step.description && (
                    <p className="text-sm text-gray-400 mt-1">{step.description}</p>
                  )}
                  {sanitizeUrl(step.resource_url) && (
                    <a
                      href={sanitizeUrl(step.resource_url)!}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={`${linkSmallClass} mt-2 inline-flex items-center gap-1`}
                    >
                      {t('roadmaps.viewResource')}
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                      </svg>
                    </a>
                  )}
                </div>

                {isOwner && (
                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => onEdit(step)}
                      className="p-1.5 text-gray-400 hover:text-blue-400 transition-colors"
                    >
                      <Pencil className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => onDelete(step.id)}
                      className={deleteIconButtonClass}
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
