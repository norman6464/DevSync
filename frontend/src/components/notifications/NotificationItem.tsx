import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Trash2, Heart, MessageCircle, UserPlus, Mail, FileText, Award, HelpCircle, type LucideIcon } from 'lucide-react';
import type { Notification } from '../../types/notification';

function getTypeIcon(type: Notification['type']): { icon: LucideIcon; color: string } {
  switch (type) {
    case 'like':
      return { icon: Heart, color: 'text-pink-400' };
    case 'comment':
      return { icon: MessageCircle, color: 'text-blue-400' };
    case 'follow':
      return { icon: UserPlus, color: 'text-purple-400' };
    case 'message':
      return { icon: Mail, color: 'text-yellow-400' };
    case 'post':
      return { icon: FileText, color: 'text-gray-400' };
    case 'badge':
      return { icon: Award, color: 'text-orange-400' };
    case 'answer':
      return { icon: HelpCircle, color: 'text-green-400' };
    default:
      return { icon: FileText, color: 'text-gray-400' };
  }
}
import Avatar from '../common/Avatar';
import { formatDistanceToNow } from '../../utils/timeFormat';

export function getNotificationLink(notification: Notification): string {
  switch (notification.type) {
    case 'post':
    case 'like':
    case 'comment':
      return notification.post_id ? `/posts/${notification.post_id}` : '/';
    case 'follow':
      return `/profile/${notification.actor.username}`;
    case 'message':
      return '/chat';
    case 'answer':
      return notification.question_id ? `/qa/${notification.question_id}` : '/';
    case 'badge':
      return `/profile/${notification.actor.username}`;
    default:
      return '/';
  }
}

export function getNotificationMessage(notification: Notification, t: (key: string, opts?: Record<string, string>) => string): string {
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
}

interface NotificationItemProps {
  notification: Notification;
  onMarkAsRead: (id: number) => void;
  onDelete: (id: number) => void;
}

export default function NotificationItem({ notification, onMarkAsRead, onDelete }: NotificationItemProps) {
  const { t } = useTranslation();

  return (
    <li
      className={`flex items-start gap-3 p-4 rounded-lg border transition-colors ${
        !notification.read
          ? 'bg-gray-800/50 border-gray-700'
          : 'bg-gray-900 border-gray-800'
      }`}
    >
      <Link
        to={getNotificationLink(notification)}
        onClick={() => {
          if (!notification.read) onMarkAsRead(notification.id);
        }}
        className="flex items-start gap-3 flex-1 min-w-0"
      >
        <div className="relative shrink-0">
          <Avatar
            name={notification.actor.name}
            avatarUrl={notification.actor.avatar_url}
            size="sm"
          />
          {(() => {
            const { icon: Icon, color } = getTypeIcon(notification.type);
            return (
              <span className={`absolute -bottom-1 -right-1 w-4 h-4 rounded-full bg-gray-900 flex items-center justify-center ring-1 ring-gray-700`}>
                <Icon className={`w-2.5 h-2.5 ${color}`} aria-hidden="true" />
              </span>
            );
          })()}
        </div>
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
          <span className="w-2 h-2 bg-blue-500 rounded-full mt-2 shrink-0" aria-hidden="true" />
        )}
      </Link>
      <button
        onClick={() => onDelete(notification.id)}
        className="p-1.5 text-gray-500 hover:text-red-400 transition-colors rounded-md shrink-0"
        aria-label={t('notifications.deleteNotification')}
      >
        <Trash2 className="w-4 h-4" aria-hidden="true" />
      </button>
    </li>
  );
}
