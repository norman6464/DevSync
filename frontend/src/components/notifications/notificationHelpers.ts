import type { Notification } from '../../types/notification';

export function getNotificationLink(notification: Notification): string {
  switch (notification.type) {
    case 'post':
    case 'like':
    case 'comment':
    case 'mention':
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
    case 'mention':
      return t('notifications.newMention', { name: notification.actor.name });
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
