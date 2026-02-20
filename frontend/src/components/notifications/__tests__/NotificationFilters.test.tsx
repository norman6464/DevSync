import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import NotificationFilters from '../NotificationFilters';

describe('NotificationFilters', () => {
  const defaultProps = {
    filterType: '' as const,
    setFilterType: vi.fn(),
    showUnreadOnly: false,
    onToggleUnreadOnly: vi.fn(),
  };

  it('全8つのフィルターボタンを表示する', () => {
    render(<NotificationFilters {...defaultProps} />);
    expect(screen.getByText('すべて')).toBeInTheDocument();
    expect(screen.getByText('投稿')).toBeInTheDocument();
    expect(screen.getByText('いいね')).toBeInTheDocument();
    expect(screen.getByText('コメント')).toBeInTheDocument();
    expect(screen.getByText('フォロー')).toBeInTheDocument();
    expect(screen.getByText('メッセージ')).toBeInTheDocument();
    expect(screen.getByText('Q&A回答')).toBeInTheDocument();
    expect(screen.getByText('バッジ')).toBeInTheDocument();
  });

  it('フィルターボタンクリック時にsetFilterTypeが呼ばれる', () => {
    const setFilterType = vi.fn();
    render(<NotificationFilters {...defaultProps} setFilterType={setFilterType} />);
    fireEvent.click(screen.getByText('いいね'));
    expect(setFilterType).toHaveBeenCalledWith('like');
  });

  it('未読のみトグルボタンを表示する', () => {
    render(<NotificationFilters {...defaultProps} />);
    expect(screen.getByText('未読のみ表示')).toBeInTheDocument();
  });

  it('未読のみトグルクリック時にonToggleUnreadOnlyが呼ばれる', () => {
    const onToggleUnreadOnly = vi.fn();
    render(<NotificationFilters {...defaultProps} onToggleUnreadOnly={onToggleUnreadOnly} />);
    fireEvent.click(screen.getByText('未読のみ表示'));
    expect(onToggleUnreadOnly).toHaveBeenCalledTimes(1);
  });

  it('選択中のフィルターにaria-pressed="true"が設定される', () => {
    render(<NotificationFilters {...defaultProps} filterType="post" />);
    expect(screen.getByText('投稿')).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByText('すべて')).toHaveAttribute('aria-pressed', 'false');
  });

  it('showUnreadOnly=trueの場合にaria-pressed="true"が設定される', () => {
    render(<NotificationFilters {...defaultProps} showUnreadOnly={true} />);
    expect(screen.getByText('未読のみ表示')).toHaveAttribute('aria-pressed', 'true');
  });
});
