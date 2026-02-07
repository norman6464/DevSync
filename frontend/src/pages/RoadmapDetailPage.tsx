import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { useRoadmapDetail } from '../hooks';
import { type RoadmapStep } from '../api/roadmaps';
import LoadingSpinner from '../components/common/LoadingSpinner';

export default function RoadmapDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const user = useAuthStore(s => s.user);
  const roadmapId = id ? parseInt(id) : null;

  const {
    roadmap, loading, saving,
    createStep, updateStep, deleteStep,
  } = useRoadmapDetail(roadmapId);

  const [showStepForm, setShowStepForm] = useState(false);
  const [editingStep, setEditingStep] = useState<RoadmapStep | null>(null);
  const [stepTitle, setStepTitle] = useState('');
  const [stepDescription, setStepDescription] = useState('');
  const [stepResourceURL, setStepResourceURL] = useState('');

  const isOwner = user?.id === roadmap?.user_id;

  const resetStepForm = () => {
    setStepTitle('');
    setStepDescription('');
    setStepResourceURL('');
    setEditingStep(null);
    setShowStepForm(false);
  };

  const handleStepSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!stepTitle.trim()) return;

    if (editingStep) {
      const result = await updateStep(editingStep.id, {
        title: stepTitle,
        description: stepDescription,
        resource_url: stepResourceURL,
      });
      if (result) resetStepForm();
    } else {
      const result = await createStep({
        title: stepTitle,
        description: stepDescription,
        resource_url: stepResourceURL,
      });
      if (result) resetStepForm();
    }
  };

  const handleEditStep = (step: RoadmapStep) => {
    setEditingStep(step);
    setStepTitle(step.title);
    setStepDescription(step.description);
    setStepResourceURL(step.resource_url);
    setShowStepForm(true);
  };

  const handleToggleComplete = async (step: RoadmapStep) => {
    await updateStep(step.id, { is_completed: !step.is_completed });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <LoadingSpinner />
      </div>
    );
  }

  if (!roadmap) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-12 text-center">
        <p className="text-gray-400">{t('roadmaps.notFound')}</p>
      </div>
    );
  }

  const inputClass = 'w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent';

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      {/* Back button */}
      <button
        onClick={() => navigate('/roadmaps')}
        className="text-sm text-gray-400 hover:text-white transition-colors mb-4 flex items-center gap-1"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
        </svg>
        {t('roadmaps.backToList')}
      </button>

      {/* Header */}
      <div className="flex items-start justify-between gap-4 mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white mb-2">{roadmap.title}</h1>
          {roadmap.description && (
            <p className="text-gray-400">{roadmap.description}</p>
          )}
          <div className="flex items-center gap-3 mt-3 text-sm">
            <span className="px-2 py-1 bg-gray-800 rounded text-gray-300">
              {t(`roadmaps.category${roadmap.category.charAt(0).toUpperCase() + roadmap.category.slice(1)}`)}
            </span>
            {roadmap.is_public && (
              <span className="px-2 py-1 bg-blue-500/10 text-blue-400 rounded text-xs">
                {t('roadmaps.public')}
              </span>
            )}
            {roadmap.status === 'completed' && (
              <span className="px-2 py-1 bg-green-500/10 text-green-400 rounded text-xs">
                {t('roadmaps.completed')}
              </span>
            )}
          </div>
        </div>
        {isOwner && (
          <button
            onClick={() => setShowStepForm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            {t('roadmaps.addStep')}
          </button>
        )}
      </div>

      {/* Progress */}
      <div className="bg-gray-800 border border-gray-700 rounded-xl p-4 mb-6">
        <div className="flex items-center justify-between mb-2">
          <span className="text-sm text-gray-400">{t('roadmaps.progress')}</span>
          <span className="text-sm font-medium text-white">
            {roadmap.completed_step_count} / {roadmap.step_count} {t('roadmaps.stepsCompleted')}
          </span>
        </div>
        <div className="h-3 bg-gray-700 rounded-full overflow-hidden">
          <div
            className={`h-full transition-all ${roadmap.status === 'completed' ? 'bg-green-500' : 'bg-blue-500'}`}
            style={{ width: `${roadmap.progress}%` }}
          />
        </div>
        <p className="text-xs text-gray-500 mt-2">{roadmap.progress}% {t('roadmaps.complete')}</p>
      </div>

      {/* Step Form Modal */}
      {showStepForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-md">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingStep ? t('roadmaps.editStep') : t('roadmaps.addStep')}
            </h2>
            <form onSubmit={handleStepSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.stepTitle')}</label>
                <input
                  type="text"
                  value={stepTitle}
                  onChange={e => setStepTitle(e.target.value)}
                  placeholder={t('roadmaps.stepTitlePlaceholder')}
                  className={inputClass}
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.stepDescription')}</label>
                <textarea
                  value={stepDescription}
                  onChange={e => setStepDescription(e.target.value)}
                  placeholder={t('roadmaps.stepDescriptionPlaceholder')}
                  rows={3}
                  className={`${inputClass} resize-none`}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.resourceURL')}</label>
                <input
                  type="url"
                  value={stepResourceURL}
                  onChange={e => setStepResourceURL(e.target.value)}
                  placeholder="https://..."
                  className={inputClass}
                />
              </div>
              <div className="flex gap-3 justify-end pt-2">
                <button
                  type="button"
                  onClick={resetStepForm}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  {t('common.cancel')}
                </button>
                <button
                  type="submit"
                  disabled={saving || !stepTitle.trim()}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
                >
                  {saving ? t('common.saving') : editingStep ? t('common.save') : t('roadmaps.create')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Steps List */}
      {roadmap.steps && roadmap.steps.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0ZM3.75 12h.007v.008H3.75V12Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm-.375 5.25h.007v.008H3.75v-.008Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
            </svg>
          </div>
          <p className="text-gray-400">{t('roadmaps.noSteps')}</p>
          {isOwner && (
            <button
              onClick={() => setShowStepForm(true)}
              className="mt-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
            >
              {t('roadmaps.addFirstStep')}
            </button>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {roadmap.steps?.map((step, index) => (
            <div
              key={step.id}
              className={`bg-gray-800 border rounded-xl p-4 transition-colors ${
                step.is_completed ? 'border-green-500/30' : 'border-gray-700'
              }`}
            >
              <div className="flex items-start gap-3">
                {/* Checkbox */}
                {isOwner && (
                  <button
                    onClick={() => handleToggleComplete(step)}
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
                      {step.description && (
                        <p className="text-sm text-gray-400 mt-1">{step.description}</p>
                      )}
                      {step.resource_url && (
                        <a
                          href={step.resource_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-xs text-blue-400 hover:text-blue-300 mt-2 inline-flex items-center gap-1"
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
                          onClick={() => handleEditStep(step)}
                          className="p-1.5 text-gray-400 hover:text-blue-400 transition-colors"
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                          </svg>
                        </button>
                        <button
                          onClick={() => deleteStep(step.id)}
                          className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
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
      )}
    </div>
  );
}
