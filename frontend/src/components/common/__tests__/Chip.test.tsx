import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Chip from '../Chip';

describe('Chip', () => {
  it('ラベルが表示される', () => {
    render(<Chip label="React" />);
    expect(screen.getByText('React')).toBeInTheDocument();
  });

  it('選択状態のスタイルが適用される', () => {
    const { container } = render(<Chip label="React" selected />);
    expect(container.querySelector('.bg-blue-600')).toBeInTheDocument();
  });

  it('非選択状態のスタイルが適用される', () => {
    const { container } = render(<Chip label="React" />);
    expect(container.querySelector('.bg-gray-800')).toBeInTheDocument();
  });

  it('クリックでコールバックが呼ばれる', async () => {
    const onClick = vi.fn();
    const user = userEvent.setup();
    render(<Chip label="React" onClick={onClick} />);
    await user.click(screen.getByText('React'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('削除ボタンが表示される', () => {
    const { container } = render(<Chip label="React" onDelete={() => {}} />);
    expect(container.querySelector('.lucide-x')).toBeInTheDocument();
  });

  it('削除ボタンクリックでコールバックが呼ばれる', async () => {
    const onDelete = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<Chip label="React" onDelete={onDelete} />);
    const deleteBtn = container.querySelector('.lucide-x')!.parentElement!;
    await user.click(deleteBtn);
    expect(onDelete).toHaveBeenCalledTimes(1);
  });

  it('アイコンが表示される', () => {
    const { container } = render(<Chip label="React" icon={<span data-testid="icon">★</span>} />);
    expect(screen.getByTestId('icon')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<Chip label="React" size="sm" />);
    expect(container.querySelector('.text-xs')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<Chip label="React" size="lg" />);
    expect(container.querySelector('.text-base')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Chip label="React" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('無効状態が適用される', () => {
    const { container } = render(<Chip label="React" disabled />);
    expect(container.querySelector('.opacity-50')).toBeInTheDocument();
  });
});
