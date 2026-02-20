import { useState, useCallback, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Bell, TrendingUp, CheckCircle2, Clock } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { usePosts, useDashboard, useBadgeNotifier, useConfirm } from '../hooks';
import { getUserBadges } from '../api/badges';
import { useAsyncData } from '../hooks/useAsyncData';
import type { BadgeResult } from '../types/badge';
import type { Post } from '../types/post';
import PostCard from '../components/posts/PostCard';
import PostForm from '../components/posts/PostForm';
import { Modal } from '../components/common';
import ConfirmDialog from '../components/common/ConfirmDialog';
import QuickPostForm from '../components/posts/QuickPostForm';
import { buttonPrimaryClass } from '../constants/styles';
import { PostCardSkeleton } from '../components/common/Skeleton';
import Avatar from '../components/common/Avatar';
import { formatDistanceToNow } from '../utils/timeFormat';
import LevelWidget from '../components/dashboard/LevelWidget';
import StreakWidget from '../components/dashboard/StreakWidget';
import DailyChallengeWidget from '../components/dashboard/DailyChallengeWidget';
import AIAdviceWidget from '../components/dashboard/AIAdviceWidget';
import RecommendedUsersWidget from '../components/dashboard/RecommendedUsersWidget';
import TrendingWidget from '../components/dashboard/TrendingWidget';
import StudyCircleWidget from '../components/dashboard/StudyCircleWidget';
import GoalsProgressWidget from '../components/dashboard/GoalsProgressWidget';

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

  const handleCreatePost = useCallback(async (
    title: string,
    content: string,
    imageUrls?: string,
    codeSnippets?: { language: string; file_name?: string; code: string }[],
    isDraft?: boolean
  ) => {
    await createPost(title, content, imageUrls, codeSnippets, isDraft);
  }, [createPost]);

  const notificationNameMap = useMemo<Record<string, string>>(() => ({
    post: 'notifications.newPost',
    message: 'notifications.newMessage',
    like: 'notifications.newLike',
    comment: 'notifications.newComment',
    follow: 'notifications.newFollow',
    answer: 'notifications.newAnswer',
    badge: 'notifications.newBadge',
  }), []);

  const getNotificationText = useCallback((notification: { type: string; actor: { name: string } }) => {
    return t(notificationNameMap[notification.type] || 'notifications.newPost', {
      name: notification.actor?.name || '',
    });
  }, [t, notificationNameMap]);

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
          <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
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
      <div className="space-y-6">
        {/* User Profile Section */}
        <div>
          <h2 className="section-heading">{t('dashboard.profile')}</h2>
          {user && (
          <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
            <div className="flex items-center gap-3 mb-3">
              <Avatar name={user.name} avatarUrl={user.avatar_url} size="md" />
              <div className="min-w-0">
                <Link to={`/profile/${user.username}`} className="font-medium text-sm hover:text-blue-400 block truncate">
                  {user.name}
                </Link>
                {user.github_username && (
                  <p className="text-xs text-gray-500 truncate">@{user.github_username}</p>
                )}
              </div>
            </div>
            <div className="space-y-0.5">
              <Link to={`/profile/${user.username}`} className="flex items-center gap-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 py-1.5 px-2 rounded-md transition-colors">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.9 17.9 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" /></svg>
                {t('dashboard.yourProfile')}
              </Link>
              <Link to="/settings" className="flex items-center gap-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 py-1.5 px-2 rounded-md transition-colors">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93l.164.076c.374.174.798.203 1.167.067l.838-.34a1.114 1.114 0 0 1 1.365.486l.547.948c.271.47.163 1.07-.258 1.414l-.69.577a1.18 1.18 0 0 0-.378 1.001c.006.1.006.2 0 .3a1.18 1.18 0 0 0 .378 1.001l.69.577c.421.345.529.944.258 1.414l-.547.948a1.114 1.114 0 0 1-1.365.486l-.838-.34a1.18 1.18 0 0 0-1.167.067l-.164.076c-.396.166-.71.506-.78.93l-.149.894c-.09.542-.56.94-1.11.94h-1.093c-.55 0-1.02-.398-1.11-.94l-.149-.894a1.18 1.18 0 0 0-.78-.93l-.164-.076a1.18 1.18 0 0 0-1.167-.067l-.838.34a1.114 1.114 0 0 1-1.365-.486l-.547-.948a1.114 1.114 0 0 1 .258-1.414l.69-.577a1.18 1.18 0 0 0 .378-1.001 2 2 0 0 1 0-.3 1.18 1.18 0 0 0-.378-1.001l-.69-.577a1.114 1.114 0 0 1-.258-1.414l.547-.948a1.114 1.114 0 0 1 1.365-.486l.838.34c.369.136.793.107 1.167-.067l.164-.076c.396-.166.71-.506.78-.93l.149-.894Z" /><path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" /></svg>
                {t('nav.settings')}
              </Link>
            </div>
            {!user.github_connected && (
              <Link
                to="/settings"
                className="mt-3 flex items-center gap-2 text-sm text-amber-400 hover:text-amber-300 py-2 px-3 bg-amber-400/10 border border-amber-400/20 rounded-lg transition-colors"
              >
                <svg className="w-4 h-4 shrink-0" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                {t('dashboard.connectGitHub')}
              </Link>
            )}
          </div>
        )}
        </div>

        {/* Progress Section */}
        <div>
          <h2 className="section-heading">{t('dashboard.progress')}</h2>
          <div className="space-y-4">
            <LevelWidget />
            <StreakWidget />
          </div>
        </div>

        {/* Activities Section */}
        <div>
          <h2 className="section-heading">{t('dashboard.activities')}</h2>
          <div className="space-y-4">
            <DailyChallengeWidget />
            <StudyCircleWidget />
          </div>
        </div>

        {/* Recommendations Section */}
        <div>
          <h2 className="section-heading">{t('dashboard.recommendations')}</h2>
          <div className="space-y-4">
            <RecommendedUsersWidget />
            <TrendingWidget />
            <AIAdviceWidget />
          </div>
        </div>

        {/* Goals Progress Widget */}
        <GoalsProgressWidget
          activeGoals={activeGoals}
          completedGoals={completedGoals}
          avgProgress={avgProgress}
          loading={goalsLoading}
        />

        {/* Recent Notifications Widget */}
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="flex items-center gap-2 text-sm font-medium text-white">
              <Bell className="w-4 h-4 text-yellow-400" aria-hidden="true" />
              {t('dashboard.recentNotifications')}
            </h3>
            <Link to="/notifications" className="text-xs text-gray-400 hover:text-blue-400 transition-colors">
              {t('dashboard.viewAll')}
            </Link>
          </div>

          {notificationsLoading ? (
            <div className="space-y-2">
              {[1, 2, 3].map((i) => (
                <div key={i} className="h-10 bg-gray-800 rounded animate-pulse" />
              ))}
            </div>
          ) : recentNotifications.length === 0 ? (
            <p className="text-xs text-gray-500 text-center py-4">{t('dashboard.noNotifications')}</p>
          ) : (
            <div className="space-y-1">
              {recentNotifications.slice(0, 5).map((notification) => (
                <div
                  key={notification.id}
                  className={`flex items-start gap-2.5 p-2 rounded-lg transition-colors ${
                    !notification.read ? 'bg-gray-800/50' : ''
                  }`}
                >
                  <Avatar
                    name={notification.actor?.name || ''}
                    avatarUrl={notification.actor?.avatar_url}
                    size="xs"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-xs text-gray-300 leading-relaxed truncate">
                      {getNotificationText(notification)}
                    </p>
                    <p className="text-[10px] text-gray-500 mt-0.5">
                      {formatDistanceToNow(notification.created_at)}
                    </p>
                  </div>
                  {!notification.read && (
                    <span className="w-1.5 h-1.5 bg-blue-500 rounded-full mt-1.5 shrink-0" />
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Quick Stats */}
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
            <TrendingUp className="w-4 h-4 text-green-400" aria-hidden="true" />
            {t('dashboard.quickStats')}
          </h3>
          <div className="space-y-2">
            <Link
              to="/goals"
              className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              <CheckCircle2 className="w-4 h-4 text-green-400" aria-hidden="true" />
              <span className="text-xs text-gray-300 flex-1">{t('dashboard.goalsCompleted')}</span>
              <span className="text-xs font-medium text-white">{completedGoals.length}</span>
            </Link>
            <Link
              to="/goals"
              className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              <Clock className="w-4 h-4 text-blue-400" aria-hidden="true" />
              <span className="text-xs text-gray-300 flex-1">{t('dashboard.goalsInProgress')}</span>
              <span className="text-xs font-medium text-white">{activeGoals.length}</span>
            </Link>
          </div>
        </div>
      </div>
      <ConfirmDialog {...dialogProps} />
    </div>
  );
}
