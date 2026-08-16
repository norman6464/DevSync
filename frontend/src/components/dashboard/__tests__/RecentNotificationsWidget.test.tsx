import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import RecentNotificationsWidget from '../RecentNotificationsWidget';
import type { Notification } from '../../../types/notification';
import type { User } from '../../../types/user';

const makeUser = (overrides: Partial<User> = {}): User => ({
  id: 2,
  username: 'test',
  name: 'テストユーザー',
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
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  ...overrides,
});

const makeNotification = (overrides: Partial<Notification> = {}): Notification => ({
  id: 1,
  user_id: 1,
  type: 'like',
  actor_id: 2,
  actor: makeUser(),
  read: false,
  created_at: new Date().toISOString(),
  ...overrides,
});

const renderWidget = (props: Partial<React.ComponentProps<typeof RecentNotificationsWidget>> = {}) => {
  const defaultProps = {
    notifications: [] as Notification[],
    loading: false,
    ...props,
  };
  return render(
    <MemoryRouter>
      <RecentNotificationsWidget {...defaultProps} />
    </MemoryRouter>
  );
};

describe('RecentNotificationsWidget', () => {
  it('ローディング中はスケルトンが表示される', () => {
    renderWidget({ loading: true });
    const skeletons = document.querySelectorAll('.animate-pulse');
    expect(skeletons.length).toBe(3);
  });

  it('通知0件の場合は空メッセージが表示される', () => {
    renderWidget({ notifications: [] });
    expect(screen.getByText('通知はありません')).toBeInTheDocument();
  });

  it('ヘッダーに「最近の通知」が表示される', () => {
    renderWidget();
    expect(screen.getByText('最近の通知')).toBeInTheDocument();
  });

  it('「すべて見る」リンクが/notificationsに遷移する', () => {
    renderWidget();
    const link = screen.getByText('すべて見る');
    expect(link.closest('a')).toHaveAttribute('href', '/notifications');
  });

  it('通知テキストが表示される（likeタイプ）', () => {
    const notifications = [makeNotification({ id: 1, type: 'like', actor: { id: 2, name: '田中太郎', username: 'tanaka', email: '', avatar_url: '', bio: '', github_username: '', github_connected: false, created_at: '', updated_at: '' } as User })];
    renderWidget({ notifications });
    expect(screen.getByText(/田中太郎/)).toBeInTheDocument();
  });

  it('未読通知にはインジケータが表示される', () => {
    const notifications = [makeNotification({ id: 1, read: false })];
    renderWidget({ notifications });
    const indicator = document.querySelector('.bg-blue-500.rounded-full');
    expect(indicator).toBeInTheDocument();
  });

  it('既読通知にはインジケータが表示されない', () => {
    const notifications = [makeNotification({ id: 1, read: true })];
    renderWidget({ notifications });
    const indicator = document.querySelector('.bg-blue-500.rounded-full');
    expect(indicator).not.toBeInTheDocument();
  });

  it('最大5件まで表示される', () => {
    const notifications = Array.from({ length: 7 }, (_, i) =>
      makeNotification({
        id: i + 1,
        actor: { id: i + 2, name: `ユーザー${i + 1}`, username: `user${i}`, email: '', avatar_url: '', bio: '', github_username: '', github_connected: false, created_at: '', updated_at: '' } as User,
      })
    );
    renderWidget({ notifications });
    // 5件分のテキストが存在し、6件目以降は存在しない
    expect(screen.getByText(/ユーザー1/)).toBeInTheDocument();
    expect(screen.getByText(/ユーザー5/)).toBeInTheDocument();
    expect(screen.queryByText(/ユーザー6/)).not.toBeInTheDocument();
  });

  it('未読通知の行にbg-gray-800/50クラスがある', () => {
    const notifications = [makeNotification({ id: 1, read: false })];
    renderWidget({ notifications });
    const row = screen.getByText(/テストユーザー/).closest('.flex.items-start');
    expect(row?.className).toContain('bg-gray-800/50');
  });
});
