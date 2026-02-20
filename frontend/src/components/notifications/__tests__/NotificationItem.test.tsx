import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import NotificationItem from '../NotificationItem';
import type { Notification } from '../../../types/notification';

const wrap = (ui: React.ReactElement) =>
  render(<BrowserRouter>{ui}</BrowserRouter>);

const baseActor = {
  id: 1,
  username: 'testuser',
  name: 'Test User',
  email: 'test@example.com',
  avatar_url: '',
  bio: '',
  github_id: 0,
  github_username: '',
  github_connected: false,
  spotify_connected: false,
  zenn_username: '',
  qiita_username: '',
  atcoder_username: '',
  paiza_rank: '',
  skills_languages: '',
  skills_frameworks: '',
  onboarding_completed: true,
  email_weekly_report: false,
  email_language: 'ja',
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
};

const makeNotification = (overrides: Partial<Notification> = {}): Notification => ({
  id: 1,
  user_id: 2,
  type: 'like',
  actor_id: 1,
  actor: baseActor,
  read: false,
  created_at: '2025-01-01T00:00:00Z',
  ...overrides,
});

describe('NotificationItem', () => {
  const defaultProps = {
    onMarkAsRead: vi.fn(),
    onDelete: vi.fn(),
  };

  it('いいね通知のメッセージを表示する', () => {
    const notification = makeNotification({ type: 'like' });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(screen.getByText('Test Userさんがあなたの投稿にいいねしました')).toBeInTheDocument();
  });

  it('フォロー通知のメッセージを表示する', () => {
    const notification = makeNotification({ type: 'follow' });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(screen.getByText('Test Userさんがあなたをフォローしました')).toBeInTheDocument();
  });

  it('未読通知に青いドットが表示される', () => {
    const notification = makeNotification({ read: false });
    const { container } = wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(container.querySelector('.bg-blue-500')).toBeInTheDocument();
  });

  it('既読通知に青いドットが表示されない', () => {
    const notification = makeNotification({ read: true });
    const { container } = wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(container.querySelector('.bg-blue-500')).not.toBeInTheDocument();
  });

  it('削除ボタンクリック時にonDeleteが呼ばれる', () => {
    const onDelete = vi.fn();
    const notification = makeNotification();
    wrap(<NotificationItem notification={notification} onMarkAsRead={vi.fn()} onDelete={onDelete} />);
    fireEvent.click(screen.getByLabelText('通知を削除'));
    expect(onDelete).toHaveBeenCalledWith(1);
  });

  it('投稿通知にpost.titleが表示される', () => {
    const notification = makeNotification({
      type: 'post',
      post_id: 10,
      post: {
        id: 10, user_id: 1, user: baseActor, title: 'テスト投稿', content: '',
        image_urls: '', is_draft: false, like_count: 0, comment_count: 0,
        bookmark_count: 0, view_count: 0, estimated_read_time: 1,
        created_at: '2025-01-01', updated_at: '2025-01-01',
      },
    });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(screen.getByText('テスト投稿')).toBeInTheDocument();
  });

  it('回答通知にquestion.titleが表示される', () => {
    const notification = makeNotification({
      type: 'answer',
      question_id: 5,
      question: { id: 5, title: 'テスト質問' },
    });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    expect(screen.getByText('テスト質問')).toBeInTheDocument();
  });

  it('投稿通知のリンク先が/posts/:idになる', () => {
    const notification = makeNotification({ type: 'post', post_id: 10 });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/posts/10');
  });

  it('フォロー通知のリンク先が/profile/:usernameになる', () => {
    const notification = makeNotification({ type: 'follow' });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/profile/testuser');
  });

  it('メッセージ通知のリンク先が/chatになる', () => {
    const notification = makeNotification({ type: 'message' });
    wrap(<NotificationItem notification={notification} {...defaultProps} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/chat');
  });
});
