import { useState, useRef, useCallback } from 'react';
import { Play, Square, RotateCcw, Flag } from 'lucide-react';

interface StopwatchProps {
  showLap?: boolean;
  label?: string;
  className?: string;
}

function formatTime(ms: number): string {
  const minutes = Math.floor(ms / 60000);
  const seconds = Math.floor((ms % 60000) / 1000);
  const centiseconds = Math.floor((ms % 1000) / 10);
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}.${String(centiseconds).padStart(2, '0')}`;
}

export default function Stopwatch({
  showLap = false,
  label,
  className = '',
}: StopwatchProps) {
  const [elapsed, setElapsed] = useState(0);
  const [running, setRunning] = useState(false);
  const [laps, setLaps] = useState<number[]>([]);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const startTimeRef = useRef(0);

  const start = useCallback(() => {
    setRunning(true);
    startTimeRef.current = Date.now() - elapsed;
    intervalRef.current = setInterval(() => {
      setElapsed(Date.now() - startTimeRef.current);
    }, 10);
  }, [elapsed]);

  const stop = useCallback(() => {
    setRunning(false);
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, []);

  const reset = useCallback(() => {
    stop();
    setElapsed(0);
    setLaps([]);
  }, [stop]);

  const lap = useCallback(() => {
    setLaps((prev) => [...prev, elapsed]);
  }, [elapsed]);

  return (
    <div className={`text-center ${className}`.trim()}>
      {label && <p className="text-sm text-gray-400 mb-2">{label}</p>}
      <p data-testid="stopwatch-time" className="text-4xl font-mono text-gray-200 mb-4">
        {formatTime(elapsed)}
      </p>
      <div className="flex justify-center gap-3">
        {!running ? (
          <button
            type="button"
            onClick={start}
            className="flex items-center gap-1 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg text-sm"
          >
            <Play className="w-4 h-4" />
            スタート
          </button>
        ) : (
          <>
            <button
              type="button"
              onClick={stop}
              className="flex items-center gap-1 px-4 py-2 bg-red-600 hover:bg-red-500 text-white rounded-lg text-sm"
            >
              <Square className="w-4 h-4" />
              ストップ
            </button>
            {showLap && (
              <button
                type="button"
                onClick={lap}
                className="flex items-center gap-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-lg text-sm"
              >
                <Flag className="w-4 h-4" />
                ラップ
              </button>
            )}
          </>
        )}
        {!running && elapsed > 0 && (
          <button
            type="button"
            onClick={reset}
            className="flex items-center gap-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-lg text-sm"
          >
            <RotateCcw className="w-4 h-4" />
            リセット
          </button>
        )}
      </div>
      {laps.length > 0 && (
        <div data-testid="lap-list" className="mt-4 space-y-1">
          {laps.map((lapTime, i) => (
            <div key={i} className="flex justify-between text-sm text-gray-400 px-4">
              <span>ラップ {i + 1}</span>
              <span className="font-mono">{formatTime(lapTime)}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
