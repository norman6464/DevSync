import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import CircularProgress from '../CircularProgress';

describe('CircularProgress', () => {
  it('SVG要素が描画される', () => {
    const { container } = render(<CircularProgress value={50} />);
    expect(container.querySelector('svg')).toBeInTheDocument();
  });

  it('パーセンテージが表示される', () => {
    render(<CircularProgress value={75} showValue />);
    expect(screen.getByText('75%')).toBeInTheDocument();
  });

  it('パーセンテージ非表示が可能', () => {
    render(<CircularProgress value={75} />);
    expect(screen.queryByText('75%')).not.toBeInTheDocument();
  });

  it('0%が正しく表示される', () => {
    render(<CircularProgress value={0} showValue />);
    expect(screen.getByText('0%')).toBeInTheDocument();
  });

  it('100%が正しく表示される', () => {
    render(<CircularProgress value={100} showValue />);
    expect(screen.getByText('100%')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<CircularProgress value={50} size="sm" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '48');
  });

  it('mdサイズが適用される', () => {
    const { container } = render(<CircularProgress value={50} size="md" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '64');
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<CircularProgress value={50} size="lg" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '96');
  });

  it('カスタムカラーが適用される', () => {
    const { container } = render(<CircularProgress value={50} color="#ff0000" />);
    const circle = container.querySelectorAll('circle')[1];
    expect(circle).toHaveAttribute('stroke', '#ff0000');
  });

  it('ラベルが表示される', () => {
    render(<CircularProgress value={50} label="進捗" />);
    expect(screen.getByText('進捗')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<CircularProgress value={50} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
