import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import {
  LoadingResults,
  SearchEmptyState,
  UserResults,
  PostResults,
  CircleResults,
} from '../SearchResults';
import type { User } from '../../../types/user';
import type { Post } from '../../../types/post';
import type { StudyCircle } from '../../../types/studyCircle';

const wrap = (ui: React.ReactElement) =>
  render(<BrowserRouter>{ui}</BrowserRouter>);

const makeUser = (id: number): User => ({
  id,
  username: `user${id}`,
  name: `User ${id}`,
  email: `user${id}@example.com`,
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
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
});

const makePost = (id: number): Post => ({
  id,
  user_id: 1,
  user: makeUser(1),
  title: `Post ${id}`,
  content: `Content ${id}`,
  image_urls: '',
  is_draft: false,
  like_count: 0,
  comment_count: 0,
  bookmark_count: 0,
  view_count: 0,
  estimated_read_time: 1,
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
});

const makeCircle = (id: number): StudyCircle => ({
  id,
  name: `Circle ${id}`,
  topic: `Topic ${id}`,
  description: `Description ${id}`,
  owner_id: 1,
  max_members: 10,
  status: 'active',
  created_at: '2025-01-01',
  updated_at: '2025-01-01',
});

describe('SearchResults', () => {
  it('SearchEmptyStateがメッセージを表示する', () => {
    wrap(<SearchEmptyState message="検索してください" />);
    expect(screen.getByText('検索してください')).toBeInTheDocument();
  });

  it('LoadingResultsがusersタブでスケルトンを表示する', () => {
    const { container } = wrap(<LoadingResults tab="users" />);
    expect(container.querySelector('.grid')).toBeInTheDocument();
  });

  it('LoadingResultsがpostsタブでスケルトンを表示する', () => {
    const { container } = wrap(<LoadingResults tab="posts" />);
    expect(container.querySelector('.space-y-4')).toBeInTheDocument();
  });

  it('UserResultsが空の場合NoResultsを表示する', () => {
    wrap(<UserResults users={[]} query="test" />);
    expect(screen.getByText('結果が見つかりませんでした')).toBeInTheDocument();
    expect(screen.getByText('"test"')).toBeInTheDocument();
  });

  it('UserResultsがユーザーを表示する', () => {
    wrap(<UserResults users={[makeUser(1), makeUser(2)]} query="test" />);
    expect(screen.getByText('User 1')).toBeInTheDocument();
    expect(screen.getByText('User 2')).toBeInTheDocument();
  });

  it('PostResultsが空の場合NoResultsを表示する', () => {
    wrap(<PostResults posts={[]} total={0} query="react" />);
    expect(screen.getByText('結果が見つかりませんでした')).toBeInTheDocument();
    expect(screen.getByText('"react"')).toBeInTheDocument();
  });

  it('PostResultsが投稿を表示する', () => {
    wrap(<PostResults posts={[makePost(1)]} total={1} query="test" />);
    expect(screen.getByText('Post 1')).toBeInTheDocument();
  });

  it('PostResultsがtotal>posts.lengthの場合に件数を表示する', () => {
    wrap(<PostResults posts={[makePost(1)]} total={50} query="test" />);
    expect(screen.getByText('50件の検索結果')).toBeInTheDocument();
  });

  it('CircleResultsが空の場合NoResultsを表示する', () => {
    wrap(<CircleResults circles={[]} query="go" />);
    expect(screen.getByText('結果が見つかりませんでした')).toBeInTheDocument();
    expect(screen.getByText('"go"')).toBeInTheDocument();
  });

  it('CircleResultsがサークルを表示する', () => {
    wrap(<CircleResults circles={[makeCircle(1)]} query="test" />);
    expect(screen.getByText('Circle 1')).toBeInTheDocument();
  });

  it('LoadingResultsがcirclesタブでスケルトンを表示する', () => {
    const { container } = wrap(<LoadingResults tab="circles" />);
    expect(container.querySelector('.space-y-4')).toBeInTheDocument();
  });
});
