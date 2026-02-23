import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DataTable from '../DataTable';

const columns = [
  { key: 'name', label: '名前', sortable: true },
  { key: 'age', label: '年齢', sortable: true },
  { key: 'role', label: '役割' },
];

const data = [
  { name: '田中', age: 30, role: 'エンジニア' },
  { name: '佐藤', age: 25, role: 'デザイナー' },
  { name: '鈴木', age: 35, role: 'マネージャー' },
];

describe('DataTable', () => {
  it('ヘッダーが表示される', () => {
    render(<DataTable columns={columns} data={data} />);
    expect(screen.getByText('名前')).toBeInTheDocument();
    expect(screen.getByText('年齢')).toBeInTheDocument();
    expect(screen.getByText('役割')).toBeInTheDocument();
  });

  it('データ行が表示される', () => {
    render(<DataTable columns={columns} data={data} />);
    expect(screen.getByText('田中')).toBeInTheDocument();
    expect(screen.getByText('佐藤')).toBeInTheDocument();
    expect(screen.getByText('鈴木')).toBeInTheDocument();
  });

  it('全セルが表示される', () => {
    render(<DataTable columns={columns} data={data} />);
    expect(screen.getByText('エンジニア')).toBeInTheDocument();
    expect(screen.getByText('デザイナー')).toBeInTheDocument();
    expect(screen.getByText('マネージャー')).toBeInTheDocument();
  });

  it('ソート可能カラムをクリックでソートされる', async () => {
    const user = userEvent.setup();
    render(<DataTable columns={columns} data={data} />);
    await user.click(screen.getByText('年齢'));
    const cells = screen.getAllByRole('cell');
    const ageCells = cells.filter((_, i) => i % 3 === 1);
    expect(ageCells[0].textContent).toBe('25');
  });

  it('2回クリックで降順ソート', async () => {
    const user = userEvent.setup();
    render(<DataTable columns={columns} data={data} />);
    await user.click(screen.getByText('年齢'));
    await user.click(screen.getByText('年齢'));
    const cells = screen.getAllByRole('cell');
    const ageCells = cells.filter((_, i) => i % 3 === 1);
    expect(ageCells[0].textContent).toBe('35');
  });

  it('空データで空メッセージが表示される', () => {
    render(<DataTable columns={columns} data={[]} emptyMessage="データがありません" />);
    expect(screen.getByText('データがありません')).toBeInTheDocument();
  });

  it('行クリックでコールバックが呼ばれる', async () => {
    const onRowClick = vi.fn();
    const user = userEvent.setup();
    render(<DataTable columns={columns} data={data} onRowClick={onRowClick} />);
    await user.click(screen.getByText('田中'));
    expect(onRowClick).toHaveBeenCalledWith(data[0], 0);
  });

  it('ストライプ行が適用される', () => {
    const { container } = render(<DataTable columns={columns} data={data} striped />);
    expect(container.querySelector('.bg-gray-800\\/30')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<DataTable columns={columns} data={data} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ソート不可カラムはクリックしてもソートされない', async () => {
    const user = userEvent.setup();
    render(<DataTable columns={columns} data={data} />);
    await user.click(screen.getByText('役割'));
    const cells = screen.getAllByRole('cell');
    expect(cells[0].textContent).toBe('田中');
  });

  it('テーブル要素が正しく描画される', () => {
    render(<DataTable columns={columns} data={data} />);
    expect(screen.getByRole('table')).toBeInTheDocument();
  });
});
