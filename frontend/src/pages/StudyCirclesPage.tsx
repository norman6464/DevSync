import { useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users, Plus, Crown } from 'lucide-react';
import { useStudyCircles } from '../hooks';
import { useAuthStore } from '../store/authStore';
import Avatar from '../components/common/Avatar';
import { Modal } from '../components/common';

export default function StudyCirclesPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const { circles, loading, saving, createCircle, deleteCircle } = useStudyCircles();
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({ name: '', topic: '', description: '', max_members: 5 });

  const handleOpenModal = useCallback(() => setShowModal(true), []);
  const handleCloseModal = useCallback(() => setShowModal(false), []);

  const handleCreate = useCallback(async () => {
    const result = await createCircle(form);
    if (result) {
      setShowModal(false);
      setForm({ name: '', topic: '', description: '', max_members: 5 });
    }
  }, [form, createCircle]);

  const handleNameChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, name: e.target.value }));
  }, []);
  const handleTopicChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, topic: e.target.value }));
  }, []);
  const handleDescriptionChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setForm((prev) => ({ ...prev, description: e.target.value }));
  }, []);
  const handleMaxMembersChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, max_members: parseInt(e.target.value) || 5 }));
  }, []);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="flex items-center gap-2 text-xl font-bold text-white">
          <Users className="w-6 h-6 text-purple-400" />
          {t('studyCircle.title')}
        </h1>
        <button
          onClick={handleOpenModal}
          className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg text-sm font-medium transition-colors"
        >
          <Plus className="w-4 h-4" />
          {t('studyCircle.create')}
        </button>
      </div>

      {/* Circle List */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="bg-gray-900 border border-gray-800 rounded-md p-5 animate-pulse">
              <div className="h-5 bg-gray-800 rounded w-1/2 mb-3" />
              <div className="h-4 bg-gray-800 rounded w-3/4 mb-2" />
              <div className="h-3 bg-gray-800 rounded w-1/3" />
            </div>
          ))}
        </div>
      ) : circles.length === 0 ? (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
          <Users className="w-16 h-16 mx-auto mb-4 text-gray-700" />
          <p className="text-gray-400 mb-4">{t('studyCircle.noCircles')}</p>
          <button
            onClick={handleOpenModal}
            className="inline-flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg text-sm font-medium transition-colors"
          >
            <Plus className="w-4 h-4" />
            {t('studyCircle.create')}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {circles.map((circle) => (
            <Link
              key={circle.id}
              to={`/study-circles/${circle.id}`}
              className="bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors group"
            >
              <div className="flex items-start justify-between mb-2">
                <h3 className="text-sm font-semibold text-white group-hover:text-purple-400 transition-colors truncate">
                  {circle.name}
                </h3>
                <span className={`text-[10px] px-2 py-0.5 rounded-full font-medium shrink-0 ml-2 ${
                  circle.status === 'active' ? 'bg-green-500/20 text-green-400' :
                  circle.status === 'completed' ? 'bg-blue-500/20 text-blue-400' :
                  'bg-gray-500/20 text-gray-400'
                }`}>
                  {t(`studyCircle.${circle.status}`)}
                </span>
              </div>
              <p className="text-xs text-purple-400 mb-2">{circle.topic}</p>
              {circle.description && (
                <p className="text-xs text-gray-500 mb-3 line-clamp-2">{circle.description}</p>
              )}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-1">
                  {circle.members?.slice(0, 5).map((member) => (
                    <Avatar
                      key={member.id}
                      name={member.user?.name || ''}
                      avatarUrl={member.user?.avatar_url}
                      size="xs"
                    />
                  ))}
                  {(circle.members?.length || 0) > 5 && (
                    <span className="text-[10px] text-gray-500 ml-1">
                      +{(circle.members?.length || 0) - 5}
                    </span>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-16 h-1.5 bg-gray-700 rounded-full overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all ${
                        (circle.members?.length || 0) >= circle.max_members
                          ? 'bg-red-400'
                          : (circle.members?.length || 0) >= circle.max_members * 0.8
                          ? 'bg-yellow-400'
                          : 'bg-purple-400'
                      }`}
                      style={{ width: `${Math.min(((circle.members?.length || 0) / circle.max_members) * 100, 100)}%` }}
                    />
                  </div>
                  <div className="flex items-center gap-1 text-[10px] text-gray-500">
                    <Users className="w-3 h-3" />
                    {circle.members?.length || 0}/{circle.max_members}
                  </div>
                </div>
              </div>
              {circle.owner_id === user?.id && (
                <div className="flex items-center gap-1 mt-2 text-[10px] text-amber-400">
                  <Crown className="w-3 h-3" />
                  {t('studyCircle.owner')}
                </div>
              )}
            </Link>
          ))}
        </div>
      )}

      {/* Create Modal */}
      <Modal
        isOpen={showModal}
        onClose={handleCloseModal}
        title={t('studyCircle.create')}
        maxWidth="max-w-md"
      >
        <div className="space-y-3">
          <div>
            <label className="text-xs text-gray-400 block mb-1">{t('studyCircle.name')}</label>
            <input
              type="text"
              value={form.name}
              onChange={handleNameChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
              placeholder={t('studyCircle.namePlaceholder')}
            />
          </div>
          <div>
            <label className="text-xs text-gray-400 block mb-1">{t('studyCircle.topic')}</label>
            <input
              type="text"
              value={form.topic}
              onChange={handleTopicChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
              placeholder={t('studyCircle.topicPlaceholder')}
            />
          </div>
          <div>
            <label className="text-xs text-gray-400 block mb-1">{t('studyCircle.description')}</label>
            <textarea
              value={form.description}
              onChange={handleDescriptionChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 h-20 resize-none"
            />
          </div>
          <div>
            <label className="text-xs text-gray-400 block mb-1">{t('studyCircle.maxMembers')}</label>
            <input
              type="number"
              min={3}
              max={10}
              value={form.max_members}
              onChange={handleMaxMembersChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
            />
          </div>
        </div>
        <div className="flex justify-end gap-2 mt-5">
          <button
            onClick={handleCloseModal}
            className="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors"
          >
            {t('common.cancel')}
          </button>
          <button
            onClick={handleCreate}
            disabled={saving || !form.name || !form.topic}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
          >
            {saving ? t('common.saving') : t('studyCircle.create')}
          </button>
        </div>
      </Modal>
    </div>
  );
}
