import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import QuickActionsWidget from '../QuickActionsWidget';

vi.mock('../../../store/authStore', () => ({
  useAuthStore: vi.fn(() => ({ user: { id: 1, name: 'Test User' } })),
}));

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('QuickActionsWidget', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('ウィジェットタイトルが表示される', () => {
    renderWithRouter(<QuickActionsWidget />);
    expect(screen.getByText('クイックアクション')).toBeInTheDocument();
  });

  it('複数のアクションカードが表示される', () => {
    renderWithRouter(<QuickActionsWidget />);

    // 主要なアクションが表示されることを確認
    expect(screen.getByText('学習ログを追加')).toBeInTheDocument();
    expect(screen.getByText('ノートを作成')).toBeInTheDocument();
    expect(screen.getByText('目標を追加')).toBeInTheDocument();
    expect(screen.getByText('投稿を作成')).toBeInTheDocument();
  });

  it('各アクションにアイコンが表示される', () => {
    const { container } = renderWithRouter(<QuickActionsWidget />);

    // lucide-reactのアイコンが存在することを確認
    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThan(0);
  });

  it('アクションカードは少なくとも6個以上ある', () => {
    renderWithRouter(<QuickActionsWidget />);

    // アクションカードを取得（リンクまたはボタン）
    const actionCards = screen.getAllByRole('link');
    expect(actionCards.length).toBeGreaterThanOrEqual(6);
  });

  it('学習ログアクションのリンクが正しい', () => {
    renderWithRouter(<QuickActionsWidget />);

    const learningLogLink = screen.getByText('学習ログを追加').closest('a');
    expect(learningLogLink).toHaveAttribute('href', '/learning-logs');
  });

  it('ノート作成アクションのリンクが正しい', () => {
    renderWithRouter(<QuickActionsWidget />);

    const noteLink = screen.getByText('ノートを作成').closest('a');
    expect(noteLink).toHaveAttribute('href', '/notes');
  });

  it('目標追加アクションのリンクが正しい', () => {
    renderWithRouter(<QuickActionsWidget />);

    const goalLink = screen.getByText('目標を追加').closest('a');
    expect(goalLink).toHaveAttribute('href', '/goals');
  });

  it('投稿作成アクションのリンクが正しい', () => {
    renderWithRouter(<QuickActionsWidget />);

    const postLink = screen.getByText('投稿を作成').closest('a');
    expect(postLink).toHaveAttribute('href', '/');
  });

  it('各アクションカードにホバーエフェクトがある', () => {
    renderWithRouter(<QuickActionsWidget />);

    const actionCard = screen.getByText('学習ログを追加').closest('a');
    expect(actionCard).toHaveClass('hover:border-blue-400/50');
  });

  it('グリッドレイアウトで表示される', () => {
    const { container } = renderWithRouter(<QuickActionsWidget />);

    const grid = container.querySelector('.grid');
    expect(grid).toBeInTheDocument();
  });

  it('各アクションに説明文が表示される', () => {
    renderWithRouter(<QuickActionsWidget />);

    // 説明文の一部を確認
    const descriptions = screen.getAllByText(/今日の|新しい|追加|作成/);
    expect(descriptions.length).toBeGreaterThan(0);
  });

  it('アクションカードは2列グリッドで表示される（モバイル以外）', () => {
    const { container } = renderWithRouter(<QuickActionsWidget />);

    const grid = container.querySelector('.grid-cols-2');
    expect(grid).toBeInTheDocument();
  });

  it('全てのアクションカードが一貫したスタイルを持つ', () => {
    renderWithRouter(<QuickActionsWidget />);

    const actionCards = screen.getAllByRole('link');

    actionCards.forEach((card) => {
      expect(card).toHaveClass('bg-gray-800/50');
      expect(card).toHaveClass('border');
      expect(card).toHaveClass('rounded-lg');
    });
  });

  it('アイコンは適切なサイズで表示される', () => {
    const { container } = renderWithRouter(<QuickActionsWidget />);

    const icons = container.querySelectorAll('svg');
    icons.forEach((icon) => {
      // lucide-reactのアイコンはw-5 h-5クラスを持つことを確認
      expect(icon.classList.contains('w-5') || icon.classList.contains('w-4')).toBe(true);
    });
  });
});
