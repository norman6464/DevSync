import { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { useRoadmapDetail } from '../hooks';
import { type RoadmapStep } from '../api/roadmaps';
import { PageLoader, Modal } from '../components/common';
import RoadmapStepList from '../components/roadmaps/RoadmapStepList';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';

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

  const handleNavigateBack = useCallback(() => navigate('/roadmaps'), [navigate]);
  const handleShowStepForm = useCallback(() => setShowStepForm(true), []);
  const handleStepTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setStepTitle(e.target.value), []);
  const handleStepDescriptionChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setStepDescription(e.target.value), []);
  const handleStepResourceURLChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setStepResourceURL(e.target.value), []);

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
    return <PageLoader />;
  }

  if (!roadmap) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-12 text-center">
        <p className="text-gray-400">{t('roadmaps.notFound')}</p>
      </div>
    );
  }


  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      {/* Back button */}
      <button
        onClick={handleNavigateBack}
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
            onClick={handleShowStepForm}
            className={`flex items-center gap-2 ${buttonSecondaryClass}`}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            {t('roadmaps.addStep')}
          </button>
        )}
      </div>

      {/* Progress */}
      <div className="bg-gray-800 border border-gray-700 rounded-md p-4 mb-6">
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
      <Modal
        isOpen={showStepForm}
        onClose={resetStepForm}
        title={editingStep ? t('roadmaps.editStep') : t('roadmaps.addStep')}
        maxWidth="max-w-md"
      >
        <form onSubmit={handleStepSubmit} className="space-y-4">
          <div>
            <label htmlFor="roadmap-step-title" className={labelClass}>{t('roadmaps.stepTitle')}</label>
            <input
              id="roadmap-step-title"
              type="text"
              value={stepTitle}
              onChange={handleStepTitleChange}
              placeholder={t('roadmaps.stepTitlePlaceholder')}
              className={inputClass}
              required
            />
          </div>
          <div>
            <label htmlFor="roadmap-step-description" className={labelClass}>{t('roadmaps.stepDescription')}</label>
            <textarea
              id="roadmap-step-description"
              value={stepDescription}
              onChange={handleStepDescriptionChange}
              placeholder={t('roadmaps.stepDescriptionPlaceholder')}
              rows={3}
              className={`${inputClass} resize-none`}
            />
          </div>
          <div>
            <label htmlFor="roadmap-step-resource-url" className={labelClass}>{t('roadmaps.resourceURL')}</label>
            <input
              id="roadmap-step-resource-url"
              type="url"
              value={stepResourceURL}
              onChange={handleStepResourceURLChange}
              placeholder="https://..."
              maxLength={2048}
              className={inputClass}
            />
          </div>
          <div className="flex gap-3 justify-end pt-2">
            <button
              type="button"
              onClick={resetStepForm}
              className={buttonSecondaryClass}
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={saving || !stepTitle.trim()}
              className={`${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
            >
              {saving ? t('common.saving') : editingStep ? t('common.save') : t('roadmaps.create')}
            </button>
          </div>
        </form>
      </Modal>

      {/* Steps Summary */}
      {roadmap.steps && roadmap.steps.length > 0 && (
        <div className="flex items-center justify-between text-sm text-gray-400 mb-2">
          <span>
            {t('roadmaps.stepsProgress', {
              completed: roadmap.completed_step_count,
              total: roadmap.step_count,
            })}
          </span>
          <span className="text-gray-500">{roadmap.progress}%</span>
        </div>
      )}

      {/* Steps List */}
      <RoadmapStepList
        steps={roadmap.steps}
        isOwner={isOwner}
        onToggleComplete={handleToggleComplete}
        onEdit={handleEditStep}
        onDelete={deleteStep}
        onAddStep={handleShowStepForm}
      />
    </div>
  );
}
