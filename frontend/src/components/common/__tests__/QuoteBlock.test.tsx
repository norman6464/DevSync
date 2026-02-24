import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import QuoteBlock from '../QuoteBlock';

describe('QuoteBlock', () => {
  it('引用テキストが表示される', () => {
    render(<QuoteBlock quote="知識は力なり" />);
    expect(screen.getByText('知識は力なり')).toBeInTheDocument();
  });

  it('著者名が表示される', () => {
    render(<QuoteBlock quote="知識は力なり" author="フランシス・ベーコン" />);
    expect(screen.getByText('フランシス・ベーコン')).toBeInTheDocument();
  });

  it('著者名がない場合は非表示', () => {
    const { container } = render(<QuoteBlock quote="知識は力なり" />);
    expect(container.querySelectorAll('cite').length).toBe(0);
  });

  it('出典が表示される', () => {
    render(<QuoteBlock quote="テスト" source="書籍名" />);
    expect(screen.getByText('書籍名')).toBeInTheDocument();
  });

  it('引用アイコンが表示される', () => {
    const { container } = render(<QuoteBlock quote="テスト" />);
    expect(container.querySelector('.lucide-quote')).toBeInTheDocument();
  });

  it('blueカラーバリアントが適用される', () => {
    const { container } = render(<QuoteBlock quote="テスト" color="blue" />);
    expect(container.querySelector('.border-blue-500')).toBeInTheDocument();
  });

  it('greenカラーバリアントが適用される', () => {
    const { container } = render(<QuoteBlock quote="テスト" color="green" />);
    expect(container.querySelector('.border-green-500')).toBeInTheDocument();
  });

  it('purpleカラーバリアントが適用される', () => {
    const { container } = render(<QuoteBlock quote="テスト" color="purple" />);
    expect(container.querySelector('.border-purple-500')).toBeInTheDocument();
  });

  it('デフォルトカラーはgray', () => {
    const { container } = render(<QuoteBlock quote="テスト" />);
    expect(container.querySelector('.border-gray-500')).toBeInTheDocument();
  });

  it('blockquote要素が使われる', () => {
    const { container } = render(<QuoteBlock quote="テスト" />);
    expect(container.querySelector('blockquote')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<QuoteBlock quote="テスト" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
