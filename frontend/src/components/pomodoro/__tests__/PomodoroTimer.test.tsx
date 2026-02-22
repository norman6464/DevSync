import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import PomodoroTimer from '../PomodoroTimer';

vi.useFakeTimers();

describe('PomodoroTimer', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.clearAllTimers();
  });

  it('タイマーが表示される', () => {
    render(<PomodoroTimer />);
    expect(screen.getByText(/25:00/)).toBeInTheDocument();
  });

  it('開始ボタンが表示される', () => {
    render(<PomodoroTimer />);
    expect(screen.getByText('開始')).toBeInTheDocument();
  });

  it('開始ボタンをクリックするとタイマーが開始する', () => {
    render(<PomodoroTimer />);

    const startButton = screen.getByText('開始');
    fireEvent.click(startButton);

    expect(screen.getByText('一時停止')).toBeInTheDocument();
  });

  it('タイマーが1秒ごとにカウントダウンする', () => {
    render(<PomodoroTimer />);

    const startButton = screen.getByText('開始');
    fireEvent.click(startButton);

    act(() => {
      vi.advanceTimersByTime(1000);
    });

    expect(screen.getByText(/24:59/)).toBeInTheDocument();
  });

  it('一時停止ボタンをクリックするとタイマーが停止する', () => {
    render(<PomodoroTimer />);

    const startButton = screen.getByText('開始');
    fireEvent.click(startButton);

    const pauseButton = screen.getByText('一時停止');
    fireEvent.click(pauseButton);

    expect(screen.getByText('再開')).toBeInTheDocument();
  });

  it('リセットボタンが表示される', () => {
    render(<PomodoroTimer />);
    expect(screen.getByText('リセット')).toBeInTheDocument();
  });

  it('リセットボタンをクリックするとタイマーが初期化される', () => {
    render(<PomodoroTimer />);

    const startButton = screen.getByText('開始');
    fireEvent.click(startButton);

    act(() => {
      vi.advanceTimersByTime(5000);
    });

    const resetButton = screen.getByText('リセット');
    fireEvent.click(resetButton);

    expect(screen.getByText(/25:00/)).toBeInTheDocument();
    expect(screen.getByText('開始')).toBeInTheDocument();
  });

  it('作業時間と休憩時間のタブが表示される', () => {
    render(<PomodoroTimer />);
    expect(screen.getByText('作業')).toBeInTheDocument();
    expect(screen.getByText('休憩')).toBeInTheDocument();
  });

  it('休憩タブをクリックすると休憩時間が表示される', () => {
    render(<PomodoroTimer />);

    const breakTab = screen.getByText('休憩');
    fireEvent.click(breakTab);

    expect(screen.getByText(/05:00/)).toBeInTheDocument();
  });

  it('進捗バーが表示される', () => {
    const { container } = render(<PomodoroTimer />);

    const progressBar = container.querySelector('[role="progressbar"]');
    expect(progressBar).toBeInTheDocument();
  });

  it('タイマーアイコンが表示される', () => {
    const { container } = render(<PomodoroTimer />);

    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThan(0);
  });

  it('タイマー終了時にコールバックが呼ばれる', async () => {
    const onComplete = vi.fn();
    render(<PomodoroTimer onComplete={onComplete} />);

    const startButton = screen.getByText('開始');
    fireEvent.click(startButton);

    // 25分 = 1500秒
    act(() => {
      vi.advanceTimersByTime(1500 * 1000);
    });

    await waitFor(() => {
      expect(onComplete).toHaveBeenCalledWith(25);
    });
  });

  it('作業中はタイマーが青色で表示される', () => {
    const { container } = render(<PomodoroTimer />);

    const timer = container.querySelector('.text-blue-400');
    expect(timer).toBeInTheDocument();
  });

  it('タイマーのカスタマイズが可能', () => {
    render(<PomodoroTimer workMinutes={30} breakMinutes={10} />);

    expect(screen.getByText(/30:00/)).toBeInTheDocument();

    const breakTab = screen.getByText('休憩');
    fireEvent.click(breakTab);

    expect(screen.getByText(/10:00/)).toBeInTheDocument();
  });
});
