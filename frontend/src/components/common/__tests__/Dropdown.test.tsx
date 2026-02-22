import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Dropdown from '../Dropdown';

const items = [
  { id: '1', label: '編集' },
  { id: '2', label: '複製' },
  { id: 'divider' as const },
  { id: '3', label: '削除' },
];

describe('Dropdown', () => {
  it('トリガーボタンが表示される', () => {
    render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    expect(screen.getByText('メニュー')).toBeInTheDocument();
  });

  it('初期状態でメニューが非表示', () => {
    render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    expect(screen.queryByText('編集')).not.toBeInTheDocument();
  });

  it('トリガークリックでメニューが表示される', async () => {
    const user = userEvent.setup();
    render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    await user.click(screen.getByText('メニュー'));

    expect(screen.getByText('編集')).toBeInTheDocument();
    expect(screen.getByText('複製')).toBeInTheDocument();
    expect(screen.getByText('削除')).toBeInTheDocument();
  });

  it('メニュー項目クリックでコールバックが呼ばれる', async () => {
    const onSelect = vi.fn();
    const user = userEvent.setup();
    render(<Dropdown trigger="メニュー" items={items} onSelect={onSelect} />);

    await user.click(screen.getByText('メニュー'));
    await user.click(screen.getByText('編集'));

    expect(onSelect).toHaveBeenCalledWith('1');
  });

  it('メニュー項目クリック後にメニューが閉じる', async () => {
    const user = userEvent.setup();
    render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    await user.click(screen.getByText('メニュー'));
    await user.click(screen.getByText('編集'));

    expect(screen.queryByText('編集')).not.toBeInTheDocument();
  });

  it('区切り線が表示される', async () => {
    const user = userEvent.setup();
    const { container } = render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    await user.click(screen.getByText('メニュー'));

    const dividers = container.querySelectorAll('.border-t');
    expect(dividers.length).toBeGreaterThanOrEqual(1);
  });

  it('無効な項目がクリックできない', async () => {
    const onSelect = vi.fn();
    const user = userEvent.setup();
    const disabledItems = [
      { id: '1', label: '無効', disabled: true },
    ];
    render(<Dropdown trigger="メニュー" items={disabledItems} onSelect={onSelect} />);

    await user.click(screen.getByText('メニュー'));
    await user.click(screen.getByText('無効'));

    expect(onSelect).not.toHaveBeenCalled();
  });

  it('トリガーを再クリックでメニューが閉じる', async () => {
    const user = userEvent.setup();
    render(<Dropdown trigger="メニュー" items={items} onSelect={() => {}} />);

    await user.click(screen.getByText('メニュー'));
    expect(screen.getByText('編集')).toBeInTheDocument();

    await user.click(screen.getByText('メニュー'));
    expect(screen.queryByText('編集')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <Dropdown trigger="メニュー" items={items} onSelect={() => {}} className="custom-class" />
    );

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ReactNodeをトリガーに使用できる', () => {
    render(
      <Dropdown
        trigger={<span data-testid="custom-trigger">カスタム</span>}
        items={items}
        onSelect={() => {}}
      />
    );

    expect(screen.getByTestId('custom-trigger')).toBeInTheDocument();
  });
});
