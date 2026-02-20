import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Plus, Trash2, Check, ExternalLink } from 'lucide-react';
import type { StudyCircle, StudyCircleMemberProgress } from '../../types/studyCircle';
import type { User } from '../../types/user';
import { sanitizeUrl } from '../../utils/url';

interface CircleRoadmapTabProps {
  circle: StudyCircle;
  progress: StudyCircleMemberProgress[];
  currentUser: User | null;
  isOwner: boolean;
  saving: boolean;
  onCreateStep: (data: { title: string; description: string; resource_url: string; order_index: number }) => Promise<unknown>;
  onDeleteStep: (stepId: number) => void;
  onToggleProgress: (stepId: number, currentlyCompleted: boolean) => void;
}

export default function CircleRoadmapTab({
  circle, progress, currentUser, isOwner, saving,
  onCreateStep, onDeleteStep, onToggleProgress,
}: CircleRoadmapTabProps) {
  const { t } = useTranslation();
  const [stepForm, setStepForm] = useState({ title: '', description: '', resource_url: '' });
  const [showStepForm, setShowStepForm] = useState(false);

  const handleCreateStep = async () => {
    const result = await onCreateStep({
      ...stepForm,
      order_index: circle.steps?.length || 0,
    });
    if (result) {
      setStepForm({ title: '', description: '', resource_url: '' });
      setShowStepForm(false);
    }
  };

  return (
    <div className="space-y-3">
      {circle.steps?.length === 0 && !showStepForm ? (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-8 text-center">
          <p className="text-sm text-gray-500 mb-3">{t('studyCircle.steps.noSteps')}</p>
          {isOwner && (
            <button
              onClick={() => setShowStepForm(true)}
              className="inline-flex items-center gap-1 px-3 py-1.5 bg-purple-600 hover:bg-purple-500 text-white rounded-lg text-xs font-medium transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
              {t('studyCircle.steps.add')}
            </button>
          )}
        </div>
      ) : (
        <>
          {circle.steps?.map((step) => {
            const myProgress = progress.find(
              (p) => p.step_id === step.id && p.user_id === currentUser?.id
            );
            const isCompleted = myProgress?.is_completed || false;
            return (
              <div
                key={step.id}
                className={`bg-gray-900 border rounded-md p-4 transition-colors ${
                  isCompleted ? 'border-green-500/30' : 'border-gray-800'
                }`}
              >
                <div className="flex items-start gap-3">
                  <button
                    onClick={() => onToggleProgress(step.id, isCompleted)}
                    aria-pressed={isCompleted}
                    aria-label={`${step.title} - ${isCompleted ? t('studyCircle.progress.completed') : t('studyCircle.progress.notCompleted')}`}
                    className={`mt-0.5 w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors ${
                      isCompleted
                        ? 'bg-green-500 border-green-500 text-white'
                        : 'border-gray-600 hover:border-purple-500'
                    }`}
                  >
                    {isCompleted && <Check className="w-3 h-3" aria-hidden="true" />}
                  </button>
                  <div className="flex-1 min-w-0">
                    <h4 className={`text-sm font-medium ${isCompleted ? 'text-green-400 line-through' : 'text-white'}`}>
                      {step.title}
                    </h4>
                    {step.description && (
                      <p className="text-xs text-gray-500 mt-1">{step.description}</p>
                    )}
                    {sanitizeUrl(step.resource_url) && (
                      <a
                        href={sanitizeUrl(step.resource_url)!}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300 mt-1"
                      >
                        <ExternalLink className="w-3 h-3" />
                        {t('studyCircle.steps.resourceUrl')}
                      </a>
                    )}
                    <div className="mt-2">
                      <div className="flex items-center gap-2">
                        {circle.members?.map((m) => {
                          const memberDone = progress.find(
                            (p) => p.step_id === step.id && p.user_id === m.user_id
                          )?.is_completed;
                          return (
                            <div
                              key={m.id}
                              className={`w-5 h-5 rounded-full flex items-center justify-center text-[8px] ${
                                memberDone ? 'bg-green-500/20 text-green-400' : 'bg-gray-800 text-gray-600'
                              }`}
                              title={m.user?.name}
                            >
                              {memberDone ? <Check className="w-3 h-3" /> : (m.user?.name?.[0] || '?')}
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  </div>
                  {isOwner && (
                    <button
                      onClick={() => onDeleteStep(step.id)}
                      aria-label={t('common.delete')}
                      className="text-gray-600 hover:text-red-400 transition-colors"
                    >
                      <Trash2 className="w-3.5 h-3.5" aria-hidden="true" />
                    </button>
                  )}
                </div>
              </div>
            );
          })}
          {isOwner && !showStepForm && (
            <button
              onClick={() => setShowStepForm(true)}
              className="w-full flex items-center justify-center gap-1 py-2 text-xs text-gray-500 hover:text-purple-400 border border-dashed border-gray-800 hover:border-purple-500/50 rounded-md transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
              {t('studyCircle.steps.add')}
            </button>
          )}
        </>
      )}

      {showStepForm && (
        <div className="bg-gray-900 border border-purple-500/30 rounded-md p-4 space-y-3">
          <input
            type="text"
            value={stepForm.title}
            onChange={(e) => setStepForm({ ...stepForm, title: e.target.value })}
            placeholder={t('studyCircle.steps.title')}
            maxLength={100}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
          />
          <textarea
            value={stepForm.description}
            onChange={(e) => setStepForm({ ...stepForm, description: e.target.value })}
            placeholder={t('studyCircle.description')}
            maxLength={500}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 h-16 resize-none"
          />
          <input
            type="url"
            value={stepForm.resource_url}
            onChange={(e) => setStepForm({ ...stepForm, resource_url: e.target.value })}
            maxLength={2048}
            placeholder={t('studyCircle.steps.resourceUrl')}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
          />
          <div className="flex justify-end gap-2">
            <button
              onClick={() => setShowStepForm(false)}
              className="px-3 py-1.5 text-xs text-gray-400 hover:text-white"
            >
              {t('common.cancel')}
            </button>
            <button
              onClick={handleCreateStep}
              disabled={saving || !stepForm.title}
              className="px-3 py-1.5 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-xs font-medium"
            >
              {t('studyCircle.steps.add')}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
