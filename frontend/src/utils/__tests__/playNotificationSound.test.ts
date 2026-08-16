import { describe, it, expect, vi, beforeEach } from 'vitest';

type MockFn = ReturnType<typeof vi.fn>;

interface MockOscillator {
  connect: MockFn;
  type: OscillatorType;
  frequency: { setValueAtTime: MockFn };
  start: MockFn;
  stop: MockFn;
}

interface MockGainNode {
  connect: MockFn;
  gain: {
    setValueAtTime: MockFn;
    linearRampToValueAtTime: MockFn;
  };
}

function createMockAudioContext() {
  const oscillators: MockOscillator[] = [];
  const gainNodes: MockGainNode[] = [];

  function createMockOscillator(): MockOscillator {
    const osc = {
      connect: vi.fn(),
      type: '' as OscillatorType,
      frequency: { setValueAtTime: vi.fn() },
      start: vi.fn(),
      stop: vi.fn(),
    };
    oscillators.push(osc);
    return osc;
  }

  function createMockGainNode(): MockGainNode {
    const gain = {
      connect: vi.fn(),
      gain: {
        setValueAtTime: vi.fn(),
        linearRampToValueAtTime: vi.fn(),
      },
    };
    gainNodes.push(gain);
    return gain;
  }

  const mockCtx = {
    currentTime: 100,
    destination: {},
    createOscillator: vi.fn(createMockOscillator),
    createGain: vi.fn(createMockGainNode),
  };

  return { mockCtx, oscillators, gainNodes };
}

describe('playNotificationSound', () => {
  beforeEach(() => {
    vi.resetModules();
    vi.unstubAllGlobals();
  });

  describe('playFocusComplete', () => {
    it('AudioContext対応環境で3つのトーンを再生する', async () => {
      const { mockCtx } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playFocusComplete } = await import('../playNotificationSound');
      playFocusComplete();

      expect(mockCtx.createOscillator).toHaveBeenCalledTimes(3);
      expect(mockCtx.createGain).toHaveBeenCalledTimes(3);
    });

    it('正しい周波数でトーンを再生する（800, 600, 800Hz）', async () => {
      const { mockCtx, oscillators } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playFocusComplete } = await import('../playNotificationSound');
      playFocusComplete();

      expect(oscillators[0].frequency.setValueAtTime).toHaveBeenCalledWith(800, 100);
      expect(oscillators[1].frequency.setValueAtTime).toHaveBeenCalledWith(600, 100.25);
      expect(oscillators[2].frequency.setValueAtTime).toHaveBeenCalledWith(800, 100.5);
    });

    it('OscillatorNodeをGainNodeとdestinationに接続する', async () => {
      const { mockCtx, oscillators, gainNodes } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playFocusComplete } = await import('../playNotificationSound');
      playFocusComplete();

      oscillators.forEach((osc, i) => {
        expect(osc.connect).toHaveBeenCalledWith(gainNodes[i]);
      });
      gainNodes.forEach((gain) => {
        expect(gain.connect).toHaveBeenCalledWith(mockCtx.destination);
      });
    });

    it('AudioContext未対応環境でエラーを投げない', async () => {
      const { playFocusComplete } = await import('../playNotificationSound');
      expect(() => playFocusComplete()).not.toThrow();
    });
  });

  describe('playBreakComplete', () => {
    it('AudioContext対応環境で2つのトーンを再生する', async () => {
      const { mockCtx } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playBreakComplete } = await import('../playNotificationSound');
      playBreakComplete();

      expect(mockCtx.createOscillator).toHaveBeenCalledTimes(2);
      expect(mockCtx.createGain).toHaveBeenCalledTimes(2);
    });

    it('正しい周波数でトーンを再生する（500, 400Hz）', async () => {
      const { mockCtx, oscillators } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playBreakComplete } = await import('../playNotificationSound');
      playBreakComplete();

      expect(oscillators[0].frequency.setValueAtTime).toHaveBeenCalledWith(500, 100);
      expect(oscillators[1].frequency.setValueAtTime).toHaveBeenCalledWith(400, 100.35);
    });

    it('正しいタイミングでstart/stopを呼ぶ', async () => {
      const { mockCtx, oscillators } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playBreakComplete } = await import('../playNotificationSound');
      playBreakComplete();

      expect(oscillators[0].start).toHaveBeenCalledWith(100);
      expect(oscillators[0].stop).toHaveBeenCalledWith(100.3);
      expect(oscillators[1].start).toHaveBeenCalledWith(100.35);
      expect(oscillators[1].stop).toHaveBeenCalledWith(100.75);
    });

    it('AudioContext未対応環境でエラーを投げない', async () => {
      const { playBreakComplete } = await import('../playNotificationSound');
      expect(() => playBreakComplete()).not.toThrow();
    });

    it('フェードイン・フェードアウトのゲイン設定が正しい', async () => {
      const { mockCtx, gainNodes } = createMockAudioContext();
      vi.stubGlobal('AudioContext', function () { return mockCtx; });

      const { playBreakComplete } = await import('../playNotificationSound');
      playBreakComplete();

      gainNodes.forEach((gain) => {
        expect(gain.gain.setValueAtTime).toHaveBeenCalled();
        expect(gain.gain.linearRampToValueAtTime).toHaveBeenCalledTimes(2);
      });
    });
  });
});
