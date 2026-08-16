import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import PostSearchCard from '../PostSearchCard';
import type { Post } from '../../../types/post';
import type { User } from '../../../types/user';

const mockPost = {
  id: 1,
  user_id: 1,
  title: 'Test Post Title',
  content: 'This is the post content for testing.',
  like_count: 5,
  comment_count: 3,
  view_count: 100,
  created_at: '2026-01-15T10:00:00Z',
  updated_at: '2026-01-15T10:00:00Z',
  user: {
    id: 1,
    name: 'Author Name',
    username: 'author',
    email: 'author@example.com',
    bio: '',
    avatar_url: '',
    github_username: '',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  } as User,
} as Post;

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>);

describe('PostSearchCard', () => {
  it('投稿タイトルが表示される', () => {
    renderWithRouter(<PostSearchCard post={mockPost} />);
    expect(screen.getByText('Test Post Title')).toBeInTheDocument();
  });

  it('投稿内容が表示される', () => {
    renderWithRouter(<PostSearchCard post={mockPost} />);
    expect(screen.getByText('This is the post content for testing.')).toBeInTheDocument();
  });

  it('著者名が表示される', () => {
    renderWithRouter(<PostSearchCard post={mockPost} />);
    expect(screen.getByText('Author Name')).toBeInTheDocument();
  });

  it('いいね数が表示される', () => {
    renderWithRouter(<PostSearchCard post={mockPost} />);
    expect(screen.getByText('5件のいいね')).toBeInTheDocument();
  });

  it('投稿へのリンクが正しいパスを持つ', () => {
    renderWithRouter(<PostSearchCard post={mockPost} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/posts/1');
  });
});
