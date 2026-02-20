import { useTranslation } from 'react-i18next';
import { MapPin } from 'lucide-react';
import { type RoadmapCategory } from '../api/roadmaps';
import { useRoadmapForm } from '../hooks';
import EmptyState from '../components/common/EmptyState';
import { Modal, PageHeader, PageLoader } from '../components/common';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';
import RoadmapTemplatesSection from '../components/roadmaps/RoadmapTemplatesSection';
import RoadmapCard from '../components/roadmaps/RoadmapCard';

const CATEGORIES: { value: RoadmapCategory; label: string; icon: string }[] = [
  { value: 'language', label: 'roadmaps.categoryLanguage', icon: '💻' },
  { value: 'framework', label: 'roadmaps.categoryFramework', icon: '🚀' },
  { value: 'skill', label: 'roadmaps.categorySkill', icon: '🎯' },
  { value: 'project', label: 'roadmaps.categoryProject', icon: '📁' },
  { value: 'other', label: 'roadmaps.categoryOther', icon: '📝' },
];

export default function RoadmapsPage() {
  const { t } = useTranslation();
  const {
    roadmaps, activeRoadmaps, completedRoadmaps, templates,
    loading, saving, templatesLoading, creating,
    showForm, setShowForm, editingRoadmap,
    title, setTitle, description, setDescription,
    category, setCategory, isPublic, setIsPublic,
    showTemplates, expandedTemplate,
    resetForm, handleSubmit, handleEdit, handleUseTemplate,
    deleteRoadmap, toggleTemplates, toggleExpandedTemplate, navigate,
  } = useRoadmapForm();

  if (loading) {
    return <PageLoader />;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <PageHeader
        title={t('roadmaps.title')}
        subtitle={t('roadmaps.description')}
        actionLabel={t('roadmaps.addRoadmap')}
        onAction={() => setShowForm(true)}
      />

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gray-800 rounded-md p-4 border border-gray-700">
          <p className="text-2xl font-bold text-white">{roadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.totalRoadmaps')}</p>
        </div>
        <div className="bg-gray-800 rounded-md p-4 border border-gray-700">
          <p className="text-2xl font-bold text-blue-400">{activeRoadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.activeRoadmaps')}</p>
        </div>
        <div className="bg-gray-800 rounded-md p-4 border border-gray-700">
          <p className="text-2xl font-bold text-green-400">{completedRoadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.completedRoadmaps')}</p>
        </div>
      </div>

      {/* Templates Section */}
      {!templatesLoading && (
        <RoadmapTemplatesSection
          templates={templates}
          showTemplates={showTemplates}
          expandedTemplate={expandedTemplate}
          creating={creating}
          toggleTemplates={toggleTemplates}
          toggleExpandedTemplate={toggleExpandedTemplate}
          handleUseTemplate={handleUseTemplate}
        />
      )}

      {/* Form Modal */}
      <Modal
        isOpen={showForm}
        onClose={resetForm}
        title={editingRoadmap ? t('roadmaps.editRoadmap') : t('roadmaps.addRoadmap')}
        maxWidth="max-w-md"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className={labelClass}>{t('roadmaps.roadmapTitle')}</label>
            <input
              type="text"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder={t('roadmaps.titlePlaceholder')}
              className={inputClass}
              required
            />
          </div>
          <div>
            <label className={labelClass}>{t('roadmaps.descriptionLabel')}</label>
            <textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              placeholder={t('roadmaps.descriptionPlaceholder')}
              rows={3}
              className={`${inputClass} resize-none`}
            />
          </div>
          <div>
            <label className={labelClass}>{t('roadmaps.category')}</label>
            <select
              value={category}
              onChange={e => setCategory(e.target.value as RoadmapCategory)}
              className={inputClass}
            >
              {CATEGORIES.map(cat => (
                <option key={cat.value} value={cat.value}>
                  {cat.icon} {t(cat.label)}
                </option>
              ))}
            </select>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="is_public"
              checked={isPublic}
              onChange={e => setIsPublic(e.target.checked)}
              className="w-4 h-4 rounded border-gray-600 bg-gray-700"
            />
            <label htmlFor="is_public" className="text-sm text-gray-300">{t('roadmaps.makePublic')}</label>
          </div>
          <div className="flex gap-3 justify-end pt-2">
            <button
              type="button"
              onClick={resetForm}
              className={buttonSecondaryClass}
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={saving || !title.trim()}
              className={`${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
            >
              {saving ? t('common.saving') : editingRoadmap ? t('common.save') : t('roadmaps.create')}
            </button>
          </div>
        </form>
      </Modal>

      {/* Content */}
      {roadmaps.length === 0 ? (
        <EmptyState
          icon={MapPin}
          message={t('roadmaps.noRoadmaps')}
          actionLabel={t('roadmaps.createFirst')}
          onAction={() => setShowForm(true)}
        />
      ) : (
        <div className="space-y-6">
          {activeRoadmaps.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-3">
                {t('roadmaps.activeRoadmaps')} ({activeRoadmaps.length})
              </h2>
              <div className="space-y-3">
                {activeRoadmaps.map(roadmap => (
                  <RoadmapCard
                    key={roadmap.id}
                    roadmap={roadmap}
                    onView={() => navigate(`/roadmaps/${roadmap.id}`)}
                    onEdit={handleEdit}
                    onDelete={deleteRoadmap}
                  />
                ))}
              </div>
            </div>
          )}

          {completedRoadmaps.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-3">
                {t('roadmaps.completedRoadmaps')} ({completedRoadmaps.length})
              </h2>
              <div className="space-y-3">
                {completedRoadmaps.map(roadmap => (
                  <RoadmapCard
                    key={roadmap.id}
                    roadmap={roadmap}
                    onView={() => navigate(`/roadmaps/${roadmap.id}`)}
                    onEdit={handleEdit}
                    onDelete={deleteRoadmap}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

