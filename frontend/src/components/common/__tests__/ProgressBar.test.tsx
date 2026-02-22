import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import ProgressBar from '../ProgressBar';

describe('ProgressBar', () => {
  it('プログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} />);

    const progressBar = container.querySelector('[role="progressbar"]');
    expect(progressBar).toBeInTheDocument();
  });

  it('パーセンテージが正しく計算される', () => {
    const { container } = render(<ProgressBar value={30} max={100} />);

    const progressBar = container.querySelector('[role="progressbar"]');
    expect(progressBar).toHaveAttribute('aria-valuenow', '30');
    expect(progressBar).toHaveAttribute('aria-valuemin', '0');
    expect(progressBar).toHaveAttribute('aria-valuemax', '100');
  });

  it('進捗バーの幅がパーセンテージに応じて設定される', () => {
    const { container } = render(<ProgressBar value={75} max={100} />);

    const bar = container.querySelector('.bg-blue-500');
    expect(bar).toHaveStyle({ width: '75%' });
  });

  it('ラベルが表示される', () => {
    render(<ProgressBar value={50} max={100} label="進捗状況" />);

    expect(screen.getByText('進捗状況')).toBeInTheDocument();
  });

  it('パーセンテージテキストが表示される', () => {
    render(<ProgressBar value={60} max={100} showPercentage />);

    expect(screen.getByText('60%')).toBeInTheDocument();
  });

  it('青色のプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} color="blue" />);

    const bar = container.querySelector('.bg-blue-500');
    expect(bar).toBeInTheDocument();
  });

  it('緑色のプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} color="green" />);

    const bar = container.querySelector('.bg-green-500');
    expect(bar).toBeInTheDocument();
  });

  it('黄色のプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} color="yellow" />);

    const bar = container.querySelector('.bg-yellow-500');
    expect(bar).toBeInTheDocument();
  });

  it('赤色のプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} color="red" />);

    const bar = container.querySelector('.bg-red-500');
    expect(bar).toBeInTheDocument();
  });

  it('小サイズのプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} size="sm" />);

    const progressBar = container.querySelector('.h-1');
    expect(progressBar).toBeInTheDocument();
  });

  it('中サイズのプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} size="md" />);

    const progressBar = container.querySelector('.h-2');
    expect(progressBar).toBeInTheDocument();
  });

  it('大サイズのプログレスバーが表示される', () => {
    const { container } = render(<ProgressBar value={50} max={100} size="lg" />);

    const progressBar = container.querySelector('.h-3');
    expect(progressBar).toBeInTheDocument();
  });

  it('100%の進捗が表示される', () => {
    const { container } = render(<ProgressBar value={100} max={100} />);

    const bar = container.querySelector('.bg-blue-500');
    expect(bar).toHaveStyle({ width: '100%' });
  });

  it('0%の進捗が表示される', () => {
    const { container } = render(<ProgressBar value={0} max={100} />);

    const bar = container.querySelector('.bg-blue-500');
    expect(bar).toHaveStyle({ width: '0%' });
  });

  it('トランジション効果がある', () => {
    const { container } = render(<ProgressBar value={50} max={100} />);

    const bar = container.querySelector('.bg-blue-500');
    expect(bar).toHaveClass('transition-all');
  });

  it('背景が灰色になっている', () => {
    const { container } = render(<ProgressBar value={50} max={100} />);

    const background = container.querySelector('.bg-gray-800');
    expect(background).toBeInTheDocument();
  });

  it('角が丸くなっている', () => {
    const { container } = render(<ProgressBar value={50} max={100} />);

    const background = container.querySelector('.rounded-full');
    expect(background).toBeInTheDocument();
  });

  it('ラベルとパーセンテージが同時に表示される', () => {
    render(<ProgressBar value={45} max={100} label="学習進捗" showPercentage />);

    expect(screen.getByText('学習進捗')).toBeInTheDocument();
    expect(screen.getByText('45%')).toBeInTheDocument();
  });

  it('最大値が100以外でも正しく計算される', () => {
    render(<ProgressBar value={3} max={5} showPercentage />);

    expect(screen.getByText('60%')).toBeInTheDocument();
  });
});
