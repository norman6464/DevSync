import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import GoalFilters from '../GoalFilters';

const defaultProps = {
  filterStatus: 'all' as const,
  setFilterStatus: vi.fn(),
  filterCategory: 'all' as const,
  setFilterCategory: vi.fn(),
  sortBy: 'newest' as const,
  setSortBy: vi.fn(),
};

const renderFilters = (props = {}) =>
  render(<GoalFilters {...defaultProps} {...props} />);

describe('GoalFilters', () => {
  it('フィルターラベルが表示される', () => {
    renderFilters();
    expect(screen.getByText('フィルター')).toBeInTheDocument();
  });

  it('ステータスラベルが表示される', () => {
    renderFilters();
    expect(screen.getByText('ステータス:')).toBeInTheDocument();
  });

  it('ステータスボタンが4つ表示される', () => {
    renderFilters();
    const allButtons = screen.getAllByText('すべて');
    expect(allButtons).toHaveLength(2); // ステータス + カテゴリ
    expect(screen.getByText('進行中')).toBeInTheDocument();
    expect(screen.getByText('一時停止')).toBeInTheDocument();
    expect(screen.getByText('完了')).toBeInTheDocument();
  });

  it('カテゴリボタンが6つ表示される', () => {
    renderFilters();
    expect(screen.getByText('カテゴリ:')).toBeInTheDocument();
    expect(screen.getByText('言語')).toBeInTheDocument();
    expect(screen.getByText('フレームワーク')).toBeInTheDocument();
    expect(screen.getByText('スキル')).toBeInTheDocument();
  });

  it('ソートボタンが4つ表示される', () => {
    renderFilters();
    expect(screen.getByText('新しい順')).toBeInTheDocument();
    expect(screen.getByText('古い順')).toBeInTheDocument();
    expect(screen.getByText('期限順')).toBeInTheDocument();
    expect(screen.getByText('進捗順')).toBeInTheDocument();
  });

  it('ステータスボタンクリックでsetFilterStatusが呼ばれる', () => {
    const setFilterStatus = vi.fn();
    renderFilters({ setFilterStatus });
    fireEvent.click(screen.getByText('一時停止'));
    expect(setFilterStatus).toHaveBeenCalledWith('paused');
  });

  it('カテゴリボタンクリックでsetFilterCategoryが呼ばれる', () => {
    const setFilterCategory = vi.fn();
    renderFilters({ setFilterCategory });
    fireEvent.click(screen.getByText('言語'));
    expect(setFilterCategory).toHaveBeenCalledWith('language');
  });

  it('ソートボタンクリックでsetSortByが呼ばれる', () => {
    const setSortBy = vi.fn();
    renderFilters({ setSortBy });
    fireEvent.click(screen.getByText('期限順'));
    expect(setSortBy).toHaveBeenCalledWith('deadline');
  });

  it('アクティブなステータスボタンにハイライトスタイルが適用される', () => {
    renderFilters({ filterStatus: 'active' });
    const btn = screen.getByText('進行中');
    expect(btn.className).toContain('border-blue-500');
  });

  it('アクティブなカテゴリボタンにハイライトスタイルが適用される', () => {
    renderFilters({ filterCategory: 'framework' });
    const btn = screen.getByText('フレームワーク');
    expect(btn.className).toContain('border-purple-500');
  });

  it('アクティブなソートボタンにハイライトスタイルが適用される', () => {
    renderFilters({ sortBy: 'progress' });
    const btn = screen.getByText('進捗順');
    expect(btn.className).toContain('border-green-500');
  });
});
