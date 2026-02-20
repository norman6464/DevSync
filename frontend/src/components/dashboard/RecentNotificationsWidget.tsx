import { useMemo, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Bell } from 'lucide-react';
import Avatar from '../common/Avatar';
import { formatDistanceToNow } from '../../utils/timeFormat';
import type { Notification } from '../../types/notification';

interface RecentNotificationsWidgetProps {
  notifications: Notification[];
  loading: boolean;
}

export default function RecentNotificationsWidget({ notifications, loading }: RecentNotificationsWidgetProps) {
  const { t } = useTranslation();

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

      {loading ? (
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-10 bg-gray-800 rounded animate-pulse" />
          ))}
        </div>
      ) : notifications.length === 0 ? (
        <p className="text-xs text-gray-500 text-center py-4">{t('dashboard.noNotifications')}</p>
      ) : (
        <div className="space-y-1">
          {notifications.slice(0, 5).map((notification) => (
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
  );
}
