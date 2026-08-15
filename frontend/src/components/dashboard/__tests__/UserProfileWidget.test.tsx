import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import UserProfileWidget from '../UserProfileWidget';
import { useAuthStore } from '../../../store/authStore';
import type { User } from '../../../types/user';

const makeUser = (overrides: Partial<User> = {}): User => ({
  id: 1,
  username: 'testuser',
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

const renderWidget = () =>
  render(
    <MemoryRouter>
      <UserProfileWidget />
    </MemoryRouter>
  );

beforeEach(() => {
  useAuthStore.setState({ user: null, isAuthenticated: false });
});

describe('UserProfileWidget', () => {
  it('ユーザーがnullの場合何も表示しない', () => {
    const { container } = renderWidget();
    expect(container.innerHTML).toBe('');
  });

  it('ユーザー名を表示する', () => {
    useAuthStore.setState({ user: makeUser(), isAuthenticated: true });
    renderWidget();
    expect(screen.getByText('テストユーザー')).toBeInTheDocument();
  });

  it('プロフィールリンクが正しいパスを持つ', () => {
    useAuthStore.setState({ user: makeUser(), isAuthenticated: true });
    renderWidget();
    const profileLinks = screen.getAllByRole('link').filter(
      (link) => link.getAttribute('href') === '/profile/testuser'
    );
    expect(profileLinks.length).toBeGreaterThanOrEqual(1);
  });

  it('設定リンクが表示される', () => {
    useAuthStore.setState({ user: makeUser(), isAuthenticated: true });
    renderWidget();
    const settingsLinks = screen.getAllByRole('link').filter(
      (link) => link.getAttribute('href') === '/settings'
    );
    expect(settingsLinks.length).toBeGreaterThanOrEqual(1);
  });

  it('GitHub未連携時にGitHub連携リンクが表示される', () => {
    useAuthStore.setState({ user: makeUser({ github_connected: false }), isAuthenticated: true });
    renderWidget();
    expect(screen.getByText('GitHub と連携する')).toBeInTheDocument();
  });

  it('GitHub連携済み時にGitHub連携リンクが表示されない', () => {
    useAuthStore.setState({ user: makeUser({ github_connected: true }), isAuthenticated: true });
    renderWidget();
    expect(screen.queryByText('GitHub と連携する')).not.toBeInTheDocument();
  });

  it('GitHubユーザー名がある場合に表示される', () => {
    useAuthStore.setState({
      user: makeUser({ github_username: 'octocat', github_connected: true }),
      isAuthenticated: true,
    });
    renderWidget();
    expect(screen.getByText('@octocat')).toBeInTheDocument();
  });

  it('GitHubユーザー名がない場合に@表示されない', () => {
    useAuthStore.setState({ user: makeUser({ github_username: '' }), isAuthenticated: true });
    renderWidget();
    expect(screen.queryByText(/@/)).not.toBeInTheDocument();
  });
});
