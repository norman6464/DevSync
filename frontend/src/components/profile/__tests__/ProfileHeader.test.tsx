import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProfileHeader from '../ProfileHeader';
import type { User } from '../../../types/user';

vi.mock('../FollowButton', () => ({
  default: ({ userId }: { userId: number }) => (
    <button data-testid="follow-button">Follow {userId}</button>
  ),
}));

const baseUser: User = {
  id: 1,
  username: 'testuser',
  name: 'テストユーザー',
  email: 'test@example.com',
  bio: 'テスト自己紹介',
  avatar_url: '',
  github_username: 'testgh',
  zenn_username: 'testzenn',
  qiita_username: '',
  atcoder_username: 'testatcoder',
  skills_languages: '',
  skills_frameworks: '',
  skills_infrastructure: '',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const defaultProps = {
  user: baseUser,
  isOwnProfile: false,
  followerCount: 10,
  followingCount: 5,
  onShareClick: vi.fn(),
  onPortfolioClick: vi.fn(),
};

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('ProfileHeader', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('ユーザー名が表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('テストユーザー')).toBeInTheDocument();
  });

  it('自己紹介が表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('テスト自己紹介')).toBeInTheDocument();
  });

  it('bioが空の場合は自己紹介が表示されない', () => {
    const noBioUser = { ...baseUser, bio: '' };
    renderWithRouter(<ProfileHeader {...defaultProps} user={noBioUser} />);
    expect(screen.queryByText('テスト自己紹介')).not.toBeInTheDocument();
  });

  it('フォロワー数・フォロー数が表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('他人のプロフィールでFollowButtonが表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} isOwnProfile={false} />);
    expect(screen.getByTestId('follow-button')).toBeInTheDocument();
  });

  it('自分のプロフィールでFollowButtonが表示されない', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} isOwnProfile={true} />);
    expect(screen.queryByTestId('follow-button')).not.toBeInTheDocument();
  });

  it('自分のプロフィールでポートフォリオボタンが表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} isOwnProfile={true} />);
    expect(screen.getByText('ポートフォリオ')).toBeInTheDocument();
  });

  it('他人のプロフィールでポートフォリオボタンが表示されない', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} isOwnProfile={false} />);
    expect(screen.queryByText('ポートフォリオ')).not.toBeInTheDocument();
  });

  it('共有ボタンクリックでonShareClickが呼ばれる', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    fireEvent.click(screen.getByText('共有'));
    expect(defaultProps.onShareClick).toHaveBeenCalledTimes(1);
  });

  it('ポートフォリオボタンクリックでonPortfolioClickが呼ばれる', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} isOwnProfile={true} />);
    fireEvent.click(screen.getByText('ポートフォリオ'));
    expect(defaultProps.onPortfolioClick).toHaveBeenCalledTimes(1);
  });

  it('GitHubリンクが表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('@testgh')).toBeInTheDocument();
    const link = screen.getByText('@testgh').closest('a');
    expect(link).toHaveAttribute('href', 'https://github.com/testgh');
    expect(link).toHaveAttribute('target', '_blank');
  });

  it('Zennリンクが表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('@testzenn')).toBeInTheDocument();
    const link = screen.getByText('@testzenn').closest('a');
    expect(link).toHaveAttribute('href', 'https://zenn.dev/testzenn');
  });

  it('Qiitaユーザー名が空の場合リンクが表示されない', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    const links = screen.queryAllByText(/@/);
    const qiitaLink = links.find(el => el.textContent?.includes('qiita'));
    expect(qiitaLink).toBeUndefined();
  });

  it('AtCoderリンクが表示される', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    expect(screen.getByText('@testatcoder')).toBeInTheDocument();
    const link = screen.getByText('@testatcoder').closest('a');
    expect(link).toHaveAttribute('href', 'https://atcoder.jp/users/testatcoder');
  });

  it('フォロワーリンクの遷移先が正しい', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    const followerLink = screen.getByText('10').closest('a');
    expect(followerLink).toHaveAttribute('href', '/profile/testuser/followers');
  });

  it('フォローリンクの遷移先が正しい', () => {
    renderWithRouter(<ProfileHeader {...defaultProps} />);
    const followingLink = screen.getByText('5').closest('a');
    expect(followingLink).toHaveAttribute('href', '/profile/testuser/following');
  });
});
