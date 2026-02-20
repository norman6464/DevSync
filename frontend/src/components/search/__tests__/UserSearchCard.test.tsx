import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import UserSearchCard from '../UserSearchCard';
import type { User } from '../../../types/user';

const mockUser: User = {
  id: 1,
  name: 'Test User',
  username: 'testuser',
  email: 'test@example.com',
  bio: 'This is a test bio',
  avatar_url: '',
  github_username: '',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>);

describe('UserSearchCard', () => {
  it('ユーザー名が表示される', () => {
    renderWithRouter(<UserSearchCard user={mockUser} />);
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('bioが表示される', () => {
    renderWithRouter(<UserSearchCard user={mockUser} />);
    expect(screen.getByText('This is a test bio')).toBeInTheDocument();
  });

  it('bioが空の場合は表示されない', () => {
    const userNoBio = { ...mockUser, bio: '' };
    renderWithRouter(<UserSearchCard user={userNoBio} />);
    expect(screen.queryByText('This is a test bio')).not.toBeInTheDocument();
  });

  it('プロフィールリンクが正しいパスを持つ', () => {
    renderWithRouter(<UserSearchCard user={mockUser} />);
    const links = screen.getAllByRole('link');
    const profileLinks = links.filter(link => link.getAttribute('href') === '/profile/testuser');
    expect(profileLinks.length).toBeGreaterThan(0);
  });

  it('プロフィール表示ボタンが表示される', () => {
    renderWithRouter(<UserSearchCard user={mockUser} />);
    expect(screen.getByText('プロフィール')).toBeInTheDocument();
  });
});
