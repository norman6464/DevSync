import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import PriceTag from '../PriceTag';

describe('PriceTag', () => {
  it('価格が表示される', () => {
    render(<PriceTag price={1000} />);
    expect(screen.getByText('¥1,000')).toBeInTheDocument();
  });

  it('通貨記号がカスタマイズできる', () => {
    render(<PriceTag price={100} currency="$" />);
    expect(screen.getByText('$100')).toBeInTheDocument();
  });

  it('元の価格が表示される', () => {
    render(<PriceTag price={800} originalPrice={1000} />);
    expect(screen.getByText('¥1,000')).toBeInTheDocument();
  });

  it('元の価格に取り消し線が付く', () => {
    const { container } = render(<PriceTag price={800} originalPrice={1000} />);
    expect(container.querySelector('.line-through')).toBeInTheDocument();
  });

  it('割引率が表示される', () => {
    render(<PriceTag price={800} originalPrice={1000} showDiscount />);
    expect(screen.getByText('-20%')).toBeInTheDocument();
  });

  it('無料の場合は無料と表示される', () => {
    render(<PriceTag price={0} />);
    expect(screen.getByText('無料')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<PriceTag price={1000} size="sm" />);
    expect(container.querySelector('.text-sm')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<PriceTag price={1000} size="lg" />);
    expect(container.querySelector('.text-2xl')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<PriceTag price={1000} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<PriceTag price={1000} label="月額" />);
    expect(screen.getByText('月額')).toBeInTheDocument();
  });
});
