import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Youtube } from 'lucide-react';
import EmptyState from '../EmptyState';

describe('EmptyState', () => {
  it('タイトルが表示される', () => {
    render(<EmptyState title="データがありません" />);

    expect(screen.getByText('データがありません')).toBeInTheDocument();
  });

  it('説明文が表示される', () => {
    render(<EmptyState title="空" description="まだデータがありません" />);

    expect(screen.getByText('まだデータがありません')).toBeInTheDocument();
  });

  it('アイコンが表示される', () => {
    const { container } = render(<EmptyState title="空" icon="inbox" />);

    const icon = container.querySelector('.lucide-inbox');
    expect(icon).toBeInTheDocument();
  });

  // ページ固有のアイコン（定義済みキーに無いもの）はコンポーネントを直接渡す。
  it('lucide のコンポーネントを直接渡してもアイコンが表示される', () => {
    const { container } = render(<EmptyState title="空" icon={Youtube} />);

    expect(container.querySelector('.lucide-youtube')).toBeInTheDocument();
  });

  it('アクションボタンが表示される', () => {
    render(
      <EmptyState
        title="空"
        actionLabel="追加する"
        onAction={() => {}}
      />
    );

    expect(screen.getByText('追加する')).toBeInTheDocument();
  });

  it('アクションボタンクリックでコールバックが呼ばれる', async () => {
    const onAction = vi.fn();
    const user = userEvent.setup();
    render(
      <EmptyState
        title="空"
        actionLabel="追加する"
        onAction={onAction}
      />
    );

    await user.click(screen.getByText('追加する'));

    expect(onAction).toHaveBeenCalledTimes(1);
  });

  it('カスタムアイコンサイズが適用される', () => {
    const { container } = render(<EmptyState title="空" icon="inbox" iconSize="lg" />);

    const icon = container.querySelector('.w-16');
    expect(icon).toBeInTheDocument();
  });

  it('説明文がない場合は表示されない', () => {
    const { container } = render(<EmptyState title="空" />);

    const paragraphs = container.querySelectorAll('p');
    expect(paragraphs.length).toBe(0);
  });

  it('アクションボタンがない場合は表示されない', () => {
    render(<EmptyState title="空" />);

    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<EmptyState title="空" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('カスタム子要素が表示される', () => {
    render(
      <EmptyState title="空">
        <span data-testid="custom">カスタムコンテンツ</span>
      </EmptyState>
    );

    expect(screen.getByTestId('custom')).toBeInTheDocument();
  });
});
