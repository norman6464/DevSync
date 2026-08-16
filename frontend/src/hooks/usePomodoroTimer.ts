import { useState, useRef, useCallback, useEffect } from 'react';
import { playFocusComplete, playBreakComplete } from '../utils/playNotificationSound';

/** ポモドーロのフェーズ */
export type PomodoroPhase = 'idle' | 'focus' | 'shortBreak' | 'longBreak';

/** 各フェーズの時間（秒） */
const DURATIONS: Record<Exclude<PomodoroPhase, 'idle'>, number> = {
  focus: 25 * 60,
  shortBreak: 5 * 60,
  longBreak: 15 * 60,
};

/** 長い休憩に切り替わるサイクル数 */
const LONG_BREAK_INTERVAL = 4;

interface UsePomodoroTimerOptions {
  /** 集中セッション完了時のコールバック */
  onFocusComplete?: () => void;
  /** 通知音を鳴らすかどうか */
  soundEnabled?: boolean;
}

export function usePomodoroTimer(options: UsePomodoroTimerOptions = {}) {
  const { onFocusComplete, soundEnabled = true } = options;

  const [phase, setPhase] = useState<PomodoroPhase>('idle');
  const [timeLeft, setTimeLeft] = useState(DURATIONS.focus);
  const [isRunning, setIsRunning] = useState(false);
  const [completedCycles, setCompletedCycles] = useState(0);

  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const targetTimeRef = useRef<number>(0);
  const onFocusCompleteRef = useRef(onFocusComplete);
  const soundEnabledRef = useRef(soundEnabled);
  const phaseRef = useRef(phase);

  // refを最新の値に同期
  useEffect(() => {
    onFocusCompleteRef.current = onFocusComplete;
  }, [onFocusComplete]);

  useEffect(() => {
    soundEnabledRef.current = soundEnabled;
  }, [soundEnabled]);

  useEffect(() => {
    phaseRef.current = phase;
  }, [phase]);

  /** フェーズ完了処理 */
  const handlePhaseComplete = useCallback((currentPhase: PomodoroPhase) => {
    if (currentPhase === 'focus') {
      const newCycles = completedCycles + 1;
      setCompletedCycles(newCycles);

      if (soundEnabledRef.current) playFocusComplete();
      onFocusCompleteRef.current?.();

      // 4サイクルごとに長い休憩
      const nextPhase: PomodoroPhase =
        newCycles % LONG_BREAK_INTERVAL === 0 ? 'longBreak' : 'shortBreak';
      setPhase(nextPhase);
      setTimeLeft(DURATIONS[nextPhase]);
      targetTimeRef.current = Date.now() + DURATIONS[nextPhase] * 1000;
    } else {
      // 休憩完了 → 次の集中フェーズへ
      if (soundEnabledRef.current) playBreakComplete();
      setPhase('focus');
      setTimeLeft(DURATIONS.focus);
      targetTimeRef.current = Date.now() + DURATIONS.focus * 1000;
    }
  }, [completedCycles]);

  const handlePhaseCompleteRef = useRef(handlePhaseComplete);
  useEffect(() => {
    handlePhaseCompleteRef.current = handlePhaseComplete;
  }, [handlePhaseComplete]);

  /** インターバルの開始 */
  const startInterval = useCallback(() => {
    if (intervalRef.current) clearInterval(intervalRef.current);

    intervalRef.current = setInterval(() => {
      const remaining = Math.max(0, Math.ceil((targetTimeRef.current - Date.now()) / 1000));
      setTimeLeft(remaining);

      // フェーズ完了はイベント（tick）起点で処理する。effect 経由で行うと
      // 「effect 内の同期 setState」になり、余計な再レンダーの連鎖も生む。
      if (remaining <= 0 && phaseRef.current !== 'idle') {
        handlePhaseCompleteRef.current(phaseRef.current);
      }
    }, 250);
  }, []);

  /** クリーンアップ */
  useEffect(() => {
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, []);

  /** タイマー開始 */
  const start = useCallback(() => {
    if (phase === 'idle') {
      setPhase('focus');
      setTimeLeft(DURATIONS.focus);
      targetTimeRef.current = Date.now() + DURATIONS.focus * 1000;
    } else {
      targetTimeRef.current = Date.now() + timeLeft * 1000;
    }
    setIsRunning(true);
    startInterval();
  }, [phase, timeLeft, startInterval]);

  /** 一時停止 */
  const pause = useCallback(() => {
    setIsRunning(false);
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, []);

  /** リセット */
  const reset = useCallback(() => {
    setIsRunning(false);
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    setPhase('idle');
    setTimeLeft(DURATIONS.focus);
    setCompletedCycles(0);
  }, []);

  /** 現在のフェーズをスキップ */
  const skip = useCallback(() => {
    if (phase === 'idle') return;
    handlePhaseComplete(phase);
    if (isRunning) {
      startInterval();
    }
  }, [phase, isRunning, handlePhaseComplete, startInterval]);

  /** 現在のフェーズの合計時間（秒） */
  const totalTime = phase === 'idle' ? DURATIONS.focus : DURATIONS[phase];

  /** 進捗率（0〜1） */
  const progress = phase === 'idle' ? 0 : 1 - timeLeft / totalTime;

  return {
    phase,
    timeLeft,
    isRunning,
    completedCycles,
    totalTime,
    progress,
    start,
    pause,
    reset,
    skip,
  };
}
