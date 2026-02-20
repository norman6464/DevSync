import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Map, MessageSquare, Trophy, Settings, ArrowLeft } from 'lucide-react';
import { useStudyCircleDetail, useStudyCircleActivity } from '../hooks';
import { useAuthStore } from '../store/authStore';
import Avatar from '../components/common/Avatar';
import CircleRoadmapTab from '../components/studyCircles/CircleRoadmapTab';
import CircleCheckinTab from '../components/studyCircles/CircleCheckinTab';
import CircleRankingTab from '../components/studyCircles/CircleRankingTab';
import CircleSettingsTab from '../components/studyCircles/CircleSettingsTab';

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
  } = useStudyCircleDetail(circleId);

  const {
    progress, checkins, streaks,
    checkin, refetchProgress,
  } = useStudyCircleActivity(circleId);

  const [tab, setTab] = useState<Tab>('roadmap');
  const isOwner = circle?.owner_id === user?.id;

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
        <button onClick={() => navigate('/study-circles')} className="mt-1 text-gray-400 hover:text-white" aria-label={t('common.back')}>
          <ArrowLeft className="w-5 h-5" aria-hidden="true" />
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
      <div className="flex items-center border-b border-gray-800 gap-1" role="tablist">
        {tabs.map(({ key, icon: Icon, label }) => (
          <button
            key={key}
            role="tab"
            aria-selected={tab === key}
            onClick={() => setTab(key)}
            className={`flex items-center gap-1.5 px-3 py-2 text-xs font-medium transition-colors relative ${
              tab === key ? 'text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            <Icon className="w-3.5 h-3.5" aria-hidden="true" />
            {label}
            {tab === key && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-purple-500 rounded-t" />
            )}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      {tab === 'roadmap' && (
        <CircleRoadmapTab
          circle={circle}
          progress={progress}
          currentUser={user}
          isOwner={isOwner}
          saving={saving}
          onCreateStep={createStep}
          onDeleteStep={deleteStep}
          onToggleProgress={handleToggleProgress}
        />
      )}

      {tab === 'checkin' && (
        <CircleCheckinTab checkins={checkins} onCheckin={checkin} />
      )}

      {tab === 'ranking' && (
        <CircleRankingTab streaks={streaks} />
      )}

      {tab === 'settings' && isOwner && (
        <CircleSettingsTab
          circle={circle}
          onAddMember={addMember}
          onRemoveMember={removeMember}
        />
      )}
    </div>
  );
}
