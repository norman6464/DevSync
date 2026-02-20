import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ResourceFilters from '../ResourceFilters';

const defaultProps = {
  searchQuery: '',
  onSearchChange: vi.fn(),
  onSearch: vi.fn(),
  categoryFilter: '' as const,
  onCategoryChange: vi.fn(),
  difficultyFilter: '' as const,
  onDifficultyChange: vi.fn(),
};

describe('ResourceFilters', () => {
  it('検索入力フィールドを表示する', () => {
    render(<ResourceFilters {...defaultProps} />);
    expect(screen.getByPlaceholderText('教材を検索...')).toBeInTheDocument();
  });

  it('カテゴリセレクトを表示する', () => {
    render(<ResourceFilters {...defaultProps} />);
    expect(screen.getByText('全カテゴリ')).toBeInTheDocument();
  });

  it('難易度セレクトを表示する', () => {
    render(<ResourceFilters {...defaultProps} />);
    expect(screen.getByText('全レベル')).toBeInTheDocument();
  });

  it('カテゴリセレクトに9個のオプションがある（全カテゴリ + 8カテゴリ）', () => {
    render(<ResourceFilters {...defaultProps} />);
    const categorySelect = screen.getByText('全カテゴリ').closest('select')!;
    expect(categorySelect.options.length).toBe(9);
  });

  it('難易度セレクトに4個のオプションがある（全レベル + 3レベル）', () => {
    render(<ResourceFilters {...defaultProps} />);
    const difficultySelect = screen.getByText('全レベル').closest('select')!;
    expect(difficultySelect.options.length).toBe(4);
  });

  it('カテゴリ変更でonCategoryChangeが呼ばれる', () => {
    const onCategoryChange = vi.fn();
    render(<ResourceFilters {...defaultProps} onCategoryChange={onCategoryChange} />);
    const categorySelect = screen.getByText('全カテゴリ').closest('select')!;
    fireEvent.change(categorySelect, { target: { value: 'book' } });
    expect(onCategoryChange).toHaveBeenCalledWith('book');
  });

  it('難易度変更でonDifficultyChangeが呼ばれる', () => {
    const onDifficultyChange = vi.fn();
    render(<ResourceFilters {...defaultProps} onDifficultyChange={onDifficultyChange} />);
    const difficultySelect = screen.getByText('全レベル').closest('select')!;
    fireEvent.change(difficultySelect, { target: { value: 'beginner' } });
    expect(onDifficultyChange).toHaveBeenCalledWith('beginner');
  });

  it('検索入力値が表示される', () => {
    render(<ResourceFilters {...defaultProps} searchQuery="React" />);
    expect(screen.getByDisplayValue('React')).toBeInTheDocument();
  });

  it('選択されたカテゴリが反映される', () => {
    render(<ResourceFilters {...defaultProps} categoryFilter="video" />);
    const categorySelect = screen.getByText('全カテゴリ').closest('select')!;
    expect(categorySelect.value).toBe('video');
  });

  it('選択された難易度が反映される', () => {
    render(<ResourceFilters {...defaultProps} difficultyFilter="advanced" />);
    const difficultySelect = screen.getByText('全レベル').closest('select')!;
    expect(difficultySelect.value).toBe('advanced');
  });
});
