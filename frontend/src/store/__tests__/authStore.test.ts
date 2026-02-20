import { describe, it, expect, beforeEach, vi } from 'vitest';
import type { User, AuthResponse } from '../../types/user';

// authApi モック
vi.mock('../../api/auth', () => ({
  login: vi.fn(),
  register: vi.fn(),
  getGitHubLoginURL: vi.fn(),
  gitHubLoginCallback: vi.fn(),
  logout: vi.fn(),
  getMe: vi.fn(),
}));

// isHttpUrl モック
vi.mock('../../utils/url', () => ({
  isHttpUrl: vi.fn(),
}));

import { useAuthStore } from '../authStore';
import * as authApi from '../../api/auth';
import { isHttpUrl } from '../../utils/url';

const mockUser: User = {
  id: 1,
  username: 'testuser',
  name: 'Test User',
  email: 'test@example.com',
  avatar_url: '',
  bio: '',
  github_id: 0,
  github_username: '',
  github_connected: false,
  spotify_connected: false,
  zenn_username: '',
  qiita_username: '',
  atcoder_username: '',
  paiza_rank: '',
  skills_languages: '',
  skills_frameworks: '',
  onboarding_completed: false,
  email_weekly_report: false,
  email_language: 'ja',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const mockAuthResponse = { data: { user: mockUser } as AuthResponse };

describe('authStore', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
      loading: true,
    });
  });

  describe('初期状態', () => {
    it('デフォルトはuser: null, isAuthenticated: false, loading: true', () => {
      const { user, isAuthenticated, loading } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(isAuthenticated).toBe(false);
      expect(loading).toBe(true);
    });
  });

  describe('login', () => {
    it('ログイン成功時にユーザー情報と認証状態を設定する', async () => {
      vi.mocked(authApi.login).mockResolvedValue(mockAuthResponse);

      await useAuthStore.getState().login('test@example.com', 'password');

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toEqual(mockUser);
      expect(isAuthenticated).toBe(true);
      expect(authApi.login).toHaveBeenCalledWith('test@example.com', 'password');
    });

    it('ログイン失敗時にエラーがスローされる', async () => {
      vi.mocked(authApi.login).mockRejectedValue(new Error('Invalid credentials'));

      await expect(useAuthStore.getState().login('bad@example.com', 'wrong')).rejects.toThrow('Invalid credentials');

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(isAuthenticated).toBe(false);
    });
  });

  describe('register', () => {
    it('登録成功時にユーザー情報と認証状態を設定する', async () => {
      vi.mocked(authApi.register).mockResolvedValue(mockAuthResponse);

      await useAuthStore.getState().register('Test User', 'testuser', 'test@example.com', 'password');

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toEqual(mockUser);
      expect(isAuthenticated).toBe(true);
      expect(authApi.register).toHaveBeenCalledWith('Test User', 'testuser', 'test@example.com', 'password');
    });

    it('登録失敗時にエラーがスローされる', async () => {
      vi.mocked(authApi.register).mockRejectedValue(new Error('Email already exists'));

      await expect(
        useAuthStore.getState().register('Test', 'test', 'dup@example.com', 'pass'),
      ).rejects.toThrow('Email already exists');

      expect(useAuthStore.getState().isAuthenticated).toBe(false);
    });
  });

  describe('loginWithGitHub', () => {
    it('有効なURLの場合にリダイレクトする', async () => {
      const githubUrl = 'https://github.com/login/oauth/authorize?client_id=xxx';
      vi.mocked(authApi.getGitHubLoginURL).mockResolvedValue({ data: { url: githubUrl } });
      vi.mocked(isHttpUrl).mockReturnValue(true);

      // window.location.hrefへの代入をキャプチャ
      const hrefSetter = vi.fn();
      Object.defineProperty(window, 'location', {
        value: { href: '' },
        writable: true,
        configurable: true,
      });
      Object.defineProperty(window.location, 'href', {
        set: hrefSetter,
        get: () => '',
        configurable: true,
      });

      await useAuthStore.getState().loginWithGitHub();

      expect(authApi.getGitHubLoginURL).toHaveBeenCalled();
      expect(isHttpUrl).toHaveBeenCalledWith(githubUrl);
      expect(hrefSetter).toHaveBeenCalledWith(githubUrl);
    });

    it('無効なURLの場合にエラーをスローする', async () => {
      vi.mocked(authApi.getGitHubLoginURL).mockResolvedValue({ data: { url: 'javascript:alert(1)' } });
      vi.mocked(isHttpUrl).mockReturnValue(false);

      await expect(useAuthStore.getState().loginWithGitHub()).rejects.toThrow('Invalid OAuth URL');
    });
  });

  describe('handleGitHubCallback', () => {
    it('コールバック処理成功時にユーザー情報を設定する', async () => {
      vi.mocked(authApi.gitHubLoginCallback).mockResolvedValue(mockAuthResponse);

      await useAuthStore.getState().handleGitHubCallback('code123', 'state456');

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toEqual(mockUser);
      expect(isAuthenticated).toBe(true);
      expect(authApi.gitHubLoginCallback).toHaveBeenCalledWith('code123', 'state456');
    });

    it('コールバック失敗時にエラーがスローされる', async () => {
      vi.mocked(authApi.gitHubLoginCallback).mockRejectedValue(new Error('Invalid code'));

      await expect(
        useAuthStore.getState().handleGitHubCallback('bad', 'state'),
      ).rejects.toThrow('Invalid code');

      expect(useAuthStore.getState().isAuthenticated).toBe(false);
    });
  });

  describe('logout', () => {
    it('ログアウト成功時に状態をクリアする', async () => {
      useAuthStore.setState({ user: mockUser, isAuthenticated: true });
      vi.mocked(authApi.logout).mockResolvedValue({ data: { message: 'ok' } });

      await useAuthStore.getState().logout();

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(isAuthenticated).toBe(false);
      expect(authApi.logout).toHaveBeenCalled();
    });

    it('APIエラー時でも状態をクリアする', async () => {
      useAuthStore.setState({ user: mockUser, isAuthenticated: true });
      vi.mocked(authApi.logout).mockRejectedValue(new Error('Network error'));
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      await useAuthStore.getState().logout();

      const { user, isAuthenticated } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(isAuthenticated).toBe(false);
      expect(consoleSpy).toHaveBeenCalled();
      consoleSpy.mockRestore();
    });
  });

  describe('loadUser', () => {
    it('ユーザー読み込み成功時に状態を設定する', async () => {
      vi.mocked(authApi.getMe).mockResolvedValue({ data: mockUser });

      await useAuthStore.getState().loadUser();

      const { user, isAuthenticated, loading } = useAuthStore.getState();
      expect(user).toEqual(mockUser);
      expect(isAuthenticated).toBe(true);
      expect(loading).toBe(false);
    });

    it('読み込み中はloadingがtrueになる', async () => {
      let resolvePromise: (value: { data: User }) => void;
      vi.mocked(authApi.getMe).mockImplementation(
        () => new Promise((resolve) => { resolvePromise = resolve; }),
      );

      const loadPromise = useAuthStore.getState().loadUser();
      expect(useAuthStore.getState().loading).toBe(true);

      resolvePromise!({ data: mockUser });
      await loadPromise;
      expect(useAuthStore.getState().loading).toBe(false);
    });

    it('読み込み失敗時に未認証状態にする', async () => {
      vi.mocked(authApi.getMe).mockRejectedValue(new Error('Unauthorized'));

      await useAuthStore.getState().loadUser();

      const { user, isAuthenticated, loading } = useAuthStore.getState();
      expect(user).toBeNull();
      expect(isAuthenticated).toBe(false);
      expect(loading).toBe(false);
    });
  });

  describe('setUser', () => {
    it('ユーザーを直接設定する', () => {
      useAuthStore.getState().setUser(mockUser);

      expect(useAuthStore.getState().user).toEqual(mockUser);
    });

    it('既存ユーザーを上書きする', () => {
      useAuthStore.setState({ user: mockUser });

      const updatedUser = { ...mockUser, name: 'Updated Name' };
      useAuthStore.getState().setUser(updatedUser);

      expect(useAuthStore.getState().user?.name).toBe('Updated Name');
    });
  });
});
