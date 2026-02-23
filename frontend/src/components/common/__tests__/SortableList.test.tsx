import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SortableList from '../SortableList';

const items = [
  { id: '1', label: 'アイテム1' },
  { id: '2', label: 'アイテム2' },
  { id: '3', label: 'アイテム3' },
];

describe('SortableList', () => {
  it('全てのアイテムが表示される', () => {
    render(<SortableList items={items} onReorder={() => {}} />);

    expect(screen.getByText('アイテム1')).toBeInTheDocument();
    expect(screen.getByText('アイテム2')).toBeInTheDocument();
    expect(screen.getByText('アイテム3')).toBeInTheDocument();
  });

  it('上移動ボタンが表示される', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} />);

    const upButtons = container.querySelectorAll('.lucide-chevron-up');
    expect(upButtons.length).toBeGreaterThan(0);
  });

  it('下移動ボタンが表示される', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} />);

    const downButtons = container.querySelectorAll('.lucide-chevron-down');
    expect(downButtons.length).toBeGreaterThan(0);
  });

  it('上移動でアイテムが上に移動する', async () => {
    const onReorder = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<SortableList items={items} onReorder={onReorder} />);

    const upButtons = container.querySelectorAll('[data-testid="move-up"]');
    await user.click(upButtons[1]);

    expect(onReorder).toHaveBeenCalledWith([
      { id: '2', label: 'アイテム2' },
      { id: '1', label: 'アイテム1' },
      { id: '3', label: 'アイテム3' },
    ]);
  });

  it('下移動でアイテムが下に移動する', async () => {
    const onReorder = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<SortableList items={items} onReorder={onReorder} />);

    const downButtons = container.querySelectorAll('[data-testid="move-down"]');
    await user.click(downButtons[0]);

    expect(onReorder).toHaveBeenCalledWith([
      { id: '2', label: 'アイテム2' },
      { id: '1', label: 'アイテム1' },
      { id: '3', label: 'アイテム3' },
    ]);
  });

  it('先頭アイテムの上移動ボタンが無効', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} />);

    const upButtons = container.querySelectorAll('[data-testid="move-up"]');
    expect(upButtons[0]).toBeDisabled();
  });

  it('末尾アイテムの下移動ボタンが無効', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} />);

    const downButtons = container.querySelectorAll('[data-testid="move-down"]');
    expect(downButtons[downButtons.length - 1]).toBeDisabled();
  });

  it('削除ボタンが表示される', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} onRemove={() => {}} />);

    const removeButtons = container.querySelectorAll('[data-testid="remove-item"]');
    expect(removeButtons.length).toBe(3);
  });

  it('削除ボタンクリックでコールバックが呼ばれる', async () => {
    const onRemove = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<SortableList items={items} onReorder={() => {}} onRemove={onRemove} />);

    const removeButtons = container.querySelectorAll('[data-testid="remove-item"]');
    await user.click(removeButtons[0]);

    expect(onRemove).toHaveBeenCalledWith('1');
  });

  it('グリップアイコンが表示される', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} />);

    const grips = container.querySelectorAll('.lucide-grip-vertical');
    expect(grips.length).toBe(3);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<SortableList items={items} onReorder={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
