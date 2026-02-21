import { useState, useRef, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Bell } from 'lucide-react';
import { useNotifications } from '../../hooks';
import Avatar from '../common/Avatar';
import { getNotificationLink, getNotificationMessage } from './NotificationItem';
import { formatDistanceToNow } from '../../utils/timeFormat';
import { linkSmallClass } from '../../constants/styles';

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

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={handleOpen}
        className="relative p-2 text-gray-400 hover:text-white transition-colors rounded-md"
        aria-label={t('notifications.title')}
        aria-expanded={isOpen}
        aria-haspopup="dialog"
      >
        <Bell className="w-5 h-5" aria-hidden="true" />
        {unreadCount > 0 && (
          <span className="absolute -top-0.5 -right-0.5 bg-red-500 text-white text-xs w-5 h-5 flex items-center justify-center rounded-full font-medium">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {isOpen && (
        <div
          className="absolute right-0 mt-2 w-80 bg-gray-800 border border-gray-700 rounded-lg shadow-sm z-50"
          role="dialog"
          aria-label={t('notifications.title')}
        >
          <div className="flex items-center justify-between p-3 border-b border-gray-700">
            <h3 className="font-semibold text-white">{t('notifications.title')}</h3>
            {unreadCount > 0 && (
              <button
                onClick={markAllAsRead}
                className={linkSmallClass}
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
                      {getNotificationMessage(notification, t)}
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
                      {formatDistanceToNow(notification.created_at)}
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
