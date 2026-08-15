import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import SearchTabs from '../SearchTabs';

describe('SearchTabs', () => {
  it('デフォルトでユーザータブが選択されている', () => {
    const mockOnChange = vi.fn();
    render(<SearchTabs activeTab="users" onTabChange={mockOnChange} counts={{ users: 10, posts: 5, circles: 3 }} />);

    const usersTab = screen.getByText('ユーザー');
    expect(usersTab.parentElement).toHaveClass('border-blue-500');
  });

  it('タブクリックで onTabChange が呼ばれる', () => {
    const mockOnChange = vi.fn();
    render(<SearchTabs activeTab="users" onTabChange={mockOnChange} counts={{ users: 10, posts: 5, circles: 3 }} />);

    const postsTab = screen.getByText('投稿');
    fireEvent.click(postsTab);

    expect(mockOnChange).toHaveBeenCalledWith('posts');
  });

  it('各タブに件数が表示される', () => {
    const mockOnChange = vi.fn();
    render(<SearchTabs activeTab="users" onTabChange={mockOnChange} counts={{ users: 10, posts: 5, circles: 3 }} />);

    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('3つのタブが表示される', () => {
    const mockOnChange = vi.fn();
    const { container } = render(<SearchTabs activeTab="users" onTabChange={mockOnChange} counts={{ users: 0, posts: 0, circles: 0 }} />);

    const tabs = container.querySelectorAll('button');
    expect(tabs.length).toBe(3);
  });
});
