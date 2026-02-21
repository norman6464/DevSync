import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { usePomodoroTimer } from '../usePomodoroTimer';

vi.mock('../../utils/playNotificationSound', () => ({
  playFocusComplete: vi.fn(),
  playBreakComplete: vi.fn(),
}));

describe('usePomodoroTimer', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('初期状態が正しいこと', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    expect(result.current.phase).toBe('idle');
    expect(result.current.timeLeft).toBe(25 * 60);
    expect(result.current.isRunning).toBe(false);
    expect(result.current.completedCycles).toBe(0);
    expect(result.current.progress).toBe(0);
    expect(result.current.totalTime).toBe(25 * 60);
  });

  it('start()でフェーズがfocusに変わりisRunningがtrueになること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.start();
    });

    expect(result.current.phase).toBe('focus');
    expect(result.current.isRunning).toBe(true);
  });

  it('pause()でisRunningがfalseになること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.start();
    });
    expect(result.current.isRunning).toBe(true);

    act(() => {
      result.current.pause();
    });
    expect(result.current.isRunning).toBe(false);
    expect(result.current.phase).toBe('focus');
  });

  it('reset()で全状態がリセットされること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.start();
    });

    act(() => {
      result.current.reset();
    });

    expect(result.current.phase).toBe('idle');
    expect(result.current.timeLeft).toBe(25 * 60);
    expect(result.current.isRunning).toBe(false);
    expect(result.current.completedCycles).toBe(0);
  });

  it('skip()でfocusからshortBreakに遷移すること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.start();
    });
    expect(result.current.phase).toBe('focus');

    act(() => {
      result.current.skip();
    });

    expect(result.current.phase).toBe('shortBreak');
    expect(result.current.timeLeft).toBe(5 * 60);
    expect(result.current.completedCycles).toBe(1);
  });

  it('skip()でshortBreakからfocusに遷移すること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    // focus開始 → skip → shortBreak
    act(() => {
      result.current.start();
    });
    act(() => {
      result.current.skip();
    });
    expect(result.current.phase).toBe('shortBreak');

    // shortBreak → skip → focus
    act(() => {
      result.current.skip();
    });
    expect(result.current.phase).toBe('focus');
    expect(result.current.timeLeft).toBe(25 * 60);
  });

  it('4サイクル後にlongBreakに遷移すること', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.start();
    });

    // 4回focusをスキップ（focus→break→focus→break→focus→break→focus→?）
    for (let i = 0; i < 3; i++) {
      act(() => { result.current.skip(); }); // focus → shortBreak
      act(() => { result.current.skip(); }); // shortBreak → focus
    }
    // 4回目のfocus skip
    act(() => { result.current.skip(); });

    expect(result.current.phase).toBe('longBreak');
    expect(result.current.timeLeft).toBe(15 * 60);
    expect(result.current.completedCycles).toBe(4);
  });

  it('idle時のskip()は何もしないこと', () => {
    const { result } = renderHook(() => usePomodoroTimer());

    act(() => {
      result.current.skip();
    });

    expect(result.current.phase).toBe('idle');
    expect(result.current.completedCycles).toBe(0);
  });

  it('onFocusCompleteコールバックがfocusスキップ時に呼ばれること', () => {
    const onFocusComplete = vi.fn();
    const { result } = renderHook(() => usePomodoroTimer({ onFocusComplete }));

    act(() => {
      result.current.start();
    });
    act(() => {
      result.current.skip();
    });

    expect(onFocusComplete).toHaveBeenCalledTimes(1);
  });
});
