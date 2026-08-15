import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import ProfileHighlightCard from '../ProfileHighlightCard';
import type { Post } from '../../../types/post';
import type { BadgeResult } from '../../../types/badge';
import type { User } from '../../../types/user';

const wrap = (ui: React.ReactElement) => render(<BrowserRouter>{ui}</BrowserRouter>);

const dummyUser = { id: 1, username: 'tester', name: 'Tester' } as User;

const makePost = (overrides: Partial<Post> = {}): Post => ({
  id: 1,
  user_id: 1,
  user: dummyUser,
  title: '投稿',
  content: '本文',
  image_urls: '',
  is_draft: false,
  like_count: 0,
  comment_count: 0,
  bookmark_count: 0,
  view_count: 0,
  estimated_read_time: 1,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  ...overrides,
});

const makeBadge = (overrides: Partial<BadgeResult> = {}): BadgeResult => ({
  id: 'first-post',
  name: '初投稿',
  description: '最初の投稿をした',
  category: 'post',
  earned: true,
  ...overrides,
});

describe('ProfileHighlightCard', () => {
  // like_count / comment_count を見ないと全投稿のスコアが 0 になり、人気の投稿が出なくなる。
  it('反応が最も多い投稿をハイライトする', () => {
    const posts = [
      makePost({ id: 1, title: '反応が少ない投稿', like_count: 1, comment_count: 0 }),
      makePost({ id: 2, title: '最も反応が多い投稿', like_count: 5, comment_count: 3 }),
      makePost({ id: 3, title: 'ふつうの投稿', like_count: 2, comment_count: 2 }),
    ];

    wrap(<ProfileHighlightCard posts={posts} badges={[]} streakDays={0} />);

    expect(screen.getByText('最も反応が多い投稿')).toBeInTheDocument();
    expect(screen.queryByText('反応が少ない投稿')).not.toBeInTheDocument();
  });

  it('反応数を投稿の実データで表示する', () => {
    const posts = [makePost({ title: '投稿A', like_count: 5, comment_count: 3 })];

    wrap(<ProfileHighlightCard posts={posts} badges={[]} streakDays={0} />);

    expect(screen.getByText(/5/)).toBeInTheDocument();
    expect(screen.getByText(/3/)).toBeInTheDocument();
  });

  // 一覧には未獲得バッジも含まれるため、絞らないと未取得のものを「最新バッジ」として出してしまう。
  it('未獲得バッジを最新バッジとして表示しない', () => {
    const badges = [
      makeBadge({ id: 'locked', name: '未獲得バッジ', earned: false }),
      makeBadge({ id: 'earned', name: '獲得済みバッジ', earned: true }),
    ];

    wrap(<ProfileHighlightCard posts={[]} badges={badges} streakDays={0} />);

    expect(screen.getByText('獲得済みバッジ')).toBeInTheDocument();
    expect(screen.queryByText('未獲得バッジ')).not.toBeInTheDocument();
  });

  it('獲得済みバッジが無ければバッジを表示しない', () => {
    const badges = [makeBadge({ name: '未獲得バッジ', earned: false })];

    wrap(<ProfileHighlightCard posts={[]} badges={badges} streakDays={5} />);

    expect(screen.queryByText('未獲得バッジ')).not.toBeInTheDocument();
  });

  it('連続学習日数を表示する', () => {
    wrap(<ProfileHighlightCard posts={[]} badges={[]} streakDays={7} />);

    expect(screen.getByText(/7/)).toBeInTheDocument();
  });

  it('ハイライトが何も無ければ描画しない', () => {
    const { container } = wrap(<ProfileHighlightCard posts={[]} badges={[]} streakDays={0} />);

    expect(container).toBeEmptyDOMElement();
  });
});
