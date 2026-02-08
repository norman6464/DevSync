import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Bell, CheckCheck, Trash2 } from 'lucide-react';
import { useNotifications } from '../hooks';
import type { Notification, NotificationType } from '../types/notification';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';

const FILTER_TYPES: { key: NotificationType | ''; labelKey: string }[] = [
  { key: '', labelKey: 'notifications.filterAll' },
  { key: 'post', labelKey: 'notifications.filterPost' },
  { key: 'like', labelKey: 'notifications.filterLike' },
  { key: 'comment', labelKey: 'notifications.filterComment' },
  { key: 'follow', labelKey: 'notifications.filterFollow' },
  { key: 'message', labelKey: 'notifications.filterMessage' },
  { key: 'answer', labelKey: 'notifications.filterAnswer' },
  { key: 'badge', labelKey: 'notifications.filterBadge' },
];

function getNotificationLink(notification: Notification): string {
  switch (notification.type) {
    case 'post':
    case 'like':
    case 'comment':
      return notification.post_id ? `/posts/${notification.post_id}` : '/';
    case 'follow':
      return `/profile/${notification.actor_id}`;
    case 'message':
      return '/chat';
    case 'answer':
      return notification.question_id ? `/qa/${notification.question_id}` : '/';
    case 'badge':
      return `/profile/${notification.actor_id}`;
    default:
      return '/';
  }
}

function formatTime(dateString: string, t: (key: string, opts?: Record<string, unknown>) => string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return t('notifications.justNow');
  if (diffMins < 60) return t('notifications.minutesAgo', { count: diffMins });
  if (diffHours < 24) return t('notifications.hoursAgo', { count: diffHours });
  return t('notifications.daysAgo', { count: diffDays });
}

export default function NotificationsPage() {
  const { t } = useTranslation();
  const {
    notifications, unreadCount, total, loading,
    page, setPage, limit,
    filterType, setFilterType,
    markAsRead, markAllAsRead, deleteNotification,
  } = useNotifications();

  const getNotificationMessage = (notification: Notification) => {
    switch (notification.type) {
      case 'post':
        return t('notifications.newPost', { name: notification.actor.name });
      case 'message':
        return t('notifications.newMessage', { name: notification.actor.name });
      case 'like':
        return t('notifications.newLike', { name: notification.actor.name });
      case 'comment':
        return t('notifications.newComment', { name: notification.actor.name });
      case 'follow':
        return t('notifications.newFollow', { name: notification.actor.name });
      case 'answer':
        return t('notifications.newAnswer', { name: notification.actor.name });
      case 'badge':
        return t('notifications.newBadge');
      default:
        return '';
    }
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('notifications.pageTitle')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('notifications.pageSubtitle')}</p>
        </div>
        {unreadCount > 0 && (
          <button
            onClick={markAllAsRead}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            <CheckCheck className="w-5 h-5" />
            {t('notifications.markAllRead')}
          </button>
        )}
      </div>

      {/* Filter Tabs */}
      <div className="flex flex-wrap gap-2 mb-6">
        {FILTER_TYPES.map(({ key, labelKey }) => (
          <button
            key={key}
            onClick={() => setFilterType(key)}
            className={`px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
              filterType === key
                ? 'bg-gray-700 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            {t(labelKey)}
          </button>
        ))}
      </div>

      {/* Content */}
      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <LoadingSpinner />
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-16 h-16 mx-auto mb-4 bg-gray-800 rounded-full flex items-center justify-center">
            <Bell className="w-8 h-8 text-gray-500" />
          </div>
          <p className="text-gray-400">{t('notifications.empty')}</p>
        </div>
      ) : (
        <>
          <div className="space-y-2">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className={`flex items-start gap-3 p-4 rounded-lg border transition-colors ${
                  !notification.read
                    ? 'bg-gray-800/50 border-gray-700'
                    : 'bg-gray-900 border-gray-800'
                }`}
              >
                <Link
                  to={getNotificationLink(notification)}
                  onClick={() => {
                    if (!notification.read) markAsRead(notification.id);
                  }}
                  className="flex items-start gap-3 flex-1 min-w-0"
                >
                  <Avatar
                    name={notification.actor.name}
                    avatarUrl={notification.actor.avatar_url}
                    size="sm"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-gray-100">
                      {getNotificationMessage(notification)}
                    </p>
                    {(notification.type === 'post' || notification.type === 'like' || notification.type === 'comment') && notification.post && (
                      <p className="text-xs text-gray-400 truncate mt-0.5">
                        {notification.post.title}
                      </p>
                    )}
                    {notification.type === 'answer' && notification.question && (
                      <p className="text-xs text-gray-400 truncate mt-0.5">
                        {notification.question.title}
                      </p>
                    )}
                    <p className="text-xs text-gray-500 mt-1">
                      {formatTime(notification.created_at, t)}
                    </p>
                  </div>
                  {!notification.read && (
                    <span className="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0" />
                  )}
                </Link>
                <button
                  onClick={() => deleteNotification(notification.id)}
                  className="p-1.5 text-gray-500 hover:text-red-400 transition-colors rounded-md shrink-0"
                  title={t('notifications.deleteNotification')}
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            ))}
          </div>

          {/* Pagination */}
          {total > limit && (
            <div className="flex justify-center gap-2 mt-8">
              <button
                onClick={() => setPage(page - 1)}
                disabled={page <= 1}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.previous')}
              </button>
              <span className="px-4 py-2 text-gray-400">
                {page} / {totalPages}
              </span>
              <button
                onClick={() => setPage(page + 1)}
                disabled={page >= totalPages}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.next')}
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
