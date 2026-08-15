import { useState, useCallback, type ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import QuickStatsWidget from '../components/dashboard/QuickStatsWidget';
import QuickActionsWidget from '../components/dashboard/QuickActionsWidget';
import { useAuthStore } from '../store/authStore';
import { usePosts, useDashboard, useBadgeNotifier, useConfirm, useWidgetSettings } from '../hooks';
import { getUserBadges } from '../api/badges';
import { useAsyncData } from '../hooks/useAsyncData';
import type { BadgeResult } from '../types/badge';
import type { Post } from '../types/post';
import PostCard from '../components/posts/PostCard';
import PostForm from '../components/posts/PostForm';
import { Modal } from '../components/common';
import ConfirmDialog from '../components/common/ConfirmDialog';
import QuickPostForm from '../components/posts/QuickPostForm';
import { emptyStateClass } from '../constants/styles';
import { buttonPrimaryClass } from '../constants/styles';
import { PostCardSkeleton } from '../components/common/Skeleton';
import RecentNotificationsWidget from '../components/dashboard/RecentNotificationsWidget';
import UserProfileWidget from '../components/dashboard/UserProfileWidget';
import LevelWidget from '../components/dashboard/LevelWidget';
import StreakWidget from '../components/dashboard/StreakWidget';
import DailyChallengeWidget from '../components/dashboard/DailyChallengeWidget';
import AIAdviceWidget from '../components/dashboard/AIAdviceWidget';
import RecommendedUsersWidget from '../components/dashboard/RecommendedUsersWidget';
import TrendingWidget from '../components/dashboard/TrendingWidget';
import StudyCircleWidget from '../components/dashboard/StudyCircleWidget';
import GoalsProgressWidget from '../components/dashboard/GoalsProgressWidget';
import QuickEntryWidget from '../components/dashboard/QuickEntryWidget';
import WeeklyChallengeWidget from '../components/dashboard/WeeklyChallengeWidget';
import WidgetSettingsPanel from '../components/dashboard/WidgetSettingsPanel';

export default function DashboardPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const { posts, loading, tab, setTab, createPost, updatePost, deletePost, refetch } = usePosts();
  const { confirm, dialogProps } = useConfirm();
  const [editingPost, setEditingPost] = useState<Post | null>(null);

  const handleDeletePost = useCallback(async (post: Post) => {
    const ok = await confirm({ title: t('common.confirm'), message: t('dashboard.confirmDeletePost'), variant: 'danger' });
    if (ok) deletePost(post);
  }, [confirm, deletePost, t]);
  const {
    activeGoals,
    completedGoals,
    avgProgress,
    goalsLoading,
    recentNotifications,
    notificationsLoading,
  } = useDashboard();

  const {
    widgets,
    editing: widgetEditing,
    localWidgets,
    saving: widgetSaving,
    startEditing: startWidgetEditing,
    cancelEditing: cancelWidgetEditing,
    toggleVisibility,
    moveUp,
    moveDown,
    saveSettings: saveWidgetSettings,
  } = useWidgetSettings();

  // Fetch badges for the current user and detect new acquisitions
  const { data: badges } = useAsyncData(
    async () => {
      if (!user) return [] as BadgeResult[];
      const res = await getUserBadges(user.id);
      return res.data?.badges || [];
    },
    { initialData: [] as BadgeResult[], deps: [user?.id], enabled: !!user }
  );
  useBadgeNotifier(badges);

  const handleQuickPost = useCallback(async (title: string, content: string, isDraft?: boolean) => {
    await createPost(title, content, undefined, undefined, isDraft);
    await refetch();
  }, [createPost, refetch]);

  // Widget component registry
  const widgetComponents: Record<string, ReactNode> = {
    userProfile: <UserProfileWidget />,
    level: <LevelWidget />,
    streak: <StreakWidget />,
    dailyChallenge: <DailyChallengeWidget />,
    weeklyChallenge: <WeeklyChallengeWidget />,
    studyCircle: <StudyCircleWidget />,
    quickEntry: <QuickEntryWidget />,
    quickActions: <QuickActionsWidget />,
    recommendedUsers: <RecommendedUsersWidget />,
    trending: <TrendingWidget />,
    aiAdvice: <AIAdviceWidget />,
    goalsProgress: (
      <GoalsProgressWidget
        activeGoals={activeGoals}
        completedGoals={completedGoals}
        avgProgress={avgProgress}
        loading={goalsLoading}
      />
    ),
    recentNotifications: (
      <RecentNotificationsWidget
        notifications={recentNotifications}
        loading={notificationsLoading}
      />
    ),
    quickStats: (
      <QuickStatsWidget
        activeCount={activeGoals.length}
        completedCount={completedGoals.length}
      />
    ),
  };

  const visibleWidgets = widgets.filter((w) => w.visible);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
      {/* Main Feed */}
      <div className="lg:col-span-3 space-y-6">
        <QuickPostForm onSubmit={handleQuickPost} />

        <div>
          <h2 className="section-heading">{t('dashboard.timeline')}</h2>
          {/* Tabs */}
          <div className="flex items-center border-b-2 border-gray-800" role="tablist" aria-label={t('dashboard.timeline')}>
          <button
            role="tab"
            aria-selected={tab === 'timeline'}
            onClick={() => setTab('timeline')}
            className={`px-4 py-2.5 text-sm font-medium transition-colors relative ${
              tab === 'timeline' ? 'text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('dashboard.following')}
            {tab === 'timeline' && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-t" />
            )}
          </button>
          <button
            role="tab"
            aria-selected={tab === 'all'}
            onClick={() => setTab('all')}
            className={`px-4 py-2.5 text-sm font-medium transition-colors relative ${
              tab === 'all' ? 'text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            {t('dashboard.allPosts')}
            {tab === 'all' && (
              <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-orange-500 rounded-t" />
            )}
          </button>
        </div>

        {loading ? (
          <div className="space-y-3">
            <PostCardSkeleton />
            <PostCardSkeleton />
            <PostCardSkeleton />
          </div>
        ) : posts.length === 0 ? (
          <div className={emptyStateClass}>
            <svg className="w-16 h-16 mx-auto mb-4 text-gray-700" fill="none" stroke="currentColor" strokeWidth="1" viewBox="0 0 24 24" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z" />
            </svg>
            <p className="text-gray-400 mb-4">
              {tab === 'timeline'
                ? t('dashboard.noPostsFollowing')
                : t('dashboard.noPostsAll')}
            </p>
            {tab === 'timeline' && (
              <Link
                to="/search"
                className={`${buttonPrimaryClass} inline-flex items-center gap-2 text-sm`}
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                </svg>
                {t('dashboard.findPeople')}
              </Link>
            )}
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                isOwner={user?.id === post.user_id}
                onEdit={() => setEditingPost(post)}
                onDelete={() => handleDeletePost(post)}
                onUpdate={refetch}
              />
            ))}
          </div>
        )}
        </div>

        {/* Edit Modal */}
        <Modal
          isOpen={!!editingPost}
          onClose={() => setEditingPost(null)}
          title={t('dashboard.editPost')}
        >
          {editingPost && (
            <PostForm
              post={editingPost}
              onSubmit={async (title, content, imageUrls) => {
                const result = await updatePost(editingPost.id, title, content, imageUrls);
                if (result) setEditingPost(null);
              }}
              onCancel={() => setEditingPost(null)}
            />
          )}
        </Modal>
      </div>

      {/* Sidebar */}
      <div className="space-y-4">
        {/* Widget Settings Toggle */}
        {!widgetEditing && (
          <div className="flex justify-end">
            <button
              onClick={startWidgetEditing}
              className="flex items-center gap-1.5 px-2.5 py-1 text-xs text-gray-400 hover:text-white border border-gray-700 hover:border-gray-500 rounded transition-colors"
            >
              <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 6h9.75M10.5 6a1.5 1.5 0 1 1-3 0m3 0a1.5 1.5 0 1 0-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m-9.75 0h9.75" />
              </svg>
              {t('widgetSettings.customize')}
            </button>
          </div>
        )}

        {/* Widget Settings Panel */}
        {widgetEditing && (
          <WidgetSettingsPanel
            widgets={localWidgets}
            saving={widgetSaving}
            onToggleVisibility={toggleVisibility}
            onMoveUp={moveUp}
            onMoveDown={moveDown}
            onSave={saveWidgetSettings}
            onCancel={cancelWidgetEditing}
          />
        )}

        {/* Dynamic Widgets */}
        {!widgetEditing && visibleWidgets.map((widget) => (
          <div key={widget.key}>
            {widgetComponents[widget.key]}
          </div>
        ))}
      </div>
      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
