import { useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import { useThemeStore } from './store/themeStore';
import ProtectedRoute from './components/auth/ProtectedRoute';
import Layout from './components/layout/Layout';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import DashboardPage from './pages/DashboardPage';
import ProfilePage from './pages/ProfilePage';
import SearchPage from './pages/SearchPage';
import SettingsPage from './pages/SettingsPage';
import RankingsPage from './pages/RankingsPage';
import ChatPage from './pages/ChatPage';
import PostDetailPage from './pages/PostDetailPage';
import FollowListPage from './pages/FollowListPage';
import GitHubCallbackPage from './pages/GitHubCallbackPage';

export default function App() {
  const { isAuthenticated, loadUser } = useAuthStore();
  const { theme, resolvedTheme } = useThemeStore();

  // Initialize theme on mount
  useEffect(() => {
    if (resolvedTheme === 'dark') {
      document.documentElement.classList.add('dark');
      document.documentElement.classList.remove('light');
    } else {
      document.documentElement.classList.add('light');
      document.documentElement.classList.remove('dark');
    }
  }, [resolvedTheme]);

  useEffect(() => {
    if (isAuthenticated) {
      loadUser();
    }
  }, []);

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/github/callback" element={<GitHubCallbackPage />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<Layout />}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/profile/:id" element={<ProfilePage />} />
          <Route path="/profile/:id/followers" element={<FollowListPage />} />
          <Route path="/profile/:id/following" element={<FollowListPage />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="/rankings" element={<RankingsPage />} />
          <Route path="/chat" element={<ChatPage />} />
          <Route path="/chat/:userId" element={<ChatPage />} />
          <Route path="/posts/:id" element={<PostDetailPage />} />
        </Route>
      </Route>
    </Routes>
  );
}
