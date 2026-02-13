import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import SearchBar from '../SearchBar';

describe('SearchBar', () => {
  it('検索バーが正しくレンダリングされる', () => {
    const mockOnSearch = vi.fn();
    const mockOnChange = vi.fn();
    render(<SearchBar value="" onChange={mockOnChange} onSearch={mockOnSearch} />);

    const input = screen.getByPlaceholderText(/search/i);
    expect(input).toBeInTheDocument();
  });

  it('入力値が変更されると onChange が呼ばれる', () => {
    const mockOnSearch = vi.fn();
    const mockOnChange = vi.fn();
    render(<SearchBar value="" onChange={mockOnChange} onSearch={mockOnSearch} />);

    const input = screen.getByPlaceholderText(/search/i);
    fireEvent.change(input, { target: { value: 'test query' } });

    expect(mockOnChange).toHaveBeenCalledWith('test query');
  });

  it('フォーム送信で onSearch が呼ばれる', () => {
    const mockOnSearch = vi.fn();
    const mockOnChange = vi.fn();
    const { container } = render(<SearchBar value="test" onChange={mockOnChange} onSearch={mockOnSearch} />);

    const form = container.querySelector('form');
    fireEvent.submit(form!);

    expect(mockOnSearch).toHaveBeenCalled();
  });

  it('検索アイコンが表示される', () => {
    const mockOnSearch = vi.fn();
    const mockOnChange = vi.fn();
    const { container } = render(<SearchBar value="" onChange={mockOnChange} onSearch={mockOnSearch} />);

    const icon = container.querySelector('svg');
    expect(icon).toBeInTheDocument();
  });
});
