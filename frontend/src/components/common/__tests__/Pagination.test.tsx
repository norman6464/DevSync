import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Pagination from '../Pagination';

describe('Pagination', () => {
  const mockOnPageChange = vi.fn();

  it('ページネーションが表示される', () => {
    render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    expect(screen.getByText('1')).toBeInTheDocument();
  });

  it('ページ番号がクリック可能', () => {
    render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const page2 = screen.getByText('2');
    fireEvent.click(page2);

    expect(mockOnPageChange).toHaveBeenCalledWith(2);
  });

  it('前へボタンが表示される', () => {
    const { container } = render(
      <Pagination currentPage={2} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const prevButton = container.querySelector('[aria-label="前のページ"]');
    expect(prevButton).toBeInTheDocument();
  });

  it('次へボタンが表示される', () => {
    const { container } = render(
      <Pagination currentPage={2} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const nextButton = container.querySelector('[aria-label="次のページ"]');
    expect(nextButton).toBeInTheDocument();
  });

  it('前へボタンをクリックすると前のページに移動', () => {
    const { container } = render(
      <Pagination currentPage={3} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const prevButton = container.querySelector('[aria-label="前のページ"]');
    fireEvent.click(prevButton!);

    expect(mockOnPageChange).toHaveBeenCalledWith(2);
  });

  it('次へボタンをクリックすると次のページに移動', () => {
    const { container } = render(
      <Pagination currentPage={2} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const nextButton = container.querySelector('[aria-label="次のページ"]');
    fireEvent.click(nextButton!);

    expect(mockOnPageChange).toHaveBeenCalledWith(3);
  });

  it('最初のページでは前へボタンが無効化される', () => {
    const { container } = render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const prevButton = container.querySelector('[aria-label="前のページ"]');
    expect(prevButton).toHaveAttribute('disabled');
  });

  it('最後のページでは次へボタンが無効化される', () => {
    const { container } = render(
      <Pagination currentPage={5} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const nextButton = container.querySelector('[aria-label="次のページ"]');
    expect(nextButton).toHaveAttribute('disabled');
  });

  it('現在のページがハイライト表示される', () => {
    render(
      <Pagination currentPage={3} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const currentPage = screen.getByText('3');
    expect(currentPage).toHaveClass('bg-blue-500');
  });

  it('省略記号が表示される（多くのページがある場合）', () => {
    render(
      <Pagination currentPage={5} totalPages={10} onPageChange={mockOnPageChange} />
    );

    expect(screen.getByText('...')).toBeInTheDocument();
  });

  it('ページが1つだけの場合は何も表示されない', () => {
    const { container } = render(
      <Pagination currentPage={1} totalPages={1} onPageChange={mockOnPageChange} />
    );

    const pagination = container.querySelector('nav');
    expect(pagination).not.toBeInTheDocument();
  });

  it('現在のページ情報が表示される', () => {
    render(
      <Pagination currentPage={3} totalPages={10} onPageChange={mockOnPageChange} showInfo />
    );

    expect(screen.getByText(/3 \/ 10/)).toBeInTheDocument();
  });

  it('無効化されたボタンはクリックできない', () => {
    const { container } = render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const prevButton = container.querySelector('[aria-label="前のページ"]');
    fireEvent.click(prevButton!);

    // 無効化されているので呼ばれない
    expect(mockOnPageChange).not.toHaveBeenCalled();
  });

  it('アクティブでないページボタンにホバーエフェクトがある', () => {
    render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const page2 = screen.getByText('2');
    expect(page2).toHaveClass('hover:bg-gray-700');
  });

  it('ページ番号が正しい順序で表示される', () => {
    render(
      <Pagination currentPage={1} totalPages={3} onPageChange={mockOnPageChange} />
    );

    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('前へ・次へボタンにアイコンが表示される', () => {
    const { container } = render(
      <Pagination currentPage={2} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThanOrEqual(2); // ChevronLeft と ChevronRight
  });

  it('複数回クリックしても正しく動作する', () => {
    const { container, rerender } = render(
      <Pagination currentPage={1} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const page2 = screen.getByText('2');
    fireEvent.click(page2);
    expect(mockOnPageChange).toHaveBeenCalledWith(2);

    rerender(
      <Pagination currentPage={2} totalPages={5} onPageChange={mockOnPageChange} />
    );

    const nextButton = container.querySelector('[aria-label="次のページ"]');
    fireEvent.click(nextButton!);
    expect(mockOnPageChange).toHaveBeenCalledWith(3);
  });
});
