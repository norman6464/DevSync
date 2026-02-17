import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import PageHeader from '../PageHeader';

describe('PageHeader', () => {
  it('タイトルが表示される', () => {
    render(<PageHeader title="テストタイトル" />);
    expect(screen.getByText('テストタイトル')).toBeInTheDocument();
  });

  it('サブタイトルが表示される', () => {
    render(<PageHeader title="タイトル" subtitle="サブタイトル" />);
    expect(screen.getByText('サブタイトル')).toBeInTheDocument();
  });

  it('サブタイトルが未指定の場合表示されない', () => {
    render(<PageHeader title="タイトル" />);
    expect(screen.queryByText('サブタイトル')).toBeNull();
  });

  it('アクションボタンが表示される', () => {
    const onAction = vi.fn();
    render(<PageHeader title="タイトル" actionLabel="追加" onAction={onAction} />);
    const button = screen.getByText('追加');
    expect(button).toBeInTheDocument();
  });

  it('アクションボタンクリックでonActionが呼ばれる', () => {
    const onAction = vi.fn();
    render(<PageHeader title="タイトル" actionLabel="追加" onAction={onAction} />);
    fireEvent.click(screen.getByText('追加'));
    expect(onAction).toHaveBeenCalledOnce();
  });

  it('actionLabelのみでonActionが未指定の場合ボタンが表示されない', () => {
    render(<PageHeader title="タイトル" actionLabel="追加" />);
    expect(screen.queryByText('追加')).toBeNull();
  });

  it('onActionのみでactionLabelが未指定の場合ボタンが表示されない', () => {
    render(<PageHeader title="タイトル" onAction={() => {}} />);
    expect(screen.queryByRole('button')).toBeNull();
  });

  it('Plusアイコンが表示される', () => {
    const { container } = render(
      <PageHeader title="タイトル" actionLabel="追加" onAction={() => {}} />
    );
    expect(container.querySelector('svg')).toBeInTheDocument();
  });
});
