import { useState, useRef, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useNotifications } from '../../hooks';
import type { Notification } from '../../types/notification';
import Avatar from '../common/Avatar';

export default function NotificationDropdown() {
  const { t } = useTranslation();
  const {
    notifications, unreadCount, loading,
    fetchNotifications, markAsRead, markAllAsRead,
  } = useNotifications();

  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleOpen = async () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      await fetchNotifications();
    }
  };

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
      default:
        return '';
    }
  };

  const getNotificationLink = (notification: Notification) => {
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
      default:
        return '/';
    }
  };

  const formatTime = (dateString: string) => {
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
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={handleOpen}
        className="relative p-2 text-gray-400 hover:text-white transition-colors rounded-md"
        title={t('notifications.title')}
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 0 0 5.454-1.31A8.967 8.967 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.967 8.967 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0" />
        </svg>
        {unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 bg-red-500 text-white text-xs w-5 h-5 flex items-center justify-center rounded-full font-medium">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-gray-800 border border-gray-700 rounded-lg shadow-xl z-50">
          <div className="flex items-center justify-between p-3 border-b border-gray-700">
            <h3 className="font-semibold text-white">{t('notifications.title')}</h3>
            {unreadCount > 0 && (
              <button
                onClick={markAllAsRead}
                className="text-xs text-blue-400 hover:text-blue-300"
              >
                {t('notifications.markAllRead')}
              </button>
            )}
          </div>

          <div className="max-h-96 overflow-y-auto">
            {loading ? (
              <div className="p-4 text-center text-gray-400">
                <div className="animate-spin w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full mx-auto" />
              </div>
            ) : notifications.length === 0 ? (
              <div className="p-4 text-center text-gray-400">
                {t('notifications.empty')}
              </div>
            ) : (
              notifications.map((notification) => (
                <Link
                  key={notification.id}
                  to={getNotificationLink(notification)}
                  onClick={() => {
                    if (!notification.read) markAsRead(notification.id);
                    setIsOpen(false);
                  }}
                  className={`flex items-start gap-3 p-3 hover:bg-gray-700/50 transition-colors border-b border-gray-700/50 last:border-0 ${
                    !notification.read ? 'bg-gray-700/30' : ''
                  }`}
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
                      {formatTime(notification.created_at)}
                    </p>
                  </div>
                  {!notification.read && (
                    <span className="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0" />
                  )}
                </Link>
              ))
            )}
          </div>

          <Link
            to="/notifications"
            onClick={() => setIsOpen(false)}
            className="block text-center py-2 text-sm text-blue-400 hover:text-blue-300 border-t border-gray-700"
          >
            {t('notifications.viewAll')}
          </Link>
        </div>
      )}
    </div>
  );
}
