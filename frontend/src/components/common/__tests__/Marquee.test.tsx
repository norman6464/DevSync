import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Marquee from '../Marquee';

describe('Marquee', () => {
  it('テキストが表示される', () => {
    render(<Marquee>スクロールテキスト</Marquee>);
    expect(screen.getByText('スクロールテキスト')).toBeInTheDocument();
  });

  it('アニメーションクラスが適用される', () => {
    const { container } = render(<Marquee>テスト</Marquee>);
    expect(container.querySelector('.animate-marquee')).toBeInTheDocument();
  });

  it('一時停止でアニメーションが停止する', () => {
    const { container } = render(<Marquee pauseOnHover>テスト</Marquee>);
    expect(container.querySelector('.hover\\:animation-paused')).toBeInTheDocument();
  });

  it('速度がカスタマイズできる', () => {
    const { container } = render(<Marquee speed="fast">テスト</Marquee>);
    const el = container.querySelector('[data-speed="fast"]');
    expect(el).toBeInTheDocument();
  });

  it('方向が変更できる', () => {
    const { container } = render(<Marquee direction="right">テスト</Marquee>);
    expect(container.querySelector('.animate-marquee-reverse')).toBeInTheDocument();
  });

  it('オーバーフローが隠される', () => {
    const { container } = render(<Marquee>テスト</Marquee>);
    expect(container.querySelector('.overflow-hidden')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Marquee className="custom-class">テスト</Marquee>);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ReactNodeを子要素に使用できる', () => {
    render(<Marquee><span data-testid="child">子要素</span></Marquee>);
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });
});
