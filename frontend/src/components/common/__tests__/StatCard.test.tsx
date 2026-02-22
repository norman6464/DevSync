import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import StatCard from '../StatCard';

describe('StatCard', () => {
  it('数値が表示される', () => {
    render(<StatCard value={1234} label="合計学習時間" />);

    expect(screen.getByText('1234')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<StatCard value={100} label="投稿数" />);

    expect(screen.getByText('投稿数')).toBeInTheDocument();
  });

  it('アイコンが表示される', () => {
    const { container } = render(<StatCard value={100} label="投稿数" icon="trending-up" />);

    expect(container.querySelector('.lucide-trending-up')).toBeInTheDocument();
  });

  it('増加の変化率が表示される', () => {
    render(<StatCard value={100} label="投稿数" change={12.5} />);

    expect(screen.getByText('+12.5%')).toBeInTheDocument();
  });

  it('減少の変化率が表示される', () => {
    render(<StatCard value={100} label="投稿数" change={-5.3} />);

    expect(screen.getByText('-5.3%')).toBeInTheDocument();
  });

  it('増加はグリーンカラー', () => {
    const { container } = render(<StatCard value={100} label="投稿数" change={10} />);

    expect(container.querySelector('.text-green-400')).toBeInTheDocument();
  });

  it('減少はレッドカラー', () => {
    const { container } = render(<StatCard value={100} label="投稿数" change={-10} />);

    expect(container.querySelector('.text-red-400')).toBeInTheDocument();
  });

  it('変化率0はグレーカラー', () => {
    const { container } = render(<StatCard value={100} label="投稿数" change={0} />);

    expect(container.querySelector('.text-gray-400')).toBeInTheDocument();
  });

  it('フォーマットされた数値が表示される', () => {
    render(<StatCard value={1234567} label="合計" formatted="1,234,567" />);

    expect(screen.getByText('1,234,567')).toBeInTheDocument();
  });

  it('サフィックスが表示される', () => {
    render(<StatCard value={120} label="学習時間" suffix="時間" />);

    expect(screen.getByText('時間')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<StatCard value={100} label="投稿数" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
