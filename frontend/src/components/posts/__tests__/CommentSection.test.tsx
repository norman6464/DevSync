import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import CommentSection from '../CommentSection';
import type { Comment } from '../../../types/post';

vi.mock('../../../api/users', () => ({
  getUsers: vi.fn().mockResolvedValue({ data: [] }),
}));

vi.mock('../../../api/posts', () => ({
  likeComment: vi.fn().mockResolvedValue({}),
  unlikeComment: vi.fn().mockResolvedValue({}),
}));

const makeComment = (overrides: Partial<Comment> = {}): Comment => ({
  id: 1,
  user_id: 10,
  user: { id: 10, name: 'TestUser', username: 'testuser', avatar_url: '', email: '' } as Comment['user'],
  post_id: 100,
  content: 'テストコメント',
  like_count: 0,
  liked: false,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  ...overrides,
});

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>);

const defaultProps = {
  comments: [] as Comment[],
  submitting: false,
  onSubmitComment: vi.fn().mockResolvedValue(true),
};

describe('CommentSection', () => {
  it('コメント数を表示する', () => {
    const comments = [makeComment({ id: 1 }), makeComment({ id: 2, content: '2番目' })];
    renderWithRouter(<CommentSection {...defaultProps} comments={comments} />);
    expect(screen.getByText('コメント (2)')).toBeInTheDocument();
  });

  it('コメントがない場合は空メッセージを表示する', () => {
    renderWithRouter(<CommentSection {...defaultProps} />);
    expect(screen.getByText('コメントはまだありません')).toBeInTheDocument();
  });

  it('コメント入力フォームを表示する', () => {
    renderWithRouter(<CommentSection {...defaultProps} />);
    expect(screen.getByPlaceholderText('コメントを書く...')).toBeInTheDocument();
  });

  it('送信ボタンを表示する', () => {
    renderWithRouter(<CommentSection {...defaultProps} />);
    expect(screen.getByRole('button', { name: 'コメント' })).toBeInTheDocument();
  });

  it('空入力時は送信ボタンが無効', () => {
    renderWithRouter(<CommentSection {...defaultProps} />);
    expect(screen.getByRole('button', { name: 'コメント' })).toBeDisabled();
  });

  it('submitting中は送信ボタンが無効', () => {
    renderWithRouter(<CommentSection {...defaultProps} submitting={true} />);
    expect(screen.getByRole('button', { name: '投稿中...' })).toBeDisabled();
  });

  it('コメント一覧を表示する', () => {
    const comments = [
      makeComment({ id: 1, content: 'コメント1' }),
      makeComment({ id: 2, content: 'コメント2' }),
    ];
    renderWithRouter(<CommentSection {...defaultProps} comments={comments} />);
    expect(screen.getByText('コメント1')).toBeInTheDocument();
    expect(screen.getByText('コメント2')).toBeInTheDocument();
  });

  it('返信を表示する', () => {
    const comments = [
      makeComment({
        id: 1,
        content: '親コメント',
        replies: [
          makeComment({ id: 2, content: '返信コメント', parent_id: 1 }),
        ],
      }),
    ];
    renderWithRouter(<CommentSection {...defaultProps} comments={comments} />);
    expect(screen.getByText('親コメント')).toBeInTheDocument();
    expect(screen.getByText('返信コメント')).toBeInTheDocument();
  });

  it('返信ボタンを表示する', () => {
    const comments = [makeComment({ id: 1 })];
    renderWithRouter(<CommentSection {...defaultProps} comments={comments} />);
    expect(screen.getByRole('button', { name: '返信' })).toBeInTheDocument();
  });

  it('コメントのユーザー名を表示する', () => {
    const comments = [makeComment({ id: 1 })];
    renderWithRouter(<CommentSection {...defaultProps} comments={comments} />);
    expect(screen.getByText('TestUser')).toBeInTheDocument();
  });
});
