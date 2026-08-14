import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import ScrollToTop from '../common/ScrollToTop';
import ScrollToTopOnNavigate from '../common/ScrollToTopOnNavigate';
import PomodoroTimer from '../common/PomodoroTimer';
import { useChatStore } from '../../store/chatStore';

export default function Layout() {
  // 通知はどの画面でもリアルタイムに届くよう、ログイン後の全画面で接続する。
  // 接続の有無はストアの最新値で判定する（描画時点の値だと二重接続になり得る）。
  useEffect(() => {
    const { socket, connect, disconnect } = useChatStore.getState();
    if (!socket) {
      connect();
    }
    return () => disconnect();
  }, []);

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      <ScrollToTopOnNavigate />
      <Header />
      <main className="max-w-7xl mx-auto px-4 py-6">
        <Outlet />
      </main>
      <PomodoroTimer />
      <ScrollToTop />
    </div>
  );
}
