import { useTranslation } from 'react-i18next';
import { List } from 'lucide-react';
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
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                      </svg>
                    </button>
                    <button
                      onClick={() => onDelete(step.id)}
                      className={deleteIconButtonClass}
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 00-7.5 0" />
                      </svg>
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
