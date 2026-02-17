import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import Pagination from '../Pagination';

describe('Pagination', () => {
  const defaultProps = {
    currentPage: 0,
    totalItems: 100,
    itemsPerPage: 10,
    onPageChange: vi.fn(),
  };

  it('現在のページ番号と総ページ数を表示する', () => {
    render(<Pagination {...defaultProps} />);
    expect(screen.getByText('1 / 10')).toBeInTheDocument();
  });

  it('前へボタンと次へボタンを表示する', () => {
    render(<Pagination {...defaultProps} />);
    expect(screen.getByRole('button', { name: '前へ' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '次へ' })).toBeInTheDocument();
  });

  it('最初のページでは前へボタンが無効', () => {
    render(<Pagination {...defaultProps} currentPage={0} />);
    expect(screen.getByRole('button', { name: '前へ' })).toBeDisabled();
  });

  it('最後のページでは次へボタンが無効', () => {
    render(<Pagination {...defaultProps} currentPage={9} />);
    expect(screen.getByRole('button', { name: '次へ' })).toBeDisabled();
  });

  it('次へボタンクリックで次のページに遷移', async () => {
    const onPageChange = vi.fn();
    render(<Pagination {...defaultProps} currentPage={0} onPageChange={onPageChange} />);

    await userEvent.click(screen.getByRole('button', { name: '次へ' }));
    expect(onPageChange).toHaveBeenCalledWith(1);
  });

  it('前へボタンクリックで前のページに遷移', async () => {
    const onPageChange = vi.fn();
    render(<Pagination {...defaultProps} currentPage={5} onPageChange={onPageChange} />);

    await userEvent.click(screen.getByRole('button', { name: '前へ' }));
    expect(onPageChange).toHaveBeenCalledWith(4);
  });

  it('totalItemsがitemsPerPage以下の場合は表示しない', () => {
    const { container } = render(
      <Pagination {...defaultProps} totalItems={10} itemsPerPage={10} />
    );
    expect(container.firstChild).toBeNull();
  });

  it('totalItemsが0の場合は表示しない', () => {
    const { container } = render(
      <Pagination {...defaultProps} totalItems={0} />
    );
    expect(container.firstChild).toBeNull();
  });

  it('端数がある場合の総ページ数を正しく計算する', () => {
    render(<Pagination {...defaultProps} totalItems={25} itemsPerPage={10} />);
    expect(screen.getByText('1 / 3')).toBeInTheDocument();
  });

  it('中間ページで両方のボタンが有効', () => {
    render(<Pagination {...defaultProps} currentPage={5} />);
    expect(screen.getByRole('button', { name: '前へ' })).not.toBeDisabled();
    expect(screen.getByRole('button', { name: '次へ' })).not.toBeDisabled();
  });

  it('ページ番号リンクをクリックで直接遷移', async () => {
    const onPageChange = vi.fn();
    render(<Pagination {...defaultProps} currentPage={0} onPageChange={onPageChange} />);

    // ページ番号ボタン「2」をクリック（0-indexedで1）
    await userEvent.click(screen.getByRole('button', { name: '2' }));
    expect(onPageChange).toHaveBeenCalledWith(1);
  });

  it('現在のページ番号がアクティブスタイルで表示される', () => {
    render(<Pagination {...defaultProps} currentPage={2} />);
    const activeButton = screen.getByRole('button', { name: '3' });
    expect(activeButton).toHaveClass('bg-gray-600');
  });

  it('ページ数が多い場合に省略記号を表示する', () => {
    render(<Pagination {...defaultProps} currentPage={5} totalItems={200} itemsPerPage={10} />);
    expect(screen.getAllByText('…').length).toBeGreaterThanOrEqual(1);
  });
});
