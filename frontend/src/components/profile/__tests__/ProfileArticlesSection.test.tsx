import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProfileArticlesSection from '../ProfileArticlesSection';
import { type ZennArticle, type ZennStats } from '../../../api/zenn';
import { type QiitaArticle, type QiitaStats } from '../../../api/qiita';

const makeZennArticle = (overrides: Partial<ZennArticle> = {}): ZennArticle => ({
  id: 1,
  user_id: 1,
  zenn_id: 100,
  title: 'Reactの基本',
  slug: 'react-basics',
  emoji: '⚛️',
  article_type: 'tech',
  liked_count: 10,
  comments_count: 3,
  published_at: '2025-01-01',
  updated_at: '2025-01-01',
  ...overrides,
});

const makeQiitaArticle = (overrides: Partial<QiitaArticle> = {}): QiitaArticle => ({
  id: 1,
  user_id: 1,
  qiita_id: 'abc123',
  title: 'Go入門ガイド',
  url: 'https://qiita.com/user/items/abc123',
  likes_count: 5,
  comments_count: 2,
  tags: 'Go,Backend',
  published_at: '2025-01-01',
  updated_at: '2025-01-01',
  ...overrides,
});

const zennStats: ZennStats = { total_articles: 10, total_likes: 50, total_comments: 20 };
const qiitaStats: QiitaStats = { total_articles: 8, total_likes: 30, total_comments: 15 };

const renderSection = (props = {}) =>
  render(
    <MemoryRouter>
      <ProfileArticlesSection
        zennUsername="testuser"
        zennArticles={[makeZennArticle()]}
        zennStats={zennStats}
        qiitaUsername="testuser"
        qiitaArticles={[makeQiitaArticle()]}
        qiitaStats={qiitaStats}
        {...props}
      />
    </MemoryRouter>
  );

describe('ProfileArticlesSection', () => {
  it('Zenn記事のタイトルが表示される', () => {
    renderSection();
    expect(screen.getByText('Reactの基本')).toBeInTheDocument();
  });

  it('Qiita記事のタイトルが表示される', () => {
    renderSection();
    expect(screen.getByText('Go入門ガイド')).toBeInTheDocument();
  });

  it('Zennセクションヘッダーが表示される', () => {
    renderSection();
    expect(screen.getByText('Zenn記事')).toBeInTheDocument();
  });

  it('Qiitaセクションヘッダーが表示される', () => {
    renderSection();
    expect(screen.getByText('Qiita記事')).toBeInTheDocument();
  });

  it('Zenn統計情報が表示される', () => {
    renderSection();
    expect(screen.getByText(/10 記事/)).toBeInTheDocument();
  });

  it('Qiita統計情報が表示される', () => {
    renderSection();
    expect(screen.getByText(/8 記事/)).toBeInTheDocument();
  });

  it('zennUsernameがない場合はZennセクションが非表示', () => {
    renderSection({ zennUsername: undefined });
    expect(screen.queryByText('Zennの記事')).not.toBeInTheDocument();
  });

  it('qiitaUsernameがない場合はQiitaセクションが非表示', () => {
    renderSection({ qiitaUsername: undefined });
    expect(screen.queryByText('Qiita記事')).not.toBeInTheDocument();
  });

  it('Zenn記事が空の場合はZennセクションが非表示', () => {
    renderSection({ zennArticles: [] });
    expect(screen.queryByText('Zennの記事')).not.toBeInTheDocument();
  });

  it('Qiita記事が空の場合はQiitaセクションが非表示', () => {
    renderSection({ qiitaArticles: [] });
    expect(screen.queryByText('Qiita記事')).not.toBeInTheDocument();
  });

  it('Zenn記事のいいね数が表示される', () => {
    renderSection();
    expect(screen.getByText('10')).toBeInTheDocument();
  });
});
