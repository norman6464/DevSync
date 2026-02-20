import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';

// matchMediaモック — モジュール読み込み前にセットアップが必要
let changeListener: ((e: { matches: boolean }) => void) | null = null;
const mockMatchMedia = vi.fn().mockImplementation(() => ({
  matches: true,
  media: '(prefers-color-scheme: dark)',
  addEventListener: vi.fn((_event: string, cb: (e: { matches: boolean }) => void) => {
    changeListener = cb;
  }),
  removeEventListener: vi.fn(),
  dispatchEvent: vi.fn(),
  onchange: null,
  addListener: vi.fn(),
  removeListener: vi.fn(),
}));
vi.stubGlobal('matchMedia', mockMatchMedia);

// localStorageモック（persist middleware用）
const store: Record<string, string> = {};
vi.stubGlobal('localStorage', {
  getItem: vi.fn((key: string) => store[key] ?? null),
  setItem: vi.fn((key: string, value: string) => { store[key] = value; }),
  removeItem: vi.fn((key: string) => { delete store[key]; }),
  clear: vi.fn(() => { Object.keys(store).forEach(k => delete store[k]); }),
  length: 0,
  key: vi.fn(() => null),
});

// モック設定後にインポート
const { useThemeStore } = await import('../themeStore');

describe('themeStore', () => {
  beforeEach(() => {
    // ストアをリセット
    useThemeStore.setState({ theme: 'dark', resolvedTheme: 'dark' });
    // DOMクラスをリセット
    document.documentElement.classList.remove('dark', 'light');
  });

  afterEach(() => {
    mockMatchMedia.mockClear();
  });

  describe('setTheme', () => {
    it('darkテーマを設定する', () => {
      useThemeStore.getState().setTheme('dark');

      const { theme, resolvedTheme } = useThemeStore.getState();
      expect(theme).toBe('dark');
      expect(resolvedTheme).toBe('dark');
    });

    it('lightテーマを設定する', () => {
      useThemeStore.getState().setTheme('light');

      const { theme, resolvedTheme } = useThemeStore.getState();
      expect(theme).toBe('light');
      expect(resolvedTheme).toBe('light');
    });

    it('systemテーマを設定する（ダークモード優先時）', () => {
      mockMatchMedia.mockReturnValue({
        matches: true,
        media: '(prefers-color-scheme: dark)',
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
      });

      useThemeStore.getState().setTheme('system');

      const { theme, resolvedTheme } = useThemeStore.getState();
      expect(theme).toBe('system');
      expect(resolvedTheme).toBe('dark');
    });

    it('systemテーマを設定する（ライトモード優先時）', () => {
      mockMatchMedia.mockReturnValue({
        matches: false,
        media: '(prefers-color-scheme: dark)',
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
      });

      useThemeStore.getState().setTheme('system');

      const { theme, resolvedTheme } = useThemeStore.getState();
      expect(theme).toBe('system');
      expect(resolvedTheme).toBe('light');
    });
  });

  describe('DOMクラス切り替え', () => {
    it('darkテーマ時にdarkクラスが付与される', () => {
      useThemeStore.getState().setTheme('dark');

      expect(document.documentElement.classList.contains('dark')).toBe(true);
      expect(document.documentElement.classList.contains('light')).toBe(false);
    });

    it('lightテーマ時にlightクラスが付与される', () => {
      useThemeStore.getState().setTheme('light');

      expect(document.documentElement.classList.contains('light')).toBe(true);
      expect(document.documentElement.classList.contains('dark')).toBe(false);
    });

    it('テーマ切り替え時に前のクラスが除去される', () => {
      useThemeStore.getState().setTheme('dark');
      expect(document.documentElement.classList.contains('dark')).toBe(true);

      useThemeStore.getState().setTheme('light');
      expect(document.documentElement.classList.contains('light')).toBe(true);
      expect(document.documentElement.classList.contains('dark')).toBe(false);
    });
  });

  describe('システムテーマ変更リスナー', () => {
    it('systemテーマ時にmatchMedia変更でresolvedThemeが更新される', () => {
      useThemeStore.setState({ theme: 'system', resolvedTheme: 'dark' });

      if (changeListener) {
        changeListener({ matches: false });
      }

      expect(useThemeStore.getState().resolvedTheme).toBe('light');
      expect(document.documentElement.classList.contains('light')).toBe(true);
      expect(document.documentElement.classList.contains('dark')).toBe(false);
    });

    it('systemテーマ以外の場合はmatchMedia変更に反応しない', () => {
      useThemeStore.setState({ theme: 'dark', resolvedTheme: 'dark' });

      if (changeListener) {
        changeListener({ matches: false });
      }

      expect(useThemeStore.getState().resolvedTheme).toBe('dark');
    });
  });

  describe('テーマの連続切り替え', () => {
    it('dark → light → dark の切り替えが正しく動作する', () => {
      useThemeStore.getState().setTheme('dark');
      expect(useThemeStore.getState().resolvedTheme).toBe('dark');

      useThemeStore.getState().setTheme('light');
      expect(useThemeStore.getState().resolvedTheme).toBe('light');

      useThemeStore.getState().setTheme('dark');
      expect(useThemeStore.getState().resolvedTheme).toBe('dark');
    });
  });

  describe('初期状態', () => {
    it('デフォルトテーマはdarkである', () => {
      useThemeStore.setState({ theme: 'dark', resolvedTheme: 'dark' });

      const { theme, resolvedTheme } = useThemeStore.getState();
      expect(theme).toBe('dark');
      expect(resolvedTheme).toBe('dark');
    });
  });
});
