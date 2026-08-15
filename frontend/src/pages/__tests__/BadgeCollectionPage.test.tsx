import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import BadgeCollectionPage from '../BadgeCollectionPage';
import * as badgesApi from '../../api/badges';
import type { BadgeResult } from '../../types/badge';

vi.mock('../../api/badges');
vi.mock('../../store/authStore', () => ({
  useAuthStore: vi.fn(() => ({ user: { id: 1, name: 'Test User' } })),
}));

// category はバックエンドが実際に返すカテゴリ名（post / streak / social など）を使う
const mockBadges: BadgeResult[] = [
  {
    id: 'first_post',
    name: '初投稿',
    description: '最初の投稿を作成',
    category: 'post',
    earned: true,
  },
  {
    id: 'week_streak_7',
    name: '7日連続',
    description: '7日連続でログイン',
    category: 'streak',
    earned: true,
  },
  {
    id: 'posts_10',
    name: '投稿10件',
    description: '10件の投稿を作成',
    category: 'post',
    earned: false,
  },
  {
    id: 'followers_100',
    name: 'フォロワー100人',
    description: '100人のフォロワーを獲得',
    category: 'social',
    earned: false,
  },
];

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('BadgeCollectionPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(badgesApi.getUserBadges).mockResolvedValue({
      data: { badges: mockBadges },
    } as any);
  });

  it('ページタイトルが表示される', async () => {
    renderWithRouter(<BadgeCollectionPage />);
    expect(
      await screen.findByRole('heading', { level: 1, name: 'バッジコレクション' })
    ).toBeInTheDocument();
  });

  it('獲得済みバッジが表示される', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      expect(screen.getByText('初投稿')).toBeInTheDocument();
      expect(screen.getByText('7日連続')).toBeInTheDocument();
    });
  });

  it('未獲得バッジが表示される', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      expect(screen.getByText('投稿10件')).toBeInTheDocument();
      expect(screen.getByText('フォロワー100人')).toBeInTheDocument();
    });
  });

  it('バッジの説明が表示される', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      expect(screen.getByText('最初の投稿を作成')).toBeInTheDocument();
    });
  });

  it('進捗パーセンテージが表示される', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      // 4バッジ中2つ獲得 = 50%
      expect(screen.getByText(/50%/)).toBeInTheDocument();
    });
  });

  it('カテゴリフィルターが機能する', async () => {
    const user = userEvent.setup();
    renderWithRouter(<BadgeCollectionPage />);

    await screen.findByText('初投稿');

    // フィルタータブが表示される
    expect(screen.getByRole('button', { name: 'すべて' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '学習' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'ストリーク' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'コミュニティ' })).toBeInTheDocument();

    // ストリークタブで絞り込むと streak カテゴリのバッジだけが表示される
    await user.click(screen.getByRole('button', { name: 'ストリーク' }));
    expect(screen.getByText('7日連続')).toBeInTheDocument();
    expect(screen.queryByText('初投稿')).not.toBeInTheDocument();
    expect(screen.queryByText('投稿10件')).not.toBeInTheDocument();
    expect(screen.queryByText('フォロワー100人')).not.toBeInTheDocument();

    // すべて に戻すと全バッジが再表示される
    await user.click(screen.getByRole('button', { name: 'すべて' }));
    expect(screen.getByText('初投稿')).toBeInTheDocument();
    expect(screen.getByText('フォロワー100人')).toBeInTheDocument();
  });

  it('獲得済みと未獲得のセクションが分かれている', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    // 見出しは「セクション名 (件数)」の形式で描画される
    expect(
      await screen.findByRole('heading', { level: 2, name: '獲得済みバッジ (2)' })
    ).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { level: 2, name: '未獲得バッジ (2)' })
    ).toBeInTheDocument();
  });

  it('ローディング状態が表示される', () => {
    vi.mocked(badgesApi.getUserBadges).mockImplementation(
      () => new Promise(() => {}) as any
    );

    renderWithRouter(<BadgeCollectionPage />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });
});
