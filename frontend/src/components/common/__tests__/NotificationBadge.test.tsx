import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import NotificationBadge from '../NotificationBadge';

describe('NotificationBadge', () => {
  it('子要素が表示される', () => {
    render(
      <NotificationBadge count={5}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.getByText('アイコン')).toBeInTheDocument();
  });

  it('カウントが表示される', () => {
    render(
      <NotificationBadge count={5}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.getByText('5')).toBeInTheDocument();
  });

  it('カウントが0の場合はバッジ非表示', () => {
    render(
      <NotificationBadge count={0}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });

  it('最大数を超えると+表示', () => {
    render(
      <NotificationBadge count={100} max={99}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.getByText('99+')).toBeInTheDocument();
  });

  it('ドットモードでは数値が表示されない', () => {
    const { container } = render(
      <NotificationBadge count={5} dot>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.queryByText('5')).not.toBeInTheDocument();
    const dot = container.querySelector('.w-2');
    expect(dot).toBeInTheDocument();
  });

  it('バッジの色がカスタマイズできる', () => {
    const { container } = render(
      <NotificationBadge count={5} color="green">
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(container.querySelector('.bg-green-500')).toBeInTheDocument();
  });

  it('デフォルトのバッジ色は赤', () => {
    const { container } = render(
      <NotificationBadge count={5}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(container.querySelector('.bg-red-500')).toBeInTheDocument();
  });

  it('右上に配置される（デフォルト）', () => {
    const { container } = render(
      <NotificationBadge count={5}>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(container.querySelector('.-top-1.-right-1')).toBeInTheDocument();
  });

  it('パルスアニメーションが適用される', () => {
    const { container } = render(
      <NotificationBadge count={5} pulse>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <NotificationBadge count={5} className="custom-class">
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('showZeroでカウント0でもバッジ表示', () => {
    render(
      <NotificationBadge count={0} showZero>
        <span>アイコン</span>
      </NotificationBadge>
    );

    expect(screen.getByText('0')).toBeInTheDocument();
  });
});
