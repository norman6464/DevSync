import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import ProfileRepositoriesSection from '../ProfileRepositoriesSection';

const makeRepo = (overrides = {}) => ({
  id: 1,
  name: 'devsync',
  full_name: 'user/devsync',
  description: 'エンジニア向けプラットフォーム',
  language: 'TypeScript',
  stars: 42,
  forks: 5,
  ...overrides,
});

const defaultProps = {
  repos: [makeRepo()],
  githubUsername: 'testuser',
};

const renderSection = (props = {}) =>
  render(<ProfileRepositoriesSection {...defaultProps} {...props} />);

describe('ProfileRepositoriesSection', () => {
  it('リポジトリが空の場合何も表示しない', () => {
    const { container } = renderSection({ repos: [] });
    expect(container.innerHTML).toBe('');
  });

  it('リポジトリ名が表示される', () => {
    renderSection();
    expect(screen.getByText('devsync')).toBeInTheDocument();
  });

  it('リポジトリの説明が表示される', () => {
    renderSection();
    expect(screen.getByText('エンジニア向けプラットフォーム')).toBeInTheDocument();
  });

  it('言語が表示される', () => {
    renderSection();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
  });

  it('スター数が表示される', () => {
    renderSection();
    expect(screen.getByText('42')).toBeInTheDocument();
  });

  it('フォーク数が表示される', () => {
    renderSection();
    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('GitHubリンクが正しいURLを持つ', () => {
    renderSection();
    const link = screen.getByText('devsync').closest('a');
    expect(link).toHaveAttribute('href', 'https://github.com/user/devsync');
  });

  it('GitHubで全て見るリンクが表示される', () => {
    renderSection();
    expect(screen.getByText('GitHubで全て見る')).toBeInTheDocument();
  });

  it('最大6件まで表示される', () => {
    const repos = Array.from({ length: 8 }, (_, i) =>
      makeRepo({ id: i + 1, name: `repo-${i + 1}`, full_name: `user/repo-${i + 1}` })
    );
    renderSection({ repos });
    expect(screen.getByText('repo-1')).toBeInTheDocument();
    expect(screen.getByText('repo-6')).toBeInTheDocument();
    expect(screen.queryByText('repo-7')).not.toBeInTheDocument();
  });
});
