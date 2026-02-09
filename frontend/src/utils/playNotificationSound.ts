/**
 * Web Audio API を使った通知音ユーティリティ。
 * ポモドーロタイマーのフェーズ完了時に使用する。
 */

let audioContext: AudioContext | null = null;

function getAudioContext(): AudioContext | null {
  try {
    if (!audioContext) {
      audioContext = new (window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext)();
    }
    return audioContext;
  } catch {
    return null;
  }
}

function playTone(frequency: number, duration: number, startTime: number, ctx: AudioContext) {
  const oscillator = ctx.createOscillator();
  const gainNode = ctx.createGain();

  oscillator.connect(gainNode);
  gainNode.connect(ctx.destination);

  oscillator.type = 'sine';
  oscillator.frequency.setValueAtTime(frequency, startTime);

  // フェードイン・フェードアウト
  gainNode.gain.setValueAtTime(0, startTime);
  gainNode.gain.linearRampToValueAtTime(0.3, startTime + 0.05);
  gainNode.gain.linearRampToValueAtTime(0, startTime + duration);

  oscillator.start(startTime);
  oscillator.stop(startTime + duration);
}

/**
 * 集中セッション完了時の通知音（明るい3音）。
 * 800Hz → 600Hz → 800Hz
 */
export function playFocusComplete() {
  const ctx = getAudioContext();
  if (!ctx) return;

  const now = ctx.currentTime;
  playTone(800, 0.2, now, ctx);
  playTone(600, 0.2, now + 0.25, ctx);
  playTone(800, 0.3, now + 0.5, ctx);
}

/**
 * 休憩セッション完了時の通知音（柔らかい2音）。
 * 500Hz → 400Hz
 */
export function playBreakComplete() {
  const ctx = getAudioContext();
  if (!ctx) return;

  const now = ctx.currentTime;
  playTone(500, 0.3, now, ctx);
  playTone(400, 0.4, now + 0.35, ctx);
}
