import { describe, it, expect, vi } from 'vitest';
import { getNotificationLink, getNotificationMessage } from '../notificationHelpers';
import type { Notification } from '../../../types/notification';

const baseNotification: Notification = {
  id: 1,
  user_id: 10,
  type: 'post',
  actor_id: 20,
  actor: { id: 20, name: 'Alice', username: 'alice' } as Notification['actor'],
  read: false,
  created_at: '2026-01-01T00:00:00Z',
};

describe('getNotificationLink', () => {
  it('post タイプで post_id ありの場合 /posts/:id を返す', () => {
    const n = { ...baseNotification, type: 'post' as const, post_id: 5 };
    expect(getNotificationLink(n)).toBe('/posts/5');
  });

  it('post タイプで post_id なしの場合 / を返す', () => {
    const n = { ...baseNotification, type: 'post' as const, post_id: undefined };
    expect(getNotificationLink(n)).toBe('/');
  });

  it('like タイプで post_id ありの場合 /posts/:id を返す', () => {
    const n = { ...baseNotification, type: 'like' as const, post_id: 3 };
    expect(getNotificationLink(n)).toBe('/posts/3');
  });

  it('comment タイプで post_id ありの場合 /posts/:id を返す', () => {
    const n = { ...baseNotification, type: 'comment' as const, post_id: 7 };
    expect(getNotificationLink(n)).toBe('/posts/7');
  });

  it('follow タイプの場合 /profile/:username を返す', () => {
    const n = { ...baseNotification, type: 'follow' as const };
    expect(getNotificationLink(n)).toBe('/profile/alice');
  });

  it('message タイプの場合 /chat を返す', () => {
    const n = { ...baseNotification, type: 'message' as const };
    expect(getNotificationLink(n)).toBe('/chat');
  });

  it('answer タイプで question_id ありの場合 /qa/:id を返す', () => {
    const n = { ...baseNotification, type: 'answer' as const, question_id: 9 };
    expect(getNotificationLink(n)).toBe('/qa/9');
  });

  it('answer タイプで question_id なしの場合 / を返す', () => {
    const n = { ...baseNotification, type: 'answer' as const, question_id: undefined };
    expect(getNotificationLink(n)).toBe('/');
  });

  it('badge タイプの場合 /profile/:username を返す', () => {
    const n = { ...baseNotification, type: 'badge' as const };
    expect(getNotificationLink(n)).toBe('/profile/alice');
  });

  it('メンション通知は投稿へ遷移する', () => {
    const n = { type: 'mention', post_id: 9 } as never;
    expect(getNotificationLink(n)).toBe('/posts/9');
  });

  it('メンション通知に投稿がなければルートへ遷移する', () => {
    const n = { type: 'mention' } as never;
    expect(getNotificationLink(n)).toBe('/');
  });

});

describe('getNotificationMessage', () => {
  const t = vi.fn((key: string, opts?: Record<string, string>) =>
    opts ? `${key}:${JSON.stringify(opts)}` : key,
  );

  it('post タイプで正しいキーとactor名を渡す', () => {
    const n = { ...baseNotification, type: 'post' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newPost', { name: 'Alice' });
  });

  it('like タイプで正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'like' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newLike', { name: 'Alice' });
  });

  it('comment タイプで正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'comment' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newComment', { name: 'Alice' });
  });

  it('follow タイプで正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'follow' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newFollow', { name: 'Alice' });
  });

  it('message タイプで正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'message' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newMessage', { name: 'Alice' });
  });

  it('answer タイプで正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'answer' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newAnswer', { name: 'Alice' });
  });

  it('badge タイプで名前なしの正しいキーを渡す', () => {
    const n = { ...baseNotification, type: 'badge' as const };
    getNotificationMessage(n, t);
    expect(t).toHaveBeenCalledWith('notifications.newBadge');
  });
});
