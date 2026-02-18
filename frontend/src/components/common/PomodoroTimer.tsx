import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Timer, Play, Pause, SkipForward, RotateCcw, Volume2, VolumeX, ChevronDown, ChevronUp } from 'lucide-react';
import { selectClass } from '../../constants/styles';
import { usePomodoroTimer } from '../../hooks/usePomodoroTimer';
import type { PomodoroPhase } from '../../hooks/usePomodoroTimer';
import { createLog } from '../../api/learningLogs';
import type { LogCategory } from '../../types/learningLog';

const CATEGORIES: { value: LogCategory; labelKey: string }[] = [
  { value: 'coding', labelKey: 'learningLogs.categoryCoding' },
  { value: 'reading', labelKey: 'learningLogs.categoryReading' },
  { value: 'course', labelKey: 'learningLogs.categoryCourse' },
  { value: 'meetup', labelKey: 'learningLogs.categoryMeetup' },
  { value: 'other', labelKey: 'learningLogs.categoryOther' },
];

/** Business UI: Unified colors */
function getPhaseColor(phase: PomodoroPhase): string {
  return phase === 'idle' ? 'text-gray-400' : 'text-blue-600';
}

function getPhaseRingColor(phase: PomodoroPhase): string {
  return phase === 'idle' ? 'stroke-gray-600' : 'stroke-blue-600';
}

/** mm:ss フォーマット */
function formatTime(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
}

