import { create } from 'zustand';
import type { User } from '../types/user';
import * as authApi from '../api/auth';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  loginWithGitHub: () => Promise<void>;
  handleGitHubCallback: (code: string, state: string) => Promise<void>;
  logout: () => Promise<void>;
  loadUser: () => Promise<void>;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  loading: true,

  login: async (email, password) => {
    const { data } = await authApi.login(email, password);
    set({ user: data.user, isAuthenticated: true });
  },

  register: async (name, email, password) => {
    const { data } = await authApi.register(name, email, password);
    set({ user: data.user, isAuthenticated: true });
  },

  loginWithGitHub: async () => {
    const { data } = await authApi.getGitHubLoginURL();
    window.location.href = data.url;
  },

  handleGitHubCallback: async (code, state) => {
    const { data } = await authApi.gitHubLoginCallback(code, state);
    set({ user: data.user, isAuthenticated: true });
  },

  logout: async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout error:', error);
    }
    set({ user: null, isAuthenticated: false });
  },

  loadUser: async () => {
    try {
      set({ loading: true });
      const { data } = await authApi.getMe();
      set({ user: data, isAuthenticated: true, loading: false });
    } catch {
      set({ user: null, isAuthenticated: false, loading: false });
    }
  },

  setUser: (user) => set({ user }),
}));
