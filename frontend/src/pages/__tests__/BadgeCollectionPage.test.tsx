import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import BadgeCollectionPage from '../BadgeCollectionPage';
import * as badgesApi from '../../api/badges';
import type { BadgeResult } from '../../types/badge';

vi.mock('../../api/badges');
vi.mock('../../store/authStore', () => ({
  useAuthStore: vi.fn(() => ({ user: { id: 1, name: 'Test User' } })),
}));

const mockBadges: BadgeResult[] = [
  {
    id: 'first_post',
    name: '初投稿',
    description: '最初の投稿を作成',
    category: 'learning',
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
    category: 'learning',
    earned: false,
  },
  {
    id: 'followers_100',
    name: 'フォロワー100人',
    description: '100人のフォロワーを獲得',
    category: 'community',
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
    expect(screen.getByText('バッジコレクション')).toBeInTheDocument();
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
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      expect(screen.getByText('すべて')).toBeInTheDocument();
      expect(screen.getByText('学習')).toBeInTheDocument();
      expect(screen.getByText('ストリーク')).toBeInTheDocument();
      expect(screen.getByText('コミュニティ')).toBeInTheDocument();
    });
  });

  it('獲得済みと未獲得のセクションが分かれている', async () => {
    renderWithRouter(<BadgeCollectionPage />);

    await waitFor(() => {
      expect(screen.getByText('獲得済みバッジ')).toBeInTheDocument();
      expect(screen.getByText('未獲得バッジ')).toBeInTheDocument();
    });
  });

  it('ローディング状態が表示される', () => {
    vi.mocked(badgesApi.getUserBadges).mockImplementation(
      () => new Promise(() => {}) as any
    );

    renderWithRouter(<BadgeCollectionPage />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });
});
