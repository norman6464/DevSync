import { useState, useEffect, useRef } from 'react';
import { Play, Pause, RotateCcw, Timer } from 'lucide-react';

interface PomodoroTimerProps {
  workMinutes?: number;
  breakMinutes?: number;
  onComplete?: (duration: number) => void;
}

type TimerMode = 'work' | 'break';

export default function PomodoroTimer({
  workMinutes = 25,
  breakMinutes = 5,
  onComplete,
}: PomodoroTimerProps) {
  const [mode, setMode] = useState<TimerMode>('work');
  const [timeLeft, setTimeLeft] = useState(workMinutes * 60);
  const [isRunning, setIsRunning] = useState(false);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const startTimeRef = useRef<number>(0);

  const totalTime = mode === 'work' ? workMinutes * 60 : breakMinutes * 60;
  const progress = ((totalTime - timeLeft) / totalTime) * 100;

  useEffect(() => {
    if (isRunning) {
      startTimeRef.current = Date.now();
      intervalRef.current = setInterval(() => {
        setTimeLeft((prev) => {
          if (prev <= 1) {
            setIsRunning(false);
            if (onComplete && mode === 'work') {
              onComplete(mode === 'work' ? workMinutes : breakMinutes);
            }
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    } else {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isRunning, mode, workMinutes, breakMinutes, onComplete]);

  // モード・時間設定の変更でタイマーをリセットする。effect で setState すると
  // 余計な再レンダーが入るため、公式の「render 中に前回値と比較して調整する」パターンで行う。
  const [prevConfig, setPrevConfig] = useState({ mode, workMinutes, breakMinutes });
  if (prevConfig.mode !== mode || prevConfig.workMinutes !== workMinutes || prevConfig.breakMinutes !== breakMinutes) {
    setPrevConfig({ mode, workMinutes, breakMinutes });
    setTimeLeft(mode === 'work' ? workMinutes * 60 : breakMinutes * 60);
    setIsRunning(false);
  }

  const handleStart = () => {
    setIsRunning(true);
  };

  const handlePause = () => {
    setIsRunning(false);
  };

  const handleReset = () => {
    setIsRunning(false);
    setTimeLeft(mode === 'work' ? workMinutes * 60 : breakMinutes * 60);
  };

  const handleModeChange = (newMode: TimerMode) => {
    setMode(newMode);
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-lg p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-2">
        <Timer className="w-5 h-5 text-blue-400" />
        <h3 className="text-lg font-semibold text-white">ポモドーロタイマー</h3>
      </div>

      {/* Mode Tabs */}
      <div className="flex gap-2">
        <button
          onClick={() => handleModeChange('work')}
          className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-colors ${
            mode === 'work'
              ? 'bg-blue-500/20 text-blue-400 border border-blue-400/30'
              : 'bg-gray-800/50 text-gray-400 hover:text-white border border-gray-700'
          }`}
        >
          作業
        </button>
        <button
          onClick={() => handleModeChange('break')}
          className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-colors ${
            mode === 'break'
              ? 'bg-green-500/20 text-green-400 border border-green-400/30'
              : 'bg-gray-800/50 text-gray-400 hover:text-white border border-gray-700'
          }`}
        >
          休憩
        </button>
      </div>

      {/* Timer Display */}
      <div className="text-center py-8">
        <div
          className={`text-6xl font-bold ${
            mode === 'work' ? 'text-blue-400' : 'text-green-400'
          }`}
        >
          {formatTime(timeLeft)}
        </div>
        <div className="text-sm text-gray-400 mt-2">
          {mode === 'work' ? '作業時間' : '休憩時間'}
        </div>
      </div>

      {/* Progress Bar */}
      <div className="space-y-2">
        <div className="h-2 bg-gray-800 rounded-full overflow-hidden" role="progressbar">
          <div
            className={`h-full transition-all ${
              mode === 'work' ? 'bg-blue-500' : 'bg-green-500'
            }`}
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      {/* Controls */}
      <div className="flex gap-3">
        {!isRunning ? (
          <button
            onClick={handleStart}
            className="flex-1 flex items-center justify-center gap-2 py-3 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium transition-colors"
          >
            <Play className="w-5 h-5" />
            {timeLeft === totalTime ? '開始' : '再開'}
          </button>
        ) : (
          <button
            onClick={handlePause}
            className="flex-1 flex items-center justify-center gap-2 py-3 bg-orange-500 hover:bg-orange-600 text-white rounded-lg font-medium transition-colors"
          >
            <Pause className="w-5 h-5" />
            一時停止
          </button>
        )}

        <button
          onClick={handleReset}
          className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium transition-colors flex items-center gap-2"
        >
          <RotateCcw className="w-5 h-5" />
          リセット
        </button>
      </div>

      {/* Hint */}
      <div className="text-xs text-gray-500 text-center">
        {mode === 'work'
          ? `${workMinutes}分間集中して作業しましょう`
          : `${breakMinutes}分間リラックスしましょう`}
      </div>
    </div>
  );
}
