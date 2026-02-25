import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Stopwatch from '../Stopwatch';

describe('Stopwatch', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('初期表示が00:00.00', () => {
    render(<Stopwatch />);
    expect(screen.getByText('00:00.00')).toBeInTheDocument();
  });

  it('スタートボタンが表示される', () => {
    render(<Stopwatch />);
    expect(screen.getByText('スタート')).toBeInTheDocument();
  });

  it('スタート後にストップボタンが表示される', () => {
    render(<Stopwatch />);
    fireEvent.click(screen.getByText('スタート'));
    expect(screen.getByText('ストップ')).toBeInTheDocument();
  });

  it('スタート後に時間が進む', () => {
    render(<Stopwatch />);
    fireEvent.click(screen.getByText('スタート'));
    act(() => { vi.advanceTimersByTime(1500); });
    expect(screen.queryByText('00:00.00')).not.toBeInTheDocument();
  });

  it('ストップで時間が止まる', () => {
    render(<Stopwatch />);
    fireEvent.click(screen.getByText('スタート'));
    act(() => { vi.advanceTimersByTime(1000); });
    fireEvent.click(screen.getByText('ストップ'));
    const timeText = screen.getByTestId('stopwatch-time').textContent;
    act(() => { vi.advanceTimersByTime(1000); });
    expect(screen.getByTestId('stopwatch-time').textContent).toBe(timeText);
  });

  it('リセットで00:00.00に戻る', () => {
    render(<Stopwatch />);
    fireEvent.click(screen.getByText('スタート'));
    act(() => { vi.advanceTimersByTime(2000); });
    fireEvent.click(screen.getByText('ストップ'));
    fireEvent.click(screen.getByText('リセット'));
    expect(screen.getByText('00:00.00')).toBeInTheDocument();
  });

  it('ラップボタンでラップタイムが記録される', () => {
    render(<Stopwatch showLap />);
    fireEvent.click(screen.getByText('スタート'));
    act(() => { vi.advanceTimersByTime(1000); });
    fireEvent.click(screen.getByText('ラップ'));
    expect(screen.getByTestId('lap-list')).toBeInTheDocument();
  });

  it('ラベルが表示される', () => {
    render(<Stopwatch label="学習時間" />);
    expect(screen.getByText('学習時間')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Stopwatch className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
