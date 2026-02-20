import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import PostFilterPanel from '../PostFilterPanel';
import type { PostSearchFilters } from '../../../api/posts';

const defaultFilters: PostSearchFilters = {
  sortBy: 'latest',
  tags: [],
};

describe('PostFilterPanel', () => {
  it('ソートボタンが3つ表示される', () => {
    const onChange = vi.fn();
    render(<PostFilterPanel filters={defaultFilters} onFiltersChange={onChange} />);

    expect(screen.getByText('最新順')).toBeInTheDocument();
    expect(screen.getByText('人気順')).toBeInTheDocument();
    expect(screen.getByText('閲覧数順')).toBeInTheDocument();
  });

  it('ソートボタンクリックでフィルターが変更される', () => {
    const onChange = vi.fn();
    render(<PostFilterPanel filters={defaultFilters} onFiltersChange={onChange} />);

    fireEvent.click(screen.getByText('人気順'));
    expect(onChange).toHaveBeenCalledWith({ ...defaultFilters, sortBy: 'popular' });
  });

  it('タグを追加できる', () => {
    const onChange = vi.fn();
    render(<PostFilterPanel filters={defaultFilters} onFiltersChange={onChange} />);

    const input = screen.getByPlaceholderText('タグを入力...');
    fireEvent.change(input, { target: { value: 'React' } });
    fireEvent.click(screen.getByText('追加'));

    expect(onChange).toHaveBeenCalledWith({ ...defaultFilters, tags: ['React'] });
  });

  it('Enterキーでタグを追加できる', () => {
    const onChange = vi.fn();
    render(<PostFilterPanel filters={defaultFilters} onFiltersChange={onChange} />);

    const input = screen.getByPlaceholderText('タグを入力...');
    fireEvent.change(input, { target: { value: 'TypeScript' } });
    fireEvent.keyDown(input, { key: 'Enter' });

    expect(onChange).toHaveBeenCalledWith({ ...defaultFilters, tags: ['TypeScript'] });
  });

  it('既存タグが表示される', () => {
    const onChange = vi.fn();
    const filters = { ...defaultFilters, tags: ['React', 'Vue'] };
    render(<PostFilterPanel filters={filters} onFiltersChange={onChange} />);

    expect(screen.getByText('#React')).toBeInTheDocument();
    expect(screen.getByText('#Vue')).toBeInTheDocument();
  });

  it('タグの削除ボタンでタグが除去される', () => {
    const onChange = vi.fn();
    const filters = { ...defaultFilters, tags: ['React', 'Vue'] };
    render(<PostFilterPanel filters={filters} onFiltersChange={onChange} />);

    const removeButtons = screen.getAllByText('×');
    fireEvent.click(removeButtons[0]);

    expect(onChange).toHaveBeenCalledWith({ ...defaultFilters, tags: ['Vue'] });
  });

  it('日付フィルターが変更される', () => {
    const onChange = vi.fn();
    const { container } = render(<PostFilterPanel filters={defaultFilters} onFiltersChange={onChange} />);

    const dateInputs = container.querySelectorAll('input[type="date"]');
    fireEvent.change(dateInputs[0], { target: { value: '2026-01-01' } });

    expect(onChange).toHaveBeenCalledWith({ ...defaultFilters, dateFrom: '2026-01-01' });
  });
});
