import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import PostCardActions from '../PostCardActions';

vi.mock('../../../api/posts', () => ({
  likePost: vi.fn().mockResolvedValue({}),
  unlikePost: vi.fn().mockResolvedValue({}),
  bookmarkPost: vi.fn().mockResolvedValue({}),
  unbookmarkPost: vi.fn().mockResolvedValue({}),
  getReactions: vi.fn().mockResolvedValue({ data: { reactions: [], user_reactions: [] } }),
  addReaction: vi.fn().mockResolvedValue({}),
  removeReaction: vi.fn().mockResolvedValue({}),
}));

const defaultProps = {
  postId: 1,
  initialLiked: false,
  initialLikeCount: 5,
  initialBookmarked: false,
  bookmarkCount: 2,
  commentCount: 3,
  viewCount: 100,
};

const renderWithRouter = (props = {}) =>
  render(
    <MemoryRouter>
      <PostCardActions {...defaultProps} {...props} />
    </MemoryRouter>
  );

describe('PostCardActions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('いいね数が表示される', () => {
    renderWithRouter();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('コメント数が表示される', () => {
    renderWithRouter();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('閲覧数が表示される', () => {
    renderWithRouter();
    expect(screen.getByText('100')).toBeInTheDocument();
  });

  it('ブックマーク数が表示される', () => {
    renderWithRouter();
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('閲覧数0の場合は非表示', () => {
    renderWithRouter({ viewCount: 0 });
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });

  it('いいねボタンクリックでlikePostが呼ばれる', async () => {
    const { likePost } = await import('../../../api/posts');
    renderWithRouter();
    fireEvent.click(screen.getByLabelText('いいね'));
    await waitFor(() => {
      expect(likePost).toHaveBeenCalledWith(1);
    });
  });

  it('いいね済みボタンクリックでunlikePostが呼ばれる', async () => {
    const { unlikePost } = await import('../../../api/posts');
    renderWithRouter({ initialLiked: true, initialLikeCount: 5 });
    fireEvent.click(screen.getByLabelText('いいね解除'));
    await waitFor(() => {
      expect(unlikePost).toHaveBeenCalledWith(1);
    });
  });

  it('ブックマークボタンクリックでbookmarkPostが呼ばれる', async () => {
    const { bookmarkPost } = await import('../../../api/posts');
    renderWithRouter();
    fireEvent.click(screen.getByLabelText('ブックマーク'));
    await waitFor(() => {
      expect(bookmarkPost).toHaveBeenCalledWith(1);
    });
  });

  it('リアクションピッカーの表示/非表示を切り替え', () => {
    renderWithRouter();
    const reactionBtn = screen.getByRole('button', { name: 'リアクション' });
    fireEvent.click(reactionBtn);
    expect(screen.getByRole('menu')).toBeInTheDocument();
  });

  it('コメントリンクが正しい遷移先を持つ', () => {
    renderWithRouter();
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/posts/1');
  });
});
