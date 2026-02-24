import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import CurrencyDisplay from '../CurrencyDisplay';

describe('CurrencyDisplay', () => {
  it('金額が表示される', () => {
    render(<CurrencyDisplay amount={1000} currency="JPY" />);
    expect(screen.getByText(/1,000/)).toBeInTheDocument();
  });

  it('JPYの通貨記号が表示される', () => {
    render(<CurrencyDisplay amount={1000} currency="JPY" />);
    expect(screen.getByText(/¥/)).toBeInTheDocument();
  });

  it('USDの通貨記号が表示される', () => {
    render(<CurrencyDisplay amount={100} currency="USD" />);
    expect(screen.getByText(/\$/)).toBeInTheDocument();
  });

  it('EURの通貨記号が表示される', () => {
    render(<CurrencyDisplay amount={100} currency="EUR" />);
    expect(screen.getByText(/€/)).toBeInTheDocument();
  });

  it('正の変動が緑色で表示される', () => {
    const { container } = render(<CurrencyDisplay amount={1000} currency="JPY" change={5.2} />);
    expect(container.querySelector('.text-green-400')).toBeInTheDocument();
  });

  it('負の変動が赤色で表示される', () => {
    const { container } = render(<CurrencyDisplay amount={1000} currency="JPY" change={-3.1} />);
    expect(container.querySelector('.text-red-400')).toBeInTheDocument();
  });

  it('変動パーセンテージが表示される', () => {
    render(<CurrencyDisplay amount={1000} currency="JPY" change={5.2} />);
    expect(screen.getByText(/5.2%/)).toBeInTheDocument();
  });

  it('変動がない場合はパーセンテージ非表示', () => {
    render(<CurrencyDisplay amount={1000} currency="JPY" />);
    expect(screen.queryByText(/%/)).not.toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<CurrencyDisplay amount={1000} currency="JPY" label="売上" />);
    expect(screen.getByText('売上')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<CurrencyDisplay amount={1000} currency="JPY" size="sm" />);
    expect(container.querySelector('.text-lg')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<CurrencyDisplay amount={1000} currency="JPY" size="lg" />);
    expect(container.querySelector('.text-3xl')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CurrencyDisplay amount={1000} currency="JPY" className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