export default function PomodoroTimer() {
  const { t } = useTranslation();
  const [expanded, setExpanded] = useState(false);
  const [soundEnabled, setSoundEnabled] = useState(true);
  const [autoLog, setAutoLog] = useState(true);
  const [category, setCategory] = useState<LogCategory>('coding');

  const handleFocusComplete = useCallback(async () => {
    if (!autoLog) return;
    try {
      await createLog({
        title: t('pomodoro.pomodoroSession'),
        content: t('pomodoro.sessionCompleted', { minutes: 25 }),
        category,
        duration: 25,
        source: 'pomodoro',
      });
    } catch (e) {
      console.warn('Failed to create pomodoro learning log:', e);
    }
  }, [autoLog, category, t]);

  const {
    phase,
    timeLeft,
    isRunning,
    completedCycles,
    progress,
    start,
    pause,
    reset,
    skip,
  } = usePomodoroTimer({
    onFocusComplete: handleFocusComplete,
    soundEnabled,
  });

  const phaseLabel = phase === 'idle'
    ? t('pomodoro.timer')
    : phase === 'focus'
      ? t('pomodoro.focus')
      : phase === 'shortBreak'
        ? t('pomodoro.shortBreak')
        : t('pomodoro.longBreak');

  // SVG円形プログレス
  const radius = 54;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference * (1 - progress);

  // 最小化時
  if (!expanded) {
    return (
      <button
        onClick={() => setExpanded(true)}
        className={`fixed bottom-20 right-6 z-50 flex items-center gap-2 px-4 py-3 rounded-full shadow-sm border transition-colors ${
          isRunning
            ? 'bg-gray-800/95 border-blue-600/50 animate-pulse'
            : 'bg-gray-800/90 border-gray-700 hover:border-gray-600'
        }`}
      >
        <Timer size={18} className={getPhaseColor(phase)} />
        <span className="text-sm font-mono" style={{ color: '#ffffff' }}>{formatTime(timeLeft)}</span>
        {completedCycles > 0 && (
          <span className="text-xs text-gray-400">x{completedCycles}</span>
        )}
      </button>
    );
  }

  return (
    <div className="fixed bottom-20 right-6 z-50 w-72 bg-gray-800/95 backdrop-blur-sm rounded-md shadow-sm border border-gray-700 overflow-hidden">
      {/* ヘッダー */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-gray-700/50">
        <div className="flex items-center gap-2">
          <Timer size={16} className={getPhaseColor(phase)} />
          <span className={`text-sm font-medium ${getPhaseColor(phase)}`}>{phaseLabel}</span>
        </div>
        <div className="flex items-center gap-1">
          <button
            onClick={() => setSoundEnabled(!soundEnabled)}
            className="p-1.5 text-gray-400 hover:text-white rounded transition-colors"
            title={soundEnabled ? t('pomodoro.soundOn') : t('pomodoro.soundOff')}
          >
            {soundEnabled ? <Volume2 size={14} /> : <VolumeX size={14} />}
          </button>
          <button
            onClick={() => setExpanded(false)}
            className="p-1.5 text-gray-400 hover:text-white rounded transition-colors"
          >
            <ChevronDown size={14} />
          </button>
        </div>
      </div>

      {/* タイマー表示 */}
      <div className="flex flex-col items-center py-6">
        <div className="relative w-32 h-32">
          <svg className="w-full h-full -rotate-90" viewBox="0 0 120 120">
            {/* 背景リング */}
            <circle
              cx="60" cy="60" r={radius}
              fill="none"
              stroke="currentColor"
              strokeWidth="6"
              className="text-gray-700"
            />
            {/* プログレスリング */}
            <circle
              cx="60" cy="60" r={radius}
              fill="none"
              strokeWidth="6"
              strokeLinecap="round"
              className={getPhaseRingColor(phase)}
              strokeDasharray={circumference}
              strokeDashoffset={strokeDashoffset}
              style={{ transition: 'stroke-dashoffset 0.3s ease' }}
            />
          </svg>
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <span className="text-3xl font-mono font-bold" style={{ color: '#ffffff' }}>
              {formatTime(timeLeft)}
            </span>
            {completedCycles > 0 && (
              <span className="text-xs text-gray-400 mt-1">
                {t('pomodoro.cyclesCompleted', { count: completedCycles })}
              </span>
            )}
          </div>
        </div>
      </div>

      {/* 操作ボタン */}
      <div className="flex items-center justify-center gap-3 pb-4">
        <button
          onClick={reset}
          className="p-2 text-gray-400 hover:text-white rounded-lg hover:bg-gray-700/50 transition-colors"
          title={t('pomodoro.reset')}
        >
          <RotateCcw size={18} />
        </button>

        <button
          onClick={isRunning ? pause : start}
          className={`p-3 rounded-full transition-colors ${
            isRunning
              ? 'bg-yellow-600 hover:bg-yellow-700'
              : 'bg-red-600 hover:bg-red-700'
          }`}
          style={{ color: '#ffffff' }}
          title={isRunning ? t('pomodoro.pause') : t('pomodoro.start')}
        >
          {isRunning ? <Pause size={20} /> : <Play size={20} />}
        </button>

        <button
          onClick={skip}
          disabled={phase === 'idle'}
          className="p-2 text-gray-400 hover:text-white rounded-lg hover:bg-gray-700/50 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
          title={t('pomodoro.skip')}
        >
          <SkipForward size={18} />
        </button>
      </div>

      {/* 設定パネル */}
      <details className="border-t border-gray-700/50">
        <summary className="flex items-center justify-center gap-1 py-2 cursor-pointer text-xs text-gray-500 hover:text-gray-400 select-none">
          <ChevronUp size={12} className="details-open:rotate-180 transition-transform" />
          {t('settings.title')}
        </summary>
        <div className="px-4 pb-4 space-y-3">
          {/* 自動ログ記録トグル */}
          <label className="flex items-center justify-between cursor-pointer">
            <span className="text-xs text-gray-400">{t('pomodoro.autoLog')}</span>
            <button
              onClick={() => setAutoLog(!autoLog)}
              className={`relative w-9 h-5 rounded-full transition-colors ${
                autoLog ? 'bg-blue-600' : 'bg-gray-600'
              }`}
            >
              <span
                className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform ${
                  autoLog ? 'translate-x-4' : ''
                }`}
              />
            </button>
          </label>

          {/* カテゴリ選択 */}
          {autoLog && (
            <div>
              <label className="text-xs text-gray-400 block mb-1">
                {t('pomodoro.categorySelect')}
              </label>
              <select
                value={category}
                onChange={(e) => setCategory(e.target.value as LogCategory)}
                className={`${selectClass} w-full`}
              >
                {CATEGORIES.map((cat) => (
                  <option key={cat.value} value={cat.value}>
                    {t(cat.labelKey)}
                  </option>
                ))}
              </select>
            </div>
          )}
        </div>
      </details>
    </div>
  );
}
