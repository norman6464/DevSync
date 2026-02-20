import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import LogFiltersBar from '../LogFiltersBar';

const defaultProps = {
  view: 'list' as const,
  filterCategory: 'all' as const,
  showFavoritesOnly: false,
  sortBy: 'latest' as const,
  filterDate: null as string | null,
  onViewList: vi.fn(),
  onViewCalendar: vi.fn(),
  onToggleFavorites: vi.fn(),
  onFilterCategory: vi.fn(),
  onSortBy: vi.fn(),
  onClearFilterDate: vi.fn(),
};

const renderBar = (props = {}) =>
  render(<LogFiltersBar {...defaultProps} {...props} />);

describe('LogFiltersBar', () => {
  it('リスト・カレンダー・お気に入りボタンが表示される', () => {
    renderBar();
    expect(screen.getByText('リスト')).toBeInTheDocument();
    expect(screen.getByText('カレンダー')).toBeInTheDocument();
    expect(screen.getByText('お気に入り')).toBeInTheDocument();
  });

  it('リストボタンクリックでonViewListが呼ばれる', () => {
    const onViewList = vi.fn();
    renderBar({ onViewList });
    fireEvent.click(screen.getByText('リスト'));
    expect(onViewList).toHaveBeenCalledOnce();
  });

  it('カレンダーボタンクリックでonViewCalendarが呼ばれる', () => {
    const onViewCalendar = vi.fn();
    renderBar({ onViewCalendar });
    fireEvent.click(screen.getByText('カレンダー'));
    expect(onViewCalendar).toHaveBeenCalledOnce();
  });

  it('お気に入りボタンクリックでonToggleFavoritesが呼ばれる', () => {
    const onToggleFavorites = vi.fn();
    renderBar({ onToggleFavorites });
    fireEvent.click(screen.getByText('お気に入り'));
    expect(onToggleFavorites).toHaveBeenCalledOnce();
  });

  it('カテゴリフィルターボタンが表示される', () => {
    renderBar();
    expect(screen.getByText('全て表示')).toBeInTheDocument();
    expect(screen.getByText('コーディング')).toBeInTheDocument();
  });

  it('カテゴリボタンクリックでonFilterCategoryが呼ばれる', () => {
    const onFilterCategory = vi.fn();
    renderBar({ onFilterCategory });
    fireEvent.click(screen.getByText('コーディング'));
    expect(onFilterCategory).toHaveBeenCalledWith('coding');
  });

  it('ソートボタンが表示される', () => {
    renderBar();
    expect(screen.getByText('新しい順')).toBeInTheDocument();
    expect(screen.getByText('古い順')).toBeInTheDocument();
  });

  it('ソートボタンクリックでonSortByが呼ばれる', () => {
    const onSortBy = vi.fn();
    renderBar({ onSortBy });
    fireEvent.click(screen.getByText('古い順'));
    expect(onSortBy).toHaveBeenCalledWith('oldest');
  });

  it('filterDateがあるとき日付フィルター表示・クリアボタンが機能する', () => {
    const onClearFilterDate = vi.fn();
    renderBar({ filterDate: '2026-01-15', onClearFilterDate });
    expect(screen.getByText(/2026-01-15/)).toBeInTheDocument();
    fireEvent.click(screen.getByText('×'));
    expect(onClearFilterDate).toHaveBeenCalledOnce();
  });

  it('filterDateがnullのとき日付フィルターが非表示', () => {
    renderBar({ filterDate: null });
    expect(screen.queryByText('×')).not.toBeInTheDocument();
  });
});
