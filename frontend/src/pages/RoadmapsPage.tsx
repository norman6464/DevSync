import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Monitor, Rocket, Target, FolderOpen, FileText, type LucideIcon } from 'lucide-react';
import { type RoadmapCategory, type Roadmap } from '../api/roadmaps';
import { useRoadmaps } from '../hooks';
import LoadingSpinner from '../components/common/LoadingSpinner';

const CATEGORIES: { value: RoadmapCategory; label: string; icon: string; Icon: LucideIcon }[] = [
  { value: 'language', label: 'roadmaps.categoryLanguage', icon: '💻', Icon: Monitor },
  { value: 'framework', label: 'roadmaps.categoryFramework', icon: '🚀', Icon: Rocket },
  { value: 'skill', label: 'roadmaps.categorySkill', icon: '🎯', Icon: Target },
  { value: 'project', label: 'roadmaps.categoryProject', icon: '📁', Icon: FolderOpen },
  { value: 'other', label: 'roadmaps.categoryOther', icon: '📝', Icon: FileText },
];

export default function RoadmapsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    roadmaps, loading, saving, activeRoadmaps, completedRoadmaps,
    createRoadmap, updateRoadmap, deleteRoadmap,
  } = useRoadmaps();

  const [showForm, setShowForm] = useState(false);
  const [editingRoadmap, setEditingRoadmap] = useState<Roadmap | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState<RoadmapCategory>('other');
  const [isPublic, setIsPublic] = useState(false);

  const resetForm = () => {
    setTitle('');
    setDescription('');
    setCategory('other');
    setIsPublic(false);
    setEditingRoadmap(null);
    setShowForm(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    if (editingRoadmap) {
      const result = await updateRoadmap(editingRoadmap.id, { title, description, category, is_public: isPublic });
      if (result) resetForm();
    } else {
      const result = await createRoadmap({ title, description, category, is_public: isPublic });
      if (result) {
        resetForm();
        navigate(`/roadmaps/${result.id}`);
      }
    }
  };

  const handleEdit = (roadmap: Roadmap) => {
    setEditingRoadmap(roadmap);
    setTitle(roadmap.title);
    setDescription(roadmap.description);
    setCategory(roadmap.category);
    setIsPublic(roadmap.is_public);
    setShowForm(true);
  };

  const getCategoryInfo = (cat: RoadmapCategory) =>
    CATEGORIES.find(c => c.value === cat) || CATEGORIES[4];

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <LoadingSpinner />
      </div>
    );
  }

  const inputClass = 'w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:ring-2 focus:ring-gray-500 focus:border-transparent';

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('roadmaps.title')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('roadmaps.description')}</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          {t('roadmaps.addRoadmap')}
        </button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-gray-800 rounded-xl p-4 border border-gray-700">
          <p className="text-2xl font-bold text-white">{roadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.totalRoadmaps')}</p>
        </div>
        <div className="bg-gray-800 rounded-xl p-4 border border-gray-700">
          <p className="text-2xl font-bold text-blue-400">{activeRoadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.activeRoadmaps')}</p>
        </div>
        <div className="bg-gray-800 rounded-xl p-4 border border-gray-700">
          <p className="text-2xl font-bold text-green-400">{completedRoadmaps.length}</p>
          <p className="text-sm text-gray-400">{t('roadmaps.completedRoadmaps')}</p>
        </div>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-md">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingRoadmap ? t('roadmaps.editRoadmap') : t('roadmaps.addRoadmap')}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.roadmapTitle')}</label>
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
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.descriptionLabel')}</label>
                <textarea
                  value={description}
                  onChange={e => setDescription(e.target.value)}
                  placeholder={t('roadmaps.descriptionPlaceholder')}
                  rows={3}
                  className={`${inputClass} resize-none`}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">{t('roadmaps.category')}</label>
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
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
                >
                  {t('common.cancel')}
                </button>
                <button
                  type="submit"
                  disabled={saving || !title.trim()}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
                >
                  {saving ? t('common.saving') : editingRoadmap ? t('common.save') : t('roadmaps.create')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Content */}
      {roadmaps.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <svg className="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 6.75V15m6-6v8.25m.503 3.498 4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 0 0-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0Z" />
            </svg>
          </div>
          <p className="text-gray-400">{t('roadmaps.noRoadmaps')}</p>
          <button
            onClick={() => setShowForm(true)}
            className="mt-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            {t('roadmaps.createFirst')}
          </button>
        </div>
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
                    getCategoryInfo={getCategoryInfo}
                    t={t}
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
                    getCategoryInfo={getCategoryInfo}
                    t={t}
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

function RoadmapCard({
  roadmap, onView, onEdit, onDelete, getCategoryInfo, t,
}: {
  roadmap: Roadmap;
  onView: () => void;
  onEdit: (r: Roadmap) => void;
  onDelete: (id: number) => void;
  getCategoryInfo: (cat: RoadmapCategory) => { icon: string; label: string; Icon: LucideIcon };
  t: (key: string) => string;
}) {
  const catInfo = getCategoryInfo(roadmap.category);
  const CategoryIcon = catInfo.Icon;

  return (
    <div
      onClick={onView}
      className="bg-gray-800 border border-gray-700 rounded-xl p-4 hover:border-gray-600 transition-colors cursor-pointer"
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
