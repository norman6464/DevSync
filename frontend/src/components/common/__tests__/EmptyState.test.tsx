import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import EmptyState from '../EmptyState';
import { HelpCircle } from 'lucide-react';

describe('EmptyState', () => {
  it('アイコンとメッセージが表示される', () => {
    render(<EmptyState icon={HelpCircle} message="データがありません" />);
    expect(screen.getByText('データがありません')).toBeInTheDocument();
  });

  it('titleが指定された場合に見出しが表示される', () => {
    render(<EmptyState icon={HelpCircle} title="空の状態" message="データがありません" />);
    expect(screen.getByText('空の状態')).toBeInTheDocument();
    expect(screen.getByText('データがありません')).toBeInTheDocument();
  });

  it('titleが未指定の場合は見出しが表示されない', () => {
    render(<EmptyState icon={HelpCircle} message="データがありません" />);
    expect(screen.queryByRole('heading')).toBeNull();
  });

  it('actionLabelとonActionが指定された場合にボタンが表示される', () => {
    const onAction = vi.fn();
    render(
      <EmptyState icon={HelpCircle} message="データがありません" actionLabel="追加する" onAction={onAction} />
    );
    const button = screen.getByText('追加する');
    expect(button).toBeInTheDocument();
    fireEvent.click(button);
    expect(onAction).toHaveBeenCalledOnce();
  });

  it('actionLabelのみ指定された場合はボタンが表示されない', () => {
    render(<EmptyState icon={HelpCircle} message="データがありません" actionLabel="追加する" />);
    expect(screen.queryByText('追加する')).toBeNull();
  });

  it('onActionのみ指定された場合はボタンが表示されない', () => {
    render(<EmptyState icon={HelpCircle} message="データがありません" onAction={() => {}} />);
    // ボタン要素がないことを確認
    expect(screen.queryByRole('button')).toBeNull();
  });
});
