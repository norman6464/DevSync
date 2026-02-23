import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import HeatMap from '../HeatMap';

const data = [
  { date: '2026-01-01', count: 0 },
  { date: '2026-01-02', count: 3 },
  { date: '2026-01-03', count: 7 },
  { date: '2026-01-04', count: 12 },
];

describe('HeatMap', () => {
  it('ヒートマップが表示される', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells.length).toBe(4);
  });

  it('活動量0はグレー', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells[0]).toHaveClass('bg-gray-800');
  });

  it('低活動量は薄い緑', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells[1]).toHaveClass('bg-green-900');
  });

  it('中活動量は中程度の緑', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells[2]).toHaveClass('bg-green-700');
  });

  it('高活動量は濃い緑', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells[3]).toHaveClass('bg-green-500');
  });

  it('凡例が表示される', () => {
    render(<HeatMap data={data} showLegend />);
    expect(screen.getByText('少')).toBeInTheDocument();
    expect(screen.getByText('多')).toBeInTheDocument();
  });

  it('セルサイズがカスタマイズできる', () => {
    const { container } = render(<HeatMap data={data} cellSize="lg" />);
    expect(container.querySelector('.w-5')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<HeatMap data={data} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('空のデータでも表示される', () => {
    const { container } = render(<HeatMap data={[]} />);
    expect(container.querySelector('.custom-class')).not.toBeInTheDocument();
  });

  it('ツールチップ用のdata属性が設定される', () => {
    const { container } = render(<HeatMap data={data} />);
    const cells = container.querySelectorAll('[data-testid="heatmap-cell"]');
    expect(cells[1]).toHaveAttribute('data-date', '2026-01-02');
    expect(cells[1]).toHaveAttribute('data-count', '3');
  });
});
