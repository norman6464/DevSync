import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Map, MessageSquare, Trophy, Settings, Plus, X, Trash2, Check, ExternalLink, Crown, UserMinus, ArrowLeft } from 'lucide-react';
import { useStudyCircleDetail, useStudyCircleActivity } from '../hooks';
import { useAuthStore } from '../store/authStore';
import { useUserSearch } from '../hooks';
import Avatar from '../components/common/Avatar';

type Tab = 'roadmap' | 'checkin' | 'ranking' | 'settings';

export default function StudyCircleDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const circleId = id ? parseInt(id) : null;

  const {
    circle, loading, saving,
    addMember, removeMember,
    createStep, deleteStep, updateProgress,
    refetch,
  } = useStudyCircleDetail(circleId);

  const {
    progress, checkins, streaks,
    checkin, refetchProgress,
  } = useStudyCircleActivity(circleId);

  const [tab, setTab] = useState<Tab>('roadmap');
  const [stepForm, setStepForm] = useState({ title: '', description: '', resource_url: '' });
  const [showStepForm, setShowStepForm] = useState(false);
  const [checkinContent, setCheckinContent] = useState('');
  const [showAddMember, setShowAddMember] = useState(false);
  const [memberSearch, setMemberSearch] = useState('');
  const { users: searchUsers } = useUserSearch(memberSearch);

  const isOwner = circle?.owner_id === user?.id;

  const handleCreateStep = async () => {
    const result = await createStep({
      title: stepForm.title,
      description: stepForm.description,
      resource_url: stepForm.resource_url,
      order_index: circle?.steps?.length || 0,
    });
    if (result) {
      setStepForm({ title: '', description: '', resource_url: '' });
      setShowStepForm(false);
    }
  };

  const handleCheckin = async () => {
    if (!checkinContent.trim()) return;
    const result = await checkin(checkinContent);
    if (result) {
      setCheckinContent('');
    }
  };

  const handleAddMember = async (userId: number) => {
    await addMember(userId);
    setShowAddMember(false);
    setMemberSearch('');
  };

  const handleToggleProgress = async (stepId: number, currentlyCompleted: boolean) => {
    await updateProgress(stepId, !currentlyCompleted);
    await refetchProgress();
  };

  const tabs: { key: Tab; icon: typeof Map; label: string }[] = [
    { key: 'roadmap', icon: Map, label: t('studyCircle.tabs.roadmap') },
    { key: 'checkin', icon: MessageSquare, label: t('studyCircle.tabs.checkin') },
    { key: 'ranking', icon: Trophy, label: t('studyCircle.tabs.ranking') },
    ...(isOwner ? [{ key: 'settings' as Tab, icon: Settings, label: t('studyCircle.tabs.settings') }] : []),
  ];

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="h-8 bg-gray-800 rounded w-1/3 animate-pulse" />
        <div className="h-40 bg-gray-900 border border-gray-800 rounded-md animate-pulse" />
      </div>
    );
  }

  if (!circle) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-400">Circle not found</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-start gap-3">
        <button onClick={() => navigate('/study-circles')} className="mt-1 text-gray-400 hover:text-white">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div className="flex-1 min-w-0">
          <h1 className="text-xl font-bold text-white truncate">{circle.name}</h1>
          <p className="text-sm text-purple-400">{circle.topic}</p>
          {circle.description && (
            <p className="text-xs text-gray-500 mt-1">{circle.description}</p>
          )}
        </div>
        <div className="flex items-center gap-1 shrink-0">
          {circle.members?.slice(0, 5).map((m) => (
            <Avatar key={m.id} name={m.user?.name || ''} avatarUrl={m.user?.avatar_url} size="xs" />
          ))}
        </div>
      </div>

      {/* Tabs */}
      <div className="flex items-center border-b border-gray-800 gap-1">
        {tabs.map(({ key, icon: Icon, label }) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`flex items-center gap-1.5 px-3 py-2 text-xs font-medium transition-colors relative ${
              tab === key ? 'text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            <Icon className="w-3.5 h-3.5" />
            {label}
            {tab === key && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-purple-500 rounded-t" />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      {tab === 'roadmap' && (
        <div className="space-y-3">
          {/* Steps */}
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
                  (p) => p.step_id === step.id && p.user_id === user?.id
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
                        onClick={() => handleToggleProgress(step.id, isCompleted)}
                        className={`mt-0.5 w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors ${
                          isCompleted
                            ? 'bg-green-500 border-green-500 text-white'
                            : 'border-gray-600 hover:border-purple-500'
                        }`}
                      >
                        {isCompleted && <Check className="w-3 h-3" />}
                      </button>
                      <div className="flex-1 min-w-0">
                        <h4 className={`text-sm font-medium ${isCompleted ? 'text-green-400 line-through' : 'text-white'}`}>
                          {step.title}
                        </h4>
                        {step.description && (
                          <p className="text-xs text-gray-500 mt-1">{step.description}</p>
                        )}
                        {step.resource_url && (
                          <a
                            href={step.resource_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300 mt-1"
                          >
                            <ExternalLink className="w-3 h-3" />
                            {t('studyCircle.steps.resourceUrl')}
                          </a>
                        )}
                        {/* Progress bar showing how many members completed this step */}
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
                          onClick={() => deleteStep(step.id)}
                          className="text-gray-600 hover:text-red-400 transition-colors"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
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

          {/* Step Form */}
          {showStepForm && (
            <div className="bg-gray-900 border border-purple-500/30 rounded-md p-4 space-y-3">
              <input
                type="text"
                value={stepForm.title}
                onChange={(e) => setStepForm({ ...stepForm, title: e.target.value })}
                placeholder={t('studyCircle.steps.title')}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
              />
              <textarea
                value={stepForm.description}
                onChange={(e) => setStepForm({ ...stepForm, description: e.target.value })}
                placeholder={t('studyCircle.description')}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 h-16 resize-none"
              />
              <input
                type="url"
                value={stepForm.resource_url}
                onChange={(e) => setStepForm({ ...stepForm, resource_url: e.target.value })}
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
      )}

      {tab === 'checkin' && (
        <div className="space-y-4">
          {/* Checkin Form */}
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <h3 className="text-sm font-medium text-white mb-2">{t('studyCircle.checkin.title')}</h3>
            <div className="flex gap-2">
              <input
                type="text"
                value={checkinContent}
                onChange={(e) => setCheckinContent(e.target.value)}
                placeholder={t('studyCircle.checkin.placeholder')}
                className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
                onKeyDown={(e) => e.key === 'Enter' && handleCheckin()}
              />
              <button
                onClick={handleCheckin}
                disabled={!checkinContent.trim()}
                className="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors shrink-0"
              >
                {t('studyCircle.checkin.submit')}
              </button>
            </div>
          </div>

          {/* Checkin History */}
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <h3 className="text-sm font-medium text-white mb-3">{t('studyCircle.checkin.history')}</h3>
            {checkins.length === 0 ? (
              <p className="text-xs text-gray-500 text-center py-4">{t('studyCircle.checkin.noCheckins')}</p>
            ) : (
              <div className="space-y-2">
                {checkins.map((ci) => (
                  <div key={ci.id} className="flex items-start gap-2.5 p-2 rounded-lg hover:bg-gray-800/50">
                    <Avatar name={ci.user?.name || ''} avatarUrl={ci.user?.avatar_url} size="xs" />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-medium text-white">{ci.user?.name}</span>
                        <span className="text-[10px] text-gray-600">{ci.date}</span>
                      </div>
                      <p className="text-xs text-gray-400 mt-0.5">{ci.content}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {tab === 'ranking' && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <h3 className="text-sm font-medium text-white mb-3">{t('studyCircle.streak.title')}</h3>
          {streaks.length === 0 ? (
            <p className="text-xs text-gray-500 text-center py-4">{t('studyCircle.streak.noStreaks')}</p>
          ) : (
            <div className="space-y-2">
              {streaks.map((s, i) => (
                <div
                  key={s.user_id}
                  className="flex items-center gap-3 p-2.5 rounded-lg bg-gray-800/30"
                >
                  <span className={`text-sm font-bold w-6 text-center ${
                    i === 0 ? 'text-yellow-400' : i === 1 ? 'text-gray-300' : i === 2 ? 'text-amber-600' : 'text-gray-600'
                  }`}>
                    {i + 1}
                  </span>
                  <Avatar name={s.user_name} avatarUrl={s.avatar_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-white truncate">{s.user_name}</p>
                    <p className="text-[10px] text-gray-500">
                      {t('studyCircle.streak.totalCheckins')}: {s.total_checkins}
                    </p>
                  </div>
                  <div className="text-right shrink-0">
                    <div className="text-lg font-bold text-orange-400">{s.current_streak}</div>
                    <div className="text-[10px] text-gray-500">{t('studyCircle.streak.currentStreak')}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {tab === 'settings' && isOwner && (
        <div className="space-y-4">
          {/* Members Management */}
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-medium text-white">{t('studyCircle.members')}</h3>
              <button
                onClick={() => setShowAddMember(!showAddMember)}
                className="flex items-center gap-1 px-2 py-1 text-xs text-purple-400 hover:text-purple-300 border border-purple-500/30 rounded-lg transition-colors"
              >
                <Plus className="w-3 h-3" />
                {t('studyCircle.addMember')}
              </button>
            </div>

            {showAddMember && (
              <div className="mb-3 p-3 bg-gray-800/50 rounded-lg">
                <input
                  type="text"
                  value={memberSearch}
                  onChange={(e) => setMemberSearch(e.target.value)}
                  placeholder="Search users..."
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 mb-2"
                />
                {searchUsers.length > 0 && (
                  <div className="space-y-1 max-h-32 overflow-y-auto">
                    {searchUsers
                      .filter((u) => !circle.members?.some((m) => m.user_id === u.id))
                      .slice(0, 5)
                      .map((u) => (
                      <button
                        key={u.id}
                        onClick={() => handleAddMember(u.id)}
                        className="flex items-center gap-2 w-full p-2 rounded-lg hover:bg-gray-800 text-left transition-colors"
                      >
                        <Avatar name={u.name} avatarUrl={u.avatar_url} size="xs" />
                        <span className="text-xs text-white">{u.name}</span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}

            <div className="space-y-1.5">
              {circle.members?.map((m) => (
                <div key={m.id} className="flex items-center gap-2.5 p-2 rounded-lg hover:bg-gray-800/50">
                  <Avatar name={m.user?.name || ''} avatarUrl={m.user?.avatar_url} size="sm" />
                  <div className="flex-1 min-w-0">
                    <p className="text-xs text-white truncate">{m.user?.name}</p>
                    <p className="text-[10px] text-gray-500">{m.role === 'owner' ? t('studyCircle.owner') : t('studyCircle.members')}</p>
                  </div>
                  {m.role !== 'owner' && (
                    <button
                      onClick={() => removeMember(m.user_id)}
                      className="text-gray-600 hover:text-red-400 transition-colors"
                      title={t('studyCircle.removeMember')}
                    >
                      <UserMinus className="w-3.5 h-3.5" />
                    </button>
                  )}
                  {m.role === 'owner' && (
                    <Crown className="w-3.5 h-3.5 text-amber-400" />
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Danger Zone */}
          <div className="bg-gray-900 border border-red-500/20 rounded-md p-4">
            <h3 className="text-sm font-medium text-red-400 mb-3">{t('common.delete')}</h3>
            <button
              onClick={async () => {
                if (confirm(t('studyCircle.confirmDelete'))) {
                  const { deleteCircle } = await import('../api/studyCircles');
                  await deleteCircle(circle.id);
                  navigate('/study-circles');
                }
              }}
              className="px-3 py-1.5 bg-red-600/20 hover:bg-red-600/30 text-red-400 rounded-lg text-xs font-medium border border-red-500/30 transition-colors"
            >
              {t('studyCircle.confirmDelete')}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
